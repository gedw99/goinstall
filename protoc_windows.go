//+build windows

package goinstall

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
)


//AddRegisterProtocolOnWindowsStep registers a new protocol that will trigger a given bash command string
func (i *installer) AddRegisterProtocolOnWindowsStep(protocol string, execCmd string) {
	setRegistryKeyValues := newSetRegistryKeyValuesFunc()
	registerRegistryKeyWithValues := newRegisterRegistryKeyWithValuesFunc(registry.CreateKey, setRegistryKeyValues)
	registerProtocolOnWindows := newRegisterProtocolOnWindowsFunc(protocol, execCmd, registerRegistryKeyWithValues)
	i.AddStep(registerProtocolOnWindows, i.printer.Sprintf(registerProtocolMsg, protocol))
}

type createRegistryKey = func(k registry.Key, path string, access uint32) (newK registry.Key, opened bool, err error)
type registerProtocolOnWindows = func() error

func newRegisterProtocolOnWindowsFunc(protocol, execCmd string, registerKeyWithValues registerRegistryKeyWithValues) registerProtocolOnWindows {
	const (
		software       = "SOFTWARE"
		classes        = "Classes"
		shell          = "shell"
		open           = "open"
		command        = "command"
		defaultKey     = ""
		urlProtocolKey = "URL Protocol"
	)
	protocolPath := filepath.Join(software, classes, protocol)
	urlProtocolValue := map[string]string{urlProtocolKey: defaultKey}
	cmdPath := filepath.Join(shell, open, command)
	cmdValues := map[string]string{defaultKey: execCmd}

	return func() error {
		protocolKey, err := registerKeyWithValues(registry.CURRENT_USER, protocolPath, urlProtocolValue)
		if err != nil {
			return fmt.Errorf("unable to install protocol key : %v", err)
		}
		defer protocolKey.Close()

		cmdKey, err := registerKeyWithValues(protocolKey, cmdPath, cmdValues)
		if err != nil {
			return fmt.Errorf("unable to install command key : %v", err)
		}
		return cmdKey.Close()
	}
}

type registerRegistryKeyWithValues = func(currentKey registry.Key, newPath string, values map[string]string) (registry.Key, error)

func newRegisterRegistryKeyWithValuesFunc(createKey createRegistryKey, setValues setRegistryKeyValues) registerRegistryKeyWithValues {
	return func(currentKey registry.Key, newPath string, values map[string]string) (registry.Key, error) {
		newK, _, err := createKey(currentKey, newPath, registry.ALL_ACCESS)
		if err != nil {
			return newK, err
		}
		return newK, setValues(newK, values)
	}
}

type setRegistryKeyValues = func(key keyEditer, values map[string]string) error
func newSetRegistryKeyValuesFunc() setRegistryKeyValues {
	return func(key keyEditer, values map[string]string) error {
		for name, val := range values {
			err := key.SetStringValue(name, val)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

type keyEditer interface {
	SetStringValue(name, val string) error
}


