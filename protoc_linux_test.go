//+build linux

package goinstall

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestExecXDGMime(t *testing.T) {
	cmdRun := func(name string, args ...string) error {return errors.New("failed to run command")}
	execXDGMime := newExecXDGMimeFunc(cmdRun, dotDesktopFile{})
	err := execXDGMime()
	assert.Error(t, errors.New("failed to run command"), err)
}

func TestNewDotDesktopFile(t *testing.T) {
	df := newDotDesktopFile(make([]byte, 5), "test")
	assert.Equal(t, "test.desktop", df.name)
	assert.Equal(t, make([]byte, 5), df.bytes)
}

func TestCopyDotDesktopFileToUserAppDir_Integration(t *testing.T) {
	//Integration test done 4/07/2021 on windows
	t.Skip()
	copyFile := newCopyFileFunc(ioutil.WriteFile)
	getHomeDir := os.UserHomeDir
	copyDotDesktop := newCopyDotDesktopFileToUserAppDirFunc(copyFile, dotDesktopFile{
		bytes: make([]byte, 5),
		name: "test.desktop",
	}, getHomeDir)
	assert.NoError(t, copyDotDesktop())
}

func TestCopyDotDesktopFileToUserAppDir(t *testing.T) {
	tests := []struct{
		description string

		copyFile copyFile
		getHomeDir getHomeDir
		dotDesktopFile dotDesktopFile
		outErr error
	}{
		{
			description: "Getting home dir should fail",
			getHomeDir: func() (string, error) {
				return "", errors.New("failed to get home dir")
			},
			outErr: errors.New("failed to get home dir"),
		},
		{
			description: "Copying desktop file should fail",
			getHomeDir: func() (string, error) {
				return "", nil
			},
			copyFile: func(dirPath string, file []byte, name string) error {
				return errors.New("failed to copy")
			},
			outErr: errors.New("failed to copy"),
		},
	}

	for _, test := range tests {
		copy := newCopyDotDesktopFileToUserAppDirFunc(test.copyFile, test.dotDesktopFile, test.getHomeDir)
		err := copy()
		assert.Equal(t, test.outErr, err)
	}
}
