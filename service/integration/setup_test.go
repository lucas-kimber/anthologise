package integration

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

const composePath = "../../compose.yml"

var baseURL string

func buildStack(ctx context.Context) (*compose.DockerCompose, error) {

	opt := compose.WithStackFiles(composePath)
	stack, err := compose.NewDockerComposeWith(opt)

	if err != nil {
		return nil, err
	}

	if err := stack.WithOsEnv().Up(ctx, compose.Wait(true)); err != nil {
		_ = stack.Down(context.Background())
		return nil, err
	}

	return stack, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	stack, err := buildStack(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to start test stack:", err)
		os.Exit(1)
	}

	baseURL = fmt.Sprintf(
		"http://localhost:%s",
		os.Getenv("ANTHOLOGISE_PORT"),
	)

	serviceContainer, err := stack.ServiceContainer(ctx, "service")
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if code != 0 {
		logs, err := serviceContainer.Logs(context.Background())

		if err != nil {
			log.Printf("failed to get service logs: %v", err)
		} else {
			defer logs.Close()

			data, err := io.ReadAll(logs)
			if err != nil {
				log.Printf("failed to read service logs: %v", err)
			} else {
				log.Printf("service logs:\n%s", data)
			}
		}
	}

	if err := stack.Down(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "failed to stop test stack:", err)

		if code == 0 {
			code = 1
		}
	}

	os.Exit(code)
}
