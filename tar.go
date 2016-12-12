package gotar

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"fmt"
)

type TarOpts struct {
	InPaths []string
	WorkDir string
	Writer  io.Writer
}

func addFile(tarWriter *tar.Writer, tarHeader *tar.Header, path string) error {
	err := tarWriter.WriteHeader(tarHeader)
	if err != nil {
		return err
	}
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = io.Copy(tarWriter, fp)
	if err != nil {
		return err
	}
	return nil
}

func Tar(to *TarOpts) error {
	// change working directory
	if (to.WorkDir != ""){
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Chdir(to.WorkDir)
		if err != nil {
			return err
		}
		defer os.Chdir(pwd)
	}

	// set the writer
	tarWriter := tar.NewWriter(to.Writer)

	// walk the file tree
	for _, path := range to.InPaths {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				tarHeader, err := tar.FileInfoHeader(info, info.Name())
				if err != nil {
					return err
				}
				tarHeader.Name = filepath.Join(filepath.Dir(path), tarHeader.Name)
				if info.IsDir() {
					err := tarWriter.WriteHeader(tarHeader)
					if err != nil {
						return err
					}
				} else {
					err = addFile(tarWriter, tarHeader, path)
					if err != nil {
						return err
					}
				}
				return nil
			})
		} else {
			fmt.Println(path)
			tarHeader, err := tar.FileInfoHeader(info, info.Name())
			tarHeader.Name = filepath.Join(filepath.Dir(path), tarHeader.Name)
			if err != nil {
				return err
			}
			err = addFile(tarWriter, tarHeader, path)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
