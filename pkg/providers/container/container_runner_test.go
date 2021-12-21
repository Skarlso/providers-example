package container

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/Skarlso/providers-example/pkg/models"
	"github.com/Skarlso/providers-example/pkg/providers/fakes"
)

type mockDockerClient struct {
	client.APIClient

	imagePullOutput io.ReadCloser
	createOutput    containertypes.ContainerCreateCreatedBody
	logsOutput      io.ReadCloser
	containerOkChan chan containertypes.ContainerWaitOKBody
}

func (mc *mockDockerClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return mc.imagePullOutput, nil
}

func (mc *mockDockerClient) ContainerCreate(ctx context.Context, config *containertypes.Config, hostConfig *containertypes.HostConfig, networkingConfig *networktypes.NetworkingConfig, platform *specs.Platform, containerName string) (containertypes.ContainerCreateCreatedBody, error) {
	return mc.createOutput, nil
}

func (mc *mockDockerClient) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	return nil
}

func (mc *mockDockerClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
	return nil
}

func (mc *mockDockerClient) ContainerWait(ctx context.Context, container string, condition containertypes.WaitCondition) (<-chan containertypes.ContainerWaitOKBody, <-chan error) {
	return mc.containerOkChan, make(<-chan error)
}

func (mc *mockDockerClient) ContainerKill(ctx context.Context, container, signal string) error {
	return nil
}

func (mc *mockDockerClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	return mc.logsOutput, nil
}

func TestCreateRun(t *testing.T) {
	logger := zerolog.New(os.Stderr)
	fakeStorer := &fakes.FakeStorer{}
	imagePullOutput := &bytes.Buffer{}
	imagePullOutput.WriteString("success")
	logsOutput := &bytes.Buffer{}
	logsOutput.WriteString("I haz logs.")
	containerOkWaitChannel := make(chan containertypes.ContainerWaitOKBody)
	createOutput := containertypes.ContainerCreateCreatedBody{
		ID: "new-container-id",
	}
	apiClient := &mockDockerClient{
		imagePullOutput: io.NopCloser(imagePullOutput),
		createOutput:    createOutput,
		logsOutput:      io.NopCloser(logsOutput),
		containerOkChan: containerOkWaitChannel,
	}
	r := Runner{
		Dependencies: Dependencies{
			Storer: fakeStorer,
			Logger: logger,
		},
		Config: Config{
			DefaultMaximumCommandRuntime: 15,
		},
		cli: apiClient,
	}
	fakeStorer.GetReturns(&models.Plugin{
		ID:   1,
		Name: "test",
		Type: models.Container,
		Container: &models.ContainerPlugin{
			Image: "test-image",
		},
	}, nil)
	go func() {
		apiClient.containerOkChan <- containertypes.ContainerWaitOKBody{
			StatusCode: 0,
		}
	}()
	err := r.Run(context.Background(), "test", []string{"arg1", "arg2"})
	assert.NoError(t, err)
}
