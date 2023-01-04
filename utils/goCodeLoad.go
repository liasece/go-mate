package utils

import (
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
)

func LoadGoInterface(path string, interfaceName string) (gocoder.Type, error) {
	return cde.GetInterfaceFromSource(path, interfaceName)
}

func LoadGoType(path string, typeName string) (gocoder.Type, error) {
	return cde.GetTypeFromSource(path, typeName)
}

func LoadGoStruct(path string, structName string) (gocoder.Type, error) {
	t, err := cde.GetTypeFromSource(path, structName)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func LoadGoMethods(path string, structName string) ([]gocoder.Func, error) {
	return cde.GetMethodsFromSource(path, structName)
}
