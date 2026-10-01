package harness

import (
	"bytes"
	"context"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/e2engine/core/model"
	"gopkg.in/yaml.v3"
)

type Harness struct {
	t          *testing.T
	cliPath    string
	configPath string
	dataDir    string
}

func New(t *testing.T, cliPath, configPath, dataDir string) *Harness {
	if configPath == "" {
		configPath = "../config/empty-config.yml"
	}

	t.Setenv("E2ENGINE_CONFIG_PATH", configPath)

	return &Harness{
		t:          t,
		cliPath:    cliPath,
		configPath: configPath,
		dataDir:    dataDir,
	}
}

func (h *Harness) RunCommand(
	ctx context.Context,
	args ...string,
) (stdout, stderr string, err error) {
	h.t.Helper()

	if h.dataDir != "" {
		h.t.Setenv("E2ENGINE_DATA_DIR", h.dataDir)
	}

	cmd := exec.CommandContext(
		ctx,
		h.cliPath,
		args...,
	)

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
