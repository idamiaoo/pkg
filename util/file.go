package util

import "os"

var File FAlias

type FAlias os.File

func (fa *FAlias) IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		return false
	}

	return true
}

func (fa *FAlias) CreateFile(name string) (*os.File, error) {
	if fa.IsExists(name) {
		f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND, os.ModeAppend)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	f, err := os.Create(name)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (fa *FAlias) WriteFile(file string, content []byte) error {
	f, err := fa.CreateFile(file)
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return f.Close()
}

func (fa *FAlias) DeleteFile(fileName string) error {
	if !fa.IsExists(fileName) {
		return nil
	}

	if err := os.Remove(fileName); err != nil {
		return err
	}

	return nil
}
