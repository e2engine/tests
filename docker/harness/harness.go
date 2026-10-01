package harness

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/e2engine/core/model"
	"gopkg.in/yaml.v3"
)

type Harness struct {
	t          *testing.T
	configPath string
	dataDir    string
	dataVolume string
	network    string
}

func New(t *testing.T, configPath, network string) *Harness {
	t.Helper()

	if configPath == "" {
		configPath = "../config/empty-config.yml"
	}

	configPath, err := filepath.Abs(configPath)
	if err != nil {
		t.Fatalf("Failed to resolve config path %q: %v", configPath, err)
	}

	dataDir, err := filepath.Abs("../data")
	if err != nil {
		t.Fatalf("Failed to resolve data directory %q: %v", dataDir, err)
	}

	volumeCmd := exec.Command("docker", "volume", "create")
	var volumeStdoutBuf bytes.Buffer
	var volumeStderrBuf bytes.Buffer
	volumeCmd.Stdout = &volumeStdoutBuf
	volumeCmd.Stderr = &volumeStderrBuf

	if err = volumeCmd.Run(); err != nil {
		t.Fatalf(
			"Failed to create Docker data volume: %v, stderr: %s",
			err,
			volumeStderrBuf.String(),
		)
	}

	dataVolume := strings.TrimSpace(volumeStdoutBuf.String())
	if dataVolume == "" {
		t.Fatal("Docker returned an empty data volume name")
	}

	t.Cleanup(func() {
		cmdCleanup := exec.Command(
			"docker",
			"volume",
			"rm",
			"-f",
			dataVolume,
		)

		var stderrBuf bytes.Buffer
		cmdCleanup.Stderr = &stderrBuf

		if errCleanup := cmdCleanup.Run(); errCleanup != nil {
			t.Errorf(
				"Failed to remove Docker data volume %q: %v, stderr: %s",
				dataVolume,
				errCleanup,
				stderrBuf.String(),
			)
		}
	})

	networkCmd := exec.Command(
		"docker",
		"network",
		"create",
		network,
	)

	var networkStdoutBuf bytes.Buffer
	var networkStderrBuf bytes.Buffer
	networkCmd.Stdout = &networkStdoutBuf
	networkCmd.Stderr = &networkStderrBuf

	if err = networkCmd.Run(); err != nil {
		t.Fatalf(
			"Failed to create Docker network %q: %v, stderr: %s",
			network,
			err,
			networkStderrBuf.String(),
		)
	}

	t.Cleanup(func() {
		cmdCleanup := exec.Command(
			"docker",
			"network",
			"rm",
			network,
		)

		var stderrBuf bytes.Buffer
		cmdCleanup.Stderr = &stderrBuf

		if err := cmdCleanup.Run(); err != nil {
			t.Errorf(
				"Failed to remove Docker network %q: %v, stderr: %s",
				network,
				err,
				stderrBuf.String(),
			)
		}
	})

	return &Harness{
		t:          t,
		configPath: configPath,
		dataDir:    dataDir,
		dataVolume: dataVolume,
		network:    network,
	}
}

func (h *Harness) RunCommand(
	ctx context.Context,
	commandArgs ...string,
) (stdout, stderr string, err error) {
	h.t.Helper()

	containerConfigDir := "/config"
	containerConfigPath := containerConfigDir + "/" + filepath.Base(h.configPath)

	args := []string{
		"run",
		"--rm",
		"--network", h.network,
		"-v", filepath.Dir(h.configPath) + ":" + containerConfigDir + ":ro",
		"-v", h.dataDir + ":/data",
		"-v", h.dataVolume + ":/testdata",
		"-e", "E2ENGINE_CONFIG_PATH=" + containerConfigPath,
		"-e", "E2ENGINE_DATA_DIR=/testdata",
	}

	args = append(args, "e2engine:local")
	args = append(args, commandArgs...)

	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	return stdoutBuf.String(), stderrBuf.String(), err
}

func (h *Harness) ListEnvironments(ctx context.Context) []model.Environment {
	h.t.Helper()

	return list[model.Environment](ctx, h, "environment")
}

func (h *Harness) ListTests(ctx context.Context) []model.Test {
	h.t.Helper()

	return list[model.Test](ctx, h, "test")
}

func (h *Harness) ListTestSuites(ctx context.Context) []model.TestSuite {
	h.t.Helper()

	return list[model.TestSuite](ctx, h, "testsuite")
}

func (h *Harness) ListTestExecutions(ctx context.Context) []model.TestExecution {
	h.t.Helper()

	return list[model.TestExecution](ctx, h, "testexecution")
}

func (h *Harness) ListTestSuiteExecutions(ctx context.Context) []model.TestSuiteExecution {
	h.t.Helper()

	return list[model.TestSuiteExecution](ctx, h, "testsuiteexecution")
}

func (h *Harness) CreateEnvironment(ctx context.Context, name, path string) string {
	h.t.Helper()

	return create(ctx, h, "environment", name, path)
}

func (h *Harness) CreateTest(ctx context.Context, name, path string) string {
	h.t.Helper()

	return create(ctx, h, "test", name, path)
}

func (h *Harness) CreateTestSuite(ctx context.Context, name, path string) string {
	h.t.Helper()

	return create(ctx, h, "testsuite", name, path)
}

func (h *Harness) GetEnvironment(ctx context.Context, ref string) model.Environment {
	h.t.Helper()

	return get[model.Environment](ctx, h, ref, "environment")
}

func (h *Harness) GetTest(ctx context.Context, ref string) model.Test {
	h.t.Helper()

	return get[model.Test](ctx, h, ref, "test")
}

func (h *Harness) GetTestSuite(ctx context.Context, ref string) model.TestSuite {
	h.t.Helper()

	return get[model.TestSuite](ctx, h, ref, "testsuite")
}

func (h *Harness) GetTestExecution(ctx context.Context, ref string) model.TestExecution {
	h.t.Helper()

	return get[model.TestExecution](ctx, h, ref, "testexecution")
}

func (h *Harness) GetTestSuiteExecution(ctx context.Context, ref string) model.TestSuiteExecution {
	h.t.Helper()

	return get[model.TestSuiteExecution](ctx, h, ref, "testsuiteexecution")
}

func (h *Harness) RunTest(ctx context.Context, testRef, envRef string) string {
	h.t.Helper()

	return run(ctx, h, "test", testRef, envRef)
}

func (h *Harness) RunTestSuite(ctx context.Context, tsRef, envRef string) string {
	h.t.Helper()

	return run(ctx, h, "testsuite", tsRef, envRef)
}

func (h *Harness) DeleteEnvironment(ctx context.Context, ref, name string) {
	h.t.Helper()

	del(ctx, h, "environment", ref, name)
}

func (h *Harness) DeleteTest(ctx context.Context, ref, name string) {
	h.t.Helper()

	del(ctx, h, "test", ref, name)
}

func (h *Harness) DeleteTestSuite(ctx context.Context, ref, name string) {
	h.t.Helper()

	del(ctx, h, "testsuite", ref, name)
}

func (h *Harness) DeleteTestExecution(ctx context.Context, ref, name string) {
	h.t.Helper()

	del(ctx, h, "testexecution", ref, name)
}

func (h *Harness) DeleteTestSuiteExecution(ctx context.Context, ref, name string) {
	h.t.Helper()

	del(ctx, h, "testsuiteexecution", ref, name)
}

func list[T any](ctx context.Context, h *Harness, kind string) []T {
	h.t.Helper()

	stdout, stderr, err := h.RunCommand(
		ctx,
		"get",
		kind+"s",
		"-v",
		"-o",
		"yaml",
	)
	if err != nil {
		h.t.Fatalf("Unexpected RunCommand error: %v", err)
	}

	if strings.TrimSpace(stderr) != "" {
		h.t.Errorf("Expected no stderr output, got: %s", stderr)
	}

	var results []T
	if errUnmarshal := yaml.Unmarshal([]byte(stdout), &results); errUnmarshal != nil {
		h.t.Fatalf("Failed to unmarshal %ss output: %v\nstdout: %s\nstderr: %s", kind, errUnmarshal, stdout, stderr)
	}

	return results
}

func create(ctx context.Context, h *Harness, kind, name, path string) string {
	h.t.Helper()

	stdout, stderr, err := h.RunCommand(
		ctx,
		"create",
		kind,
		path,
	)
	if err != nil {
		h.t.Fatalf("Unexpected RunCommand error: %v, stderr: %s", err, stderr)
	}

	if strings.TrimSpace(stderr) != "" {
		h.t.Errorf("Expected no stderr output, got: %s", stderr)
	}

	regex := `^created ` + kind + `, name: ` + name + ` id: ([a-f0-9]{64}) version: 1\.0\.0$`

	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(strings.TrimSpace(stdout))
	if matches == nil {
		h.t.Fatalf(
			"Expected stdout to match regex: %q, got: %q",
			regex,
			stdout,
		)
	}

	if strings.TrimSpace(stdout) != matches[0] {
		h.t.Fatalf(
			"Expected stdout to be: %q, got: %q",
			matches[0],
			stdout,
		)
	}

	return matches[1] // id
}

func get[T any](ctx context.Context, h *Harness, ref, kind string) T {
	h.t.Helper()

	stdout, stderr, err := h.RunCommand(
		ctx,
		"get",
		kind,
		ref,
		"-v",
		"-o",
		"yaml",
	)
	if err != nil {
		h.t.Fatalf("Unexpected RunCommand error: %v, stderr: %s", err, stderr)
	}

	if strings.TrimSpace(stderr) != "" {
		h.t.Errorf("Expected no stderr output, got: %s", stderr)
	}

	var result T
	if errUnmarshal := yaml.Unmarshal([]byte(stdout), &result); errUnmarshal != nil {
		h.t.Fatalf(
			"Failed to unmarshal %s output: %v\nstdout: %s\nstderr: %s",
			kind,
			errUnmarshal,
			stdout,
			stderr,
		)
	}

	return result
}

func run(ctx context.Context, h *Harness, kind, testRef, envRef string) string {
	h.t.Helper()

	stdout, stderr, err := h.RunCommand(
		ctx,
		"run",
		kind,
		testRef,
		envRef,
	)
	if err != nil {
		h.t.Fatalf("Unexpected RunCommand error: %v, stderr: %s", err, stderr)
	}

	if strings.TrimSpace(stderr) != "" {
		h.t.Errorf("Expected no stderr output, got: %s", stderr)
	}

	regex := `^created ` + kind + ` execution with id: ([a-f0-9]{64})$`

	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(strings.TrimSpace(stdout))
	if matches == nil {
		h.t.Fatalf(
			"Expected stdout to match regex: %q, got: %q",
			regex,
			stdout,
		)
	}

	if strings.TrimSpace(stdout) != matches[0] {
		h.t.Fatalf(
			"Expected stdout to be: %q, got: %q",
			matches[0],
			stdout,
		)
	}

	return matches[1] // id
}

func del(ctx context.Context, h *Harness, kind, ref, name string) {
	h.t.Helper()

	stdout, stderr, err := h.RunCommand(
		ctx,
		"delete",
		kind,
		ref,
	)
	if err != nil {
		h.t.Fatalf("Unexpected RunCommand error: %v", err)
	}

	if strings.TrimSpace(stderr) != "" {
		h.t.Errorf("Expected no stderr output, got: %s", stderr)
	}

	expectedStdout := "deleted " + kind + ", name: " + name + " id: " + ref + " version: 1.0.0"
	switch kind {
	case "testexecution":
		expectedStdout = "deleted test execution with id: " + ref
	case "testsuiteexecution":
		expectedStdout = "deleted testsuite execution with id: " + ref
	}

	if strings.TrimSpace(stdout) != expectedStdout {
		h.t.Fatalf(
			"Expected stdout to be: %q, got: %q",
			expectedStdout,
			stdout,
		)
	}
}
