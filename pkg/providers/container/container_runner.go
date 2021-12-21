package container

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog"

	"github.com/Skarlso/providers-example/pkg/models"
	"github.com/Skarlso/providers-example/pkg/providers"
)

// Config defines parameters for the Runner.
type Config struct {
	DefaultMaximumCommandRuntime int
}

// Dependencies defines the provider dependencies this provider has.
type Dependencies struct {
	Next   providers.Runner
	Storer providers.Storer
	Logger zerolog.Logger
}

// Runner implements the Run interface for container based runtimes.
type Runner struct {
	Config
	Dependencies

	cli client.APIClient
}

// NewRunner creates a new container based runtime.
func NewRunner(cfg Config, deps Dependencies) (*Runner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		deps.Logger.Debug().Err(err).Msg("Failed to create docker client.")
		return nil, err
	}
	return &Runner{
		Config:       cfg,
		Dependencies: deps,
		cli:          cli,
	}, nil
}

// Run implements the container based runtime details, using Docker as an engine.
func (cr *Runner) Run(ctx context.Context, name string, args []string) error {
	// Find the plugin, get the location, find the type, if it's not container, call next.
	cmd, err := cr.Storer.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("plugin not found: %w", err)
	}
	if cmd.Type != models.Container {
		cr.Logger.Info().Msg("Unknown plugin type, calling next in line.")
		if cr.Next == nil {
			return fmt.Errorf("no next provider configured")
		}
		return cr.Next.Run(ctx, name, args)
	}
	if err := cr.runCommand(cmd.Name, cmd.Container.Image, args); err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}
	return nil
}

// runCommand takes a command name and an image and the necessary arguments and runs the container and waits for output.
func (cr *Runner) runCommand(commandName, image string, args []string) error {
	output, err := cr.cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		cr.Logger.Debug().Err(err).Msg("Failed to pull image.")
		return err
	}
	if _, err := io.Copy(os.Stdout, output); err != nil {
		cr.Logger.Debug().Err(err).Msg("Failed to pull image.")
		return err
	}

	cr.Logger.Info().Msg("Creating container...")
	cont, err := cr.cli.ContainerCreate(context.Background(), &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Image:        image,
		Cmd:          args,
	}, nil, nil, nil, "")
	if err != nil {
		cr.Logger.Debug().Err(err).Strs("warnings", cont.Warnings).Msg("Failed to create container.")
		return err
	}
	cr.startAndWaitForContainer(commandName, cont.ID)
	return nil
}

// runCommand takes a single command and executes it, waiting for it to finish,
// or tcr out. Either way, it will update the corresponding command row.
func (cr *Runner) startAndWaitForContainer(commandName, containerID string) {
	cr.Logger.Info().Str("name", commandName).Msg("Starting running command...")
	done := make(chan error, 1)
	defer func() {
		// we remove the container in a `defer` instead of autoRemove, to be able to read out the logs.
		// If we use AutoRemove, the container is gone by the tcr we want to read the output.
		// Could try streaming the logs. But this is enough for now.
		if err := cr.cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			cr.Logger.Debug().Err(err).Str("container_id", containerID).Msg("Failed to remove container.")
		}
	}()

	cr.Logger.Info().Msg("Starting container...")
	if err := cr.cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{}); err != nil {
		return
	}

	go func() {
		exit, err := cr.cli.ContainerWait(context.Background(), containerID, container.WaitConditionNotRunning)
		select {
		case e := <-err:
			done <- e
		case e := <-exit:
			if e.StatusCode != 0 {
				if e.Error != nil {
					done <- errors.New(e.Error.Message)
				} else {
					done <- fmt.Errorf("status code: %d", e.StatusCode)
				}
			} else {
				done <- nil
			}
		}
	}()

	for {
		select {
		case err := <-done:
			log, logErr := cr.cli.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
				ShowStderr: true,
				ShowStdout: true,
			})
			if logErr != nil {
				return
			}
			buffer := &bytes.Buffer{}
			logs := "no logs available"
			if _, err := stdcopy.StdCopy(buffer, buffer, log); err != nil {
				cr.Logger.Debug().Err(err).Msg("Failed to de-multiplex the docker log.")
			} else {
				logs = buffer.String()
			}

			if err != nil {
				cr.Logger.Debug().Err(err).Msg("Failed to run command.")
				cr.Logger.Debug().Str("logs", logs).Msg("Logs from the attached container.")
				return
			}
			cr.Logger.Info().Msg("Successfully finished command. Output: ")
			fmt.Println(buffer.String())
			return
		case <-time.After(time.Duration(cr.DefaultMaximumCommandRuntime) * time.Second):
			// update entry
			cr.Logger.Error().Msg("Command tcrd out.")
			if err := cr.cli.ContainerKill(context.Background(), containerID, "SIGKILL"); err != nil {
				cr.Logger.Error().Str("container_id", containerID).Msg("Failed to kill process with pid.")
			}
			return
		}
	}
}
