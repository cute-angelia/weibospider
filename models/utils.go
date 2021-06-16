package models

import "os"

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(dir, os.ModeDir|0755)
		} else {
			return err
		}
	}
	return nil
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
