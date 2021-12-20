package repository_test

import (
	"os"
	"testing"

	"github.com/Misha-blue/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fileName = "test.file"

func TestGetAll(t *testing.T) {
	type want struct {
		err   string
		value string
	}

	tests := []struct {
		name string
		text string
		item string
		want want
	}{
		{
			name: "file with 2 elements without new line",
			text: "item1;1\nitem2;2",
			item: "item1",
			want: want{
				value: "1",
				err:   "",
			},
		},
		{
			name: "file with 5 elements with new line",
			text: "item1;1\nitem2;2\nitem3;3\nitem4;4\nitem5;5\n",
			item: "item2",
			want: want{
				value: "2",
				err:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := repository.New(fileName)
			require.NoError(t, err)
			require.NoError(t, prepareFile(tt.text))

			actual, err := repo.Load(tt.item)
			require.NoError(t, err)

			assert.Equal(t, tt.want.value, actual, "file should have expected content")
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
