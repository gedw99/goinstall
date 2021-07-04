//+build windows

package goinstall

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows/registry"
	"testing"
)

func TestRegisterNewRegistryKeyWithValues(t *testing.T) {
	tests := []struct{
		description string

		createKey createRegistryKey
		setValues setRegistryKeyValues

		currentKey registry.Key
		newPath string
		values map[string]string

		outKey registry.Key
		outErr error
	}{
		{
			description: "should fail to create Key",
			createKey: func(k registry.Key, path string, access uint32) (newK registry.Key, opened bool, err error) {
				return registry.CLASSES_ROOT, false, errors.New("failed to create key")
			},
			outKey: registry.CLASSES_ROOT,
			outErr: errors.New("failed to create key"),
		},
		{
			description: "should fail to set values",
			createKey: func(k registry.Key, path string, access uint32) (newK registry.Key, opened bool, err error) {
				return registry.CLASSES_ROOT, false, nil
			},
			setValues: func(key keyEditer, values map[string]string) error {
				return errors.New("fail setting values")
			},
			outKey: registry.CLASSES_ROOT,
			outErr: errors.New("fail setting values"),
		},
	}
	for _, test := range tests {
		registerKey := newRegisterRegistryKeyWithValuesFunc(test.createKey, test.setValues)
		k, err := registerKey(test.currentKey,test.newPath, test.values)
		assert.Equal(t, test.outKey, k)
		assert.Equal(t, test.outErr, err)
	}
}

func TestSetRegistryKeyValues(t *testing.T) {
	tests := []struct{
		description string

		key keyEditer
		values map[string]string
		outErr error
	}{
		{
			description: "should fail to set string value",
			key: &mockKeyEditer{err: errors.New("fail to set string val")},
			values: map[string]string{"test":"test"},
			outErr: errors.New("fail to set string val"),
		},
	}
	for _, test := range tests {
		setRegistryKeyValues := newSetRegistryKeyValuesFunc()
		err := setRegistryKeyValues(test.key, test.values)
		assert.Equal(t, test.outErr, err)
	}
}

func TestRegisterProtocolOnWindows_Integration(t *testing.T) {
	//test done 04/07/2021 windows 10
	t.Skip()
	setRegistryKeyValues := newSetRegistryKeyValuesFunc()
	registerNewRegistryKeyWithValues := newRegisterRegistryKeyWithValuesFunc(registry.CreateKey, setRegistryKeyValues)
	registerProtocol := newRegisterProtocolOnWindowsFunc("test", "testCmd", registerNewRegistryKeyWithValues)
	assert.NoError(t, registerProtocol())
}

type mockKeyEditer struct {
	err error
}

func (k *mockKeyEditer) SetStringValue(name, val string) error {
	return k.err
}
