package file

import (
	"fmt"
	"os"
)

// MustOpen :
func MustOpen(filelName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permimssion denied src: %s", src)
	}

	err = IsNotExistMkDir(dir + "/wwwroot/")
	err = IsNotExistMkDir(dir + "/wwwroot/runtime/")
	err = IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+filelName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

// Open :
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// IsNotExistMkDir :
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}

	return nil
}

// MkDir :
func MkDir(src string) error {
	err := os.Mkdir(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// CheckNotExist :
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// CheckPermission :
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}
