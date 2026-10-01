package harness

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type Service struct {
	name string
}

func (s *Service) Stop(t *testing.T) {
	t.Helper()

	cmd := exec.Command("docker", "rm", "-f", s.name)

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Errorf(
			"Failed to stop helper service %q: %v, output: %s",
			s.name,
			err,
			output,
		)
	}
}

func (h *Harness) StartService(
	ctx context.Context,
	serviceName string,
	commandArgs ...string,
) *Service {
	h.t.Helper()

	args := []string{
		"run",
		"-d",
		"--network", h.network,
		"--name", serviceName,
		"e2engine-test-" + serviceName + ":local",
	}
	args = append(args, commandArgs...)

	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		h.t.Fatalf(
			"Failed to start helper service %q: %v, stderr: %s",
			serviceName,
			err,
			stderr.String(),
		)
	}

	service := &Service{name: serviceName}

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()

	for {
		inspect := exec.Command(
			"docker",
			"inspect",
			"--format={{.State.Running}}",
			serviceName,
		)

		output, err := inspect.CombinedOutput()
		if err != nil || strings.TrimSpace(string(output)) != "true" {
			logs, _ := exec.Command(
				"docker",
				"logs",
				serviceName,
			).CombinedOutput()

			_ = removeContainer(serviceName)

			h.t.Fatalf(
				"Helper service %q exited before becoming ready\nlogs: %s",
				serviceName,
				string(logs),
			)
		}

		probe := exec.CommandContext(
			ctx,
			"docker",
			"run",
			"--rm",
			"--network", h.network,
			"busybox",
			"nc", "-z", serviceName, "9000",
		)

		if err := probe.Run(); err == nil {
			return service
		}

		select {
		case <-ctx.Done():
			_ = removeContainer(serviceName)

			h.t.Fatalf(
				"Context cancelled while waiting for helper service %q: %v",
				serviceName,
				ctx.Err(),
			)

		case <-timer.C:
			logs, _ := exec.Command(
				"docker",
				"logs",
				serviceName,
			).CombinedOutput()

			_ = removeContainer(serviceName)

			h.t.Fatalf(
				"Timed out waiting for helper service %q\nlogs: %s",
				serviceName,
				string(logs),
			)

		case <-ticker.C:
		}
	}
}

func removeContainer(name string) error {
	return exec.Command("docker", "rm", "-f", name).Run()
}
