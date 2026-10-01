package harness

import (
	"context"
	"time"

	"github.com/e2engine/core/model"
)

func (h *Harness) WaitTestExecution(
	ctx context.Context,
	testExecutionID string,
) model.TestExecution {
	h.t.Helper()

	actualTestExecution := h.GetTestExecution(
		ctx,
		testExecutionID,
	)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for actualTestExecution.Status == model.ExecutionStatusRunning {
		select {
		case <-ctx.Done():
			h.t.Fatalf(
				"Timed out waiting for test execution %s to complete",
				testExecutionID,
			)

		case <-ticker.C:
		}

		actualTestExecution = h.GetTestExecution(
			ctx,
			testExecutionID,
		)
	}

	return actualTestExecution
}
