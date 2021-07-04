package goinstall

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemakeDir(t *testing.T) {
	tests := []struct {
		description string

		removeDir removeDir
		createDir createDir

		path   string
		outErr error
	}{
		{
			description: "should fail to remove dir",

			removeDir: func(path string) error {
				return errors.New("failed to remove dir")
			},
			outErr: errors.New("failed to remove dir"),
		},
		{
			description: "should fail to createDir",

			removeDir: func(path string) error { return nil },
			createDir: func(path string, perm fs.FileMode) error {
				return errors.New("fail to create dir")
			},
			outErr: errors.New("fail to create dir"),
		},
	}

	for _, test := range tests {
		rmkDir := newRemakeDirFunc(test.removeDir, test.createDir)
		err := rmkDir(test.path)
		assert.Equal(t, test.outErr, err)
	}
}

func TestCopyFile(t *testing.T) {
	tests := []struct {
		description string

		writeFile writeFile

		dirPath string
		file    []byte
		name    string

		outErr error
	}{
		{
			description: "should fail to write file",
			writeFile: func(filePath string, file []byte, perm fs.FileMode) error {
				return errors.New("failed to write")
			},
			outErr: errors.New("failed to write"),
		},
	}

	for _, test := range tests {
		copyFile := newCopyFileFunc(test.writeFile)
		err := copyFile(test.dirPath, test.file, test.name)
		assert.Equal(t, test.outErr, err)
	}
}

func TestRemakeDirIntegration(t *testing.T) {
	//integration test done 4/07/2021
	t.Skip()
	tests := []struct {
		path   string
		outErr error
	}{
		{
			path:   "tests/testdir",
			outErr: nil,
		},
	}

	for _, test := range tests {
		rmkDir := newRemakeDirFunc(os.RemoveAll, os.MkdirAll)
		err := rmkDir(test.path)
		assert.Equal(t, test.outErr, err)
	}
}

func TestCopyFileIntegration(t *testing.T) {
	//integration test done 4/07/2021
	t.Skip()
	tests := []struct {
		dirPath string
		file    []byte
		name    string

		outErr error
	}{
		{
			dirPath: "tests",
			file:    make([]byte, 5),
			name:    "test.txt",

			outErr: nil,
		},
	}

	for _, test := range tests {
		copy := newCopyFileFunc(ioutil.WriteFile)
		err := copy(test.dirPath, test.file, test.name)
		assert.Equal(t, test.outErr, err)
	}
}


