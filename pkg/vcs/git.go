package vcs

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"

	"code-intelligence.com/cifuzz/pkg/log"
	"code-intelligence.com/cifuzz/util/executil"
)

func runCommand(cmd *executil.Cmd) (string, error) {
	pipeR, pipeW, err := os.Pipe()
	if err != nil {
		return "", errors.WithStack(err)
	}
	cmd.Stdout = pipeW
	defer pipeR.Close()

	log.Debugf("Command: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		return "", err
	}
	pipeW.Close()

	out, err := ioutil.ReadAll(pipeR)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return strings.TrimSpace(string(out)), nil
}

func GetGitCommit() (string, error) {
	cmd := executil.Command("git", "rev-parse", "HEAD")
	// TODO handle invalid return value
	commit, err := runCommand(cmd)
	if err != nil {
		return "", err
	}
	log.Debugf("Fetched git commit: %s", commit)
	return commit, nil
}

func GetGitBranch() (string, error) {
	cmd := executil.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	// TODO handle invalid return value
	branch, err := runCommand(cmd)
	if err != nil {
		return "", err
	}
	log.Debugf("Fetched git branch: %s", branch)
	return branch, nil
}
