package jsonoutput

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewManagerFromFile(t *testing.T) {
	previewManager, err := NewManagerFromFile("testdata/preview-changes.json")
	require.NoError(t, err)

	t.Log(previewManager.ShortSummaryString())

	t.Log(previewManager.TreeString())

	errorManager, err := NewManagerFromFile("testdata/error.json")
	require.NoError(t, err)

	t.Log(errorManager.ShortSummaryString())
}
