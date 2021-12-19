package livestore

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"

	"github.com/Skarlso/providers-example/pkg/providers/storer"
)

var (
	testDbLocation string
)

// TestMain runs the tests for the package and allows us to bring up any external dependencies required.
func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	temp, err := os.MkdirTemp("", "lite_test")
	if err != nil {
		fmt.Println("failed to create temporary folder: ", err)
		return 1
	}

	testDbLocation = temp

	if _, err := storer.NewLiteStorer(zerolog.New(os.Stderr), temp); err != nil {
		fmt.Println("failed to initialize db: ", err)
	}

	defer func() {
		if err := os.RemoveAll(temp); err != nil {
			fmt.Println("failed to remove temp location: ", err)
		}
	}()

	return m.Run()
}
