package goinstall

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (i *installer) AddRemakeDirStep(path string) {
	rmkDir := newRemakeDirFunc(os.RemoveAll, os.MkdirAll)
	process := func() error {
		return rmkDir(path)
	}
	i.AddStep(process, i.printer.Sprintf(remakeDirMsg, path))
}

func (i *installer) AddCopyFilesIntoDirStep(path string, files map[string][]byte) {
	copyFile := newCopyFileFunc(ioutil.WriteFile)
	process := func() error {
		for fileName, file := range files {
			if err := copyFile(path, file, fileName); err != nil {
				return err
			}
		}
		return nil
	}
	i.AddStep(process, i.printer.Sprintf(copyingFilesMsg, path))
}

//remakeDir removes given directory if it exists and recreates it
type remakeDir = func(path string) error
type removeDir = func(path string) error
type createDir = func(path string, perm fs.FileMode) error

func newRemakeDirFunc(removeDir removeDir, createDir createDir) remakeDir {
	return func(path string) error {
		err := removeDir(path)
		if err != nil {
			return err
		}
		return createDir(path, os.ModePerm)
	}
}

type copyFile = func(dirPath string, file []byte, name string) error
type writeFile = func(filePath string, file []byte, perm fs.FileMode) error

func newCopyFileFunc(writeFile writeFile) copyFile {
	return func(dirPath string, file []byte, name string) error {
		filePath := filepath.Join(dirPath, name)
		return writeFile(filePath, file, os.ModePerm)
	}
}
