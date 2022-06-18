package libfuzzer

import (
	"testing"

	"code-intelligence.com/cifuzz/integration/runner/libfuzzer/testutils"
)

func TestIntegration_CasesWrittenToCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	testutils.TestWithAndWithoutMinijail(t, func(t *testing.T, disableMinijail bool) {
		test := testutils.NewLibfuzzerTest(t, "new_paths_fuzzer", disableMinijail)

		_, _, reports := test.Run(t)

		testutils.CheckReports(t, reports, &testutils.CheckReportOptions{
			NumFindings: 0,
		})

		test.RequireSeedCorpusNotEmpty(t)
	})
}
