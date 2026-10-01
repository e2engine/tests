package tests

import (
	"bufio"
	"context"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/e2engine/core/model"

	"github.com/e2engine/tests/cli/harness"
	"github.com/e2engine/tests/util"
)

func TestTransportKind(t *testing.T) {
	dataDir := "../data/transport-kind"

	tests := []struct {
		name                string
		configPath          string
		envSpec             string
		testSpec            string
		testExecSpec        string
		testExecRunningSpec string
		env                 map[string]string
		services            []Service
		long                bool
		captureLogs         bool
		crashRunner         bool
		allowExpiry         bool
	}{
		{
			name:         "default", // direct transport kind is default
			envSpec:      "real-with-mocked-dependency",
			testSpec:     "gateway-calls-entitlements",
			testExecSpec: "http-real-mocked.yml",
			env:          map[string]string{},
			services: []Service{
				{Name: "gateway-http-path"},
			},
		},
		{
			name:         "direct",
			envSpec:      "real-with-mocked-dependency",
			testSpec:     "gateway-calls-entitlements",
			testExecSpec: "http-real-mocked.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND": "direct",
			},
			services: []Service{
				{Name: "gateway-http-path"},
			},
		},
		{
			name:         "socket",
			envSpec:      "real-with-mocked-dependency",
			testSpec:     "gateway-calls-entitlements",
			testExecSpec: "http-real-mocked.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "127.0.0.1:9001",
			},
			services: []Service{
				{Name: "gateway-http-path"},
			},
		},
		{
			name:                "socket, long",
			envSpec:             "real-with-mocked-dependency",
			testSpec:            "gateway-calls-entitlements",
			testExecSpec:        "http-real-mocked.yml",
			testExecRunningSpec: "http-real-mocked-running.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "127.0.0.1:9001",
			},
			services: []Service{
				{Name: "gateway-http-path", Args: []string{"2"}},
			},
			long: true,
		},
		{
			name:         "socket, grpc, error",
			envSpec:      "socket-grpc-error",
			testSpec:     "socket-grpc-error",
			testExecSpec: "socket-grpc-error.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "127.0.0.1:9001",
			},
			services: []Service{
				{Name: "ok-grpc"},
			},
		},
		{
			name:                "socket, child process crash",
			envSpec:             "real-with-mocked-dependency",
			testSpec:            "gateway-calls-entitlements",
			testExecSpec:        "http-real-mocked.yml",
			testExecRunningSpec: "http-real-mocked-running.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "127.0.0.1:9001",
			},
			services: []Service{
				{Name: "gateway-http-path", Args: []string{"10"}},
			},
			long:        true,
			captureLogs: true,
			crashRunner: true,
			allowExpiry: true,
		},
		{
			name:         "direct, grpc, error",
			envSpec:      "socket-grpc-error",
			testSpec:     "socket-grpc-error",
			testExecSpec: "socket-grpc-error.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND": "direct",
			},
			services: []Service{
				{Name: "ok-grpc"},
			},
		},
		{
			name:         "socket, custom address",
			envSpec:      "real-with-mocked-dependency",
			testSpec:     "gateway-calls-entitlements",
			testExecSpec: "http-real-mocked.yml",
			env: map[string]string{
				"E2ENGINE_TRANSPORT_KIND":    "socket",
				"E2ENGINE_TRANSPORT_ADDRESS": "127.0.0.1:19001",
			},
			services: []Service{
				{Name: "gateway-http-path"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				defaultTestTimeout,
			)
			defer cancel()

			for key, value := range test.env {
				t.Setenv(key, value)
			}

			tempDir := t.TempDir()

			if test.captureLogs {
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_KIND", "file")
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_LEVEL", "-4") // TODO: "debug" string should be acceptable
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_PATH", path.Join(tempDir, "e2engine.log"))
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_FORMAT", "text")
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_QUEUE_SIZE", "100")
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_BUFFER_SIZE", "1024")
				t.Setenv("E2ENGINE_LOGGER_SINKS_0_FLUSH_INTERVAL", "10ms")
			}

			getTestExecution := func(t *testing.T) model.TestExecution {
				return util.Load[model.TestExecution](
					t,
					path.Join(dataDir, test.testExecSpec),
				)
			}

			getTestExecutionRunning := func(t *testing.T) model.TestExecution {
				return util.Load[model.TestExecution](
					t,
					path.Join(dataDir, test.testExecRunningSpec),
				)
			}

			h := harness.New(t, "../bin/cli", test.configPath, path.Join(t.TempDir(), "data"))
			if len(test.services) > 0 {
				serviceCtx, serviceCancel := context.WithCancel(context.Background())
				t.Cleanup(serviceCancel)

				for _, service := range test.services {
					s := h.StartService(serviceCtx, service.Name, service.Args...)
					t.Cleanup(func() {
						s.Stop(t)
					})
				}
			}

			// create environment
			environmentID := h.CreateEnvironment(
				ctx,
				test.envSpec,
				path.Join(dataDir, test.envSpec+".env.yml"),
			)

			// create test
			testID := h.CreateTest(
				ctx,
				test.testSpec,
				path.Join(dataDir, test.testSpec+".test.yml"),
			)

			// execute test
			testExecutionID := h.RunTest(
				ctx,
				testID,
				environmentID,
			)

			// verify running state for long asynchronous execution
			if test.long {
				actualTestExecution := h.GetTestExecution(
					ctx,
					testExecutionID,
				)

				expectedTestExecution := getTestExecutionRunning(t)
				expectedTestExecution.ID = testExecutionID
				expectedTestExecution.TestID = testID
				expectedTestExecution.EnvironmentID = environmentID
				expectedTestExecution.Status = model.ExecutionStatusRunning
				expectedTestExecution.StartedAt = actualTestExecution.StartedAt
				expectedTestExecution.FinishedAt = time.Time{}

				util.AssertEqualTestExecutions(
					t,
					&expectedTestExecution,
					&actualTestExecution,
				)

				if test.crashRunner {
					logPath := path.Join(tempDir, "e2engine.log")
					pid := getRunnerPid(t, ctx, logPath)

					process, err := os.FindProcess(pid)
					if err != nil {
						t.Fatalf("Failed to find runner process %d: %v", pid, err)
					}

					if err := process.Kill(); err != nil {
						t.Fatalf("Failed to kill runner process %d: %v", pid, err)
					}

					actual := h.GetTestExecution(ctx, testExecutionID)

					if actual.Status != model.ExecutionStatusRunning {
						t.Fatalf(
							"Expected execution to remain running after runner crash, got %s",
							actual.Status,
						)
					}
				}
			}

			if !test.allowExpiry {
				// wait for terminal state
				actualTestExecution2 := h.WaitTestExecution(
					ctx,
					testExecutionID,
				)

				expectedTestExecution2 := getTestExecution(t)
				expectedTestExecution2.ID = testExecutionID
				expectedTestExecution2.TestID = testID
				expectedTestExecution2.EnvironmentID = environmentID
				expectedTestExecution2.StartedAt = actualTestExecution2.StartedAt
				expectedTestExecution2.FinishedAt = actualTestExecution2.FinishedAt

				util.AssertEqualTestExecutions(
					t,
					&expectedTestExecution2,
					&actualTestExecution2,
				)
			}
		})
	}
}

func getRunnerPid(
	t *testing.T,
	ctx context.Context,
	logPath string,
) int {
	t.Helper()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		if pid, ok := readRunnerPid(t, logPath); ok {
			return pid
		}

		select {
		case <-ctx.Done():
			t.Fatalf("Runner PID not found in log file %s", logPath)

		case <-ticker.C:
		}
	}
}

func readRunnerPid(t *testing.T, logPath string) (int, bool) {
	t.Helper()

	f, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false
		}

		t.Fatalf("Failed to open log file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, `Started internal-runner`) {
			continue
		}

		const prefix = `pid="`
		start := strings.Index(line, prefix)
		if start == -1 {
			continue
		}
		start += len(prefix)

		end := strings.IndexByte(line[start:], '"')
		if end == -1 {
			continue
		}

		pid, err := strconv.Atoi(line[start : start+end])
		if err != nil {
			t.Fatalf("Failed to parse runner PID: %v", err)
		}

		return pid, true
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	return 0, false
}

func TestDirectTransportWaitsForExecution(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTestTimeout,
	)
	defer cancel()

	t.Setenv("E2ENGINE_TRANSPORT_KIND", "direct")

	dataDir := "../data/transport-kind"

	h := harness.New(
		t,
		"../bin/cli",
		"",
		path.Join(t.TempDir(), "data"),
	)

	serviceCtx, serviceCancel := context.WithCancel(context.Background())
	t.Cleanup(serviceCancel)

	service := h.StartService(
		serviceCtx,
		"gateway-http-path",
		"2",
	)
	t.Cleanup(func() {
		service.Stop(t)
	})

	environmentID := h.CreateEnvironment(
		ctx,
		"real-with-mocked-dependency",
		path.Join(
			dataDir,
			"real-with-mocked-dependency.env.yml",
		),
	)

	testID := h.CreateTest(
		ctx,
		"gateway-calls-entitlements",
		path.Join(
			dataDir,
			"gateway-calls-entitlements.test.yml",
		),
	)

	executionID := h.RunTest(
		ctx,
		testID,
		environmentID,
	)

	actual := h.GetTestExecution(
		ctx,
		executionID,
	)

	if actual.Status != model.ExecutionStatusPassed {
		t.Fatalf(
			"Expected execution to be completed when direct run returns, got %s",
			actual.Status,
		)
	}

	if actual.FinishedAt.IsZero() {
		t.Fatal("Expected completed execution to have FinishedAt")
	}

	if actual.Summary == nil {
		t.Fatal("Expected completed execution to have summary")
	}
}

// Test verifies that detached runner #1 actually terminates and
// releases its listening socket before runner #2 starts.
func TestSocketTransportSequentialExecutions(t *testing.T) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultTestTimeout,
	)
	defer cancel()

	t.Setenv("E2ENGINE_TRANSPORT_KIND", "socket")
	t.Setenv("E2ENGINE_TRANSPORT_ADDRESS", "127.0.0.1:19002")

	dataDir := "../data/transport-kind"

	h := harness.New(
		t,
		"../bin/cli",
		"",
		path.Join(t.TempDir(), "data"),
	)

	serviceCtx, serviceCancel := context.WithCancel(context.Background())
	t.Cleanup(serviceCancel)

	service := h.StartService(serviceCtx, "gateway-http-path")
	t.Cleanup(func() {
		service.Stop(t)
	})

	environmentID := h.CreateEnvironment(
		ctx,
		"real-with-mocked-dependency",
		path.Join(
			dataDir,
			"real-with-mocked-dependency.env.yml",
		),
	)

	testID := h.CreateTest(
		ctx,
		"gateway-calls-entitlements",
		path.Join(
			dataDir,
			"gateway-calls-entitlements.test.yml",
		),
	)

	for i := range 2 {
		executionID := h.RunTest(
			ctx,
			testID,
			environmentID,
		)

		actual := h.WaitTestExecution(
			ctx,
			executionID,
		)

		if actual.Status != model.ExecutionStatusPassed {
			t.Fatalf(
				"Execution %d: expected status %s, got %s",
				i+1,
				model.ExecutionStatusPassed,
				actual.Status,
			)
		}

		if actual.Summary == nil {
			t.Fatalf(
				"Execution %d: expected non-nil summary",
				i+1,
			)
		}
	}
}
