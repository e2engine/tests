package harness

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type ServiceProcess struct {
	cmd    *exec.Cmd
	name   string
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	done   chan error
}

func (p *ServiceProcess) Stop(t *testing.T) {
	t.Helper()

	if p == nil ||
		p.cmd == nil ||
		p.cmd.Process == nil {
		return
	}

	select {
	case err := <-p.done:
		if err != nil {
			t.Errorf(
				"Helper service %q exited unexpectedly: %v\nstdout: %s\nstderr: %s",
				p.name,
				err,
				p.stdout.String(),
				p.stderr.String(),
			)
		}

		return

	default:
	}

	if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
		t.Errorf(
			"Failed to stop helper service %q: %v",
			p.name,
			err,
		)
		return
	}

	select {
	case <-p.done:
		return

	case <-time.After(2 * time.Second):
		if err := p.cmd.Process.Kill(); err != nil {
			t.Errorf(
				"Failed to kill helper service %q: %v",
				p.name,
				err,
			)
			return
		}

		<-p.done
	}
}

const defaultServiceAddress = "127.0.0.1:9000"

var serviceAddresses = map[string]string{
	"gateway-http-path":    defaultServiceAddress,
	"gateway-http-body":    defaultServiceAddress,
	"gateway-grpc-message": defaultServiceAddress,
	"gateway-grpc-method":  defaultServiceAddress,
	"ok-http":              defaultServiceAddress,
	"ok-grpc":              defaultServiceAddress,
}

func (h *Harness) StartService(
	ctx context.Context,
	serviceName string,
	args ...string,
) *ServiceProcess {
	h.t.Helper()

	address, ok := serviceAddresses[serviceName]
	if !ok {
		h.t.Fatalf(
			"Unknown helper service: %s",
			serviceName,
		)
	}

	servicePath := filepath.Join(
		"..",
		"..",
		"services",
		"bin",
		serviceName,
	)

	cmd := exec.CommandContext(
		ctx,
		servicePath,
		args...,
	)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd.Stdout = io.MultiWriter(stdout, os.Stdout)
	cmd.Stderr = io.MultiWriter(stderr, os.Stderr)

	if err := cmd.Start(); err != nil {
		h.t.Fatalf(
			"Failed to start helper service %q: %v",
			serviceName,
			err,
		)
	}

	process := &ServiceProcess{
		cmd:    cmd,
		name:   serviceName,
		stdout: stdout,
		stderr: stderr,
		done:   make(chan error, 1),
	}

	go func() {
		process.done <- cmd.Wait()
		close(process.done)
	}()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()

	dialer := net.Dialer{
		Timeout: 100 * time.Millisecond,
	}

	for {
		conn, err := dialer.DialContext(
			ctx,
			"tcp",
			address,
		)
		if err == nil {
			_ = conn.Close()
			return process
		}

		select {
		case err := <-process.done:
			h.t.Fatalf(
				"Helper service %q exited before becoming ready: %v\nstdout: %s\nstderr: %s",
				serviceName,
				err,
				stdout.String(),
				stderr.String(),
			)

		case <-ctx.Done():
			_ = cmd.Process.Kill()

			h.t.Fatalf(
				"Context cancelled while waiting for helper service %q: %v",
				serviceName,
				ctx.Err(),
			)

		case <-timer.C:
			_ = cmd.Process.Kill()

			h.t.Fatalf(
				"Timed out waiting for helper service %q at %s\nstdout: %s\nstderr: %s",
				serviceName,
				address,
				stdout.String(),
				stderr.String(),
			)

		case <-ticker.C:
		}
	}
}
