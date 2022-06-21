package vcs

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	viper.Set("verbose", true)
	m.Run()
}

func TestGitBranch(t *testing.T) {
	branch, err := GetGitBranch()
	require.NoError(t, err)
	assert.NotEmpty(t, branch)
}

func TestGitCommit(t *testing.T) {
	commit, err := GetGitCommit()
	require.NoError(t, err)
	assert.Len(t, commit, 40)
}
