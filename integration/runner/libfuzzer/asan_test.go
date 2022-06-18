package libfuzzer

import (
	"testing"

	"code-intelligence.com/cifuzz/integration/runner/libfuzzer/testutils"
	"code-intelligence.com/cifuzz/pkg/report"
)

func TestIntegration_ASAN(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutils.TestWithAndWithoutMinijail(t, func(t *testing.T, disableMinijail bool) {
		test := testutils.NewLibfuzzerTest(t, "trigger_asan", disableMinijail)

		_, _, reports := test.Run(t)

		testutils.CheckReports(t, reports, &testutils.CheckReportOptions{
			ErrorType:   report.ErrorType_CRASH,
			SourceFile:  "trigger_asan.c",
			Details:     "heap-buffer-overflow",
			NumFindings: 1,
		})
	})
}
