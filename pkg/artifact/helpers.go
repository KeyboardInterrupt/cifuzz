package artifact

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"

	"code-intelligence.com/cifuzz/pkg/log"
	"code-intelligence.com/cifuzz/pkg/vcs"
)

func getCodeRevision() (*CodeRevision, error) {
	_, err := exec.LookPath("git")
	if errors.Is(err, exec.ErrNotFound) {
		log.Warn("git not found")
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	branch, err := vcs.GetGitBranch()
	if err != nil {
		return nil, err
	}

	commit, err := vcs.GetGitCommit()
	if err != nil {
		return nil, err
	}

	return &CodeRevision{
		Git: &GitRevision{
			Branch: branch,
			Commit: commit,
		},
	}, nil
}

func getCMakeFuzzers() ([]*Fuzzer, error) {
	matches, err := zglob.Glob(".cifuzz-build/**/.cifuzz/fuzz_tests/*")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fuzzers := []*Fuzzer{}
	for _, match := range matches {

		// TODO replace this with the "correct" way (CMake Cache Variables maybe)
		// parse engine & sanitizer out of fuzz target path
		re := regexp.MustCompile(`\/(?P<engine>[[:alpha:]]+)\/(?P<sanitizer>[[:alpha:]+]+)\/`)
		groups := re.FindStringSubmatch(match)
		engine := strings.ToUpper(groups[1])
		// TODO according to the artifact spec there should be only one sanitizer per fuzzer
		//sanitizer :=  strings.Split(groups[2], "+")

		fuzzer := &Fuzzer{
			Target: filepath.Base(match),
			Path:   match,
			Engine: engine,
		}
		fuzzers = append(fuzzers, fuzzer)
	}

	return fuzzers, nil
}
