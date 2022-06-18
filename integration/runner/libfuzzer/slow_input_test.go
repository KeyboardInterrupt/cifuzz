package libfuzzer

import (
	"testing"

	"code-intelligence.com/cifuzz/integration/runner/libfuzzer/testutils"
	"code-intelligence.com/cifuzz/pkg/report"
)

func TestIntegration_SlowInput(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutils.TestWithAndWithoutMinijail(t, func(t *testing.T, disableMinijail bool) {
		test := testutils.NewLibfuzzerTest(t, "trigger_slow_input", disableMinijail)
		// The input timeout should be reported on the first input
		test.RunsLimit = 1
		test.EngineArgs = append(test.EngineArgs, "-report_slow_units=1")

		_, _, reports := test.Run(t)

		testutils.CheckReports(t, reports, &testutils.CheckReportOptions{
			ErrorType:   report.ErrorType_WARNING,
			Details:     "Slow input detected",
			NumFindings: 1,
		})
	})
}
