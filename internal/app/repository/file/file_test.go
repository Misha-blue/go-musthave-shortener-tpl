package file_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fileName = "test.file"

func TestGetAll(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "empty file",
			text: "",
		},
		{
			name: "file with empty elements",
			text: ";",
		},
		{
			name: "file with 1 element without new line",
			text: "1;1",
		},
		{
			name: "file with 1 element with new line",
			text: "1;1\n",
		},
		{
			name: "file with 2 elements without new line",
			text: "item1;1\nitem2;2",
		},
		{
			name: "file with 2 elements with new line",
			text: "item1;1\nitem2;2\n",
		},
		{
			name: "file with 5 elements with new line",
			text: "item1;1\nitem2;2\nitem3;3\nitem4;4\nitem5;5\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := file.New(fileName)
			require.NoError(t, err)
			require.NoError(t, prepareFile(tt.text))

			info, err := file.GetAll()
			require.NoError(t, err)

			expected := getExpected(tt.text)
			assert.Equal(t, info, expected, "file should have expected content")
			require.NoError(t, os.Remove(fileName))
		})
	}
}

func prepareFile(fileText string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Write([]byte(fileText))
	return nil
}

func getExpected(s string) map[string]string {
	expected := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		items := strings.Split(line, ";")
		if len(items) > 1 {
			expected[items[0]] = items[1]
		}
	}
	return expected
}
