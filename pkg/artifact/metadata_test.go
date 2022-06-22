package artifact

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToYaml_Minimal(t *testing.T) {
	artifact := &Metadata{
		RunEnvironment: &RunEnvironment{
			Docker: "my-image",
		},
	}

	yaml, err := artifact.ToYaml()
	require.NoError(t, err)
	require.NotEmpty(t, yaml)

	assert.Contains(t, string(yaml), "docker: my-image")
}
