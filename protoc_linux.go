//+build linux

package goinstall

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//AddRegisterProtocolOnLinuxStep using xdg-mime app.
//A .desktop file has to be provided to handle registered protocol
func (i *installer) AddRegisterProtocolOnLinuxStep(protocol string, dotDesktopFileBytes []byte) {
	dotDesktopFile := newDotDesktopFile(dotDesktopFileBytes, protocol)
	copy := newCopyFileFunc(ioutil.WriteFile)

	copyDotDesktopFile := newCopyDotDesktopFileToUserAppDirFunc(copy, dotDesktopFile, os.UserHomeDir)
	runXDGMime := newExecXDGMimeFunc(newCmdRunFunc(), dotDesktopFile)

	process := newRegisterProtocolOnLinuxProcess(copyDotDesktopFile, runXDGMime)
	i.AddStep(process, i.printer.Sprintf(registerProtocolMsg, protocol))
}

func newRegisterProtocolOnLinuxProcess(copyDotDesktopFile func() error, runXDGMime func() error) func() error {
	return func() error {
		err := copyDotDesktopFile()
		if err != nil {
			return err
		}
		return runXDGMime()
	}
}

func newDotDesktopFile(fileBytes []byte, name string) dotDesktopFile {
	const desktopExtension = ".desktop"
	return dotDesktopFile{
		bytes: fileBytes,
		name: name + desktopExtension,
	}
}

type dotDesktopFile struct {
	bytes []byte
	name string
}

func newExecXDGMimeFunc(cmdRun cmdRun, file dotDesktopFile) func() error {
	const (
		xdgMimeProg = "xdg-mime"
		defaultXDGMimeOption = "default"
		schemeHandler = "x-scheme-handler"
		schemeHandlerSeparator = "/"
	)
	fileSchemeHandler := schemeHandler + schemeHandlerSeparator + file.name
	return func() error {
		return cmdRun(xdgMimeProg, defaultXDGMimeOption, file.name, fileSchemeHandler)
	}
}

type cmdRun = func(name string, args ...string) error
func newCmdRunFunc() cmdRun {
	return func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		return cmd.Run()
	}
}

func newCopyDotDesktopFileToUserAppDirFunc(copy copyFile, dotDesktopFile dotDesktopFile, getHomeDir getHomeDir) func() error {
	const shareAppPath = ".local/share/applications"
	return func() error {
		homeDir, err := getHomeDir()
		if err != nil {
			return err
		}
		return copy(filepath.Join(homeDir, shareAppPath), dotDesktopFile.bytes, dotDesktopFile.name)
	}
}



type getHomeDir = func() (string, error)