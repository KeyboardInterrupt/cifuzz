package dialog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	input := []byte("\n")

	r, w, err := os.Pipe()
	require.NoError(t, err)

	_, err = w.Write(input)
	require.NoError(t, err)
	w.Close()

	items := map[string]string{
		"Item No1": "item1",
	}

	userInput, err := Select("Test", items, r)
	assert.NoError(t, err)
	assert.Equal(t, "item1", userInput)
}

func TestInputFilename(t *testing.T) {
	err := os.Setenv("USE_PTY", "true")
	assert.NoError(t, err)

	input := []byte("my input\n")
	r, w, err := os.Pipe()
	require.NoError(t, err)

	_, err = w.Write(input)
	require.NoError(t, err)
	w.Close()

	userInput, err := InputFilename(r, "Test", "default")
	assert.NoError(t, err)
	assert.Equal(t, "my input", userInput)
}

// Should return default value if user just presses "enter"
func TestInputFilename_Default(t *testing.T) {
	err := os.Setenv("USE_PTY", "true")
	assert.NoError(t, err)

	input := []byte("\n")
	r, w, err := os.Pipe()
	require.NoError(t, err)

	_, err = w.Write(input)
	require.NoError(t, err)
	w.Close()

	userInput, err := InputFilename(r, "Test", "default")
	assert.NoError(t, err)
	assert.Equal(t, "default", userInput)
}
