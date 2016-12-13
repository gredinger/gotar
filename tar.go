package gotar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type TarOpts struct {
	InPaths []string
	Gzip    bool
	Verbose bool
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

func Tar(opts *TarOpts) error {
	// change working directory
	if opts.WorkDir != "" {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = os.Chdir(opts.WorkDir)
		if err != nil {
			return err
		}
		defer os.Chdir(pwd)
	}

	// set the writer
	var tarWriter *tar.Writer
	var gzipWriter *gzip.Writer
	if opts.Gzip {
		gzipWriter = gzip.NewWriter(opts.Writer)
		tarWriter = tar.NewWriter(gzipWriter)
	} else {
		tarWriter = tar.NewWriter(opts.Writer)
	}
	defer func() {
		tarWriter.Close()
		if gzipWriter != nil {
			gzipWriter.Close()
		}
	}()

	// walk the file tree
	for _, path := range opts.InPaths {
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
				if opts.Verbose {
					fmt.Println(tarHeader.Name)
				}
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
			tarHeader, err := tar.FileInfoHeader(info, info.Name())
			tarHeader.Name = filepath.Join(filepath.Dir(path), tarHeader.Name)
			if opts.Verbose {
				fmt.Println(tarHeader.Name)
			}
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
