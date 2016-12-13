package gotar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

type UntarOpts struct {
	Gzip    bool
	Reader  io.Reader
	Verbose bool
	WorkDir string
}

func Untar(opts *UntarOpts) error {
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

	// set the reader
	var r io.Reader
	if opts.Gzip {
		var err error
		r, err = gzip.NewReader(opts.Reader)
		if err != nil {
			return err
		}
	} else {
		r = opts.Reader
	}
	tarReader := tar.NewReader(r)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		info := hdr.FileInfo()
		if opts.Verbose {
			fmt.Println(hdr.Name)
		}
		if info.IsDir() {
			err := os.MkdirAll(hdr.Name, info.Mode())
			if err != nil {
				return err
			}
		} else {
			// tar files don't always store directories, attempt to create first
			err := os.MkdirAll(path.Dir(hdr.Name), 0755)
			if err != nil {
				return err
			}
			outFp, err := os.Create(hdr.Name)
			if err != nil {
				return err
			}
			_, err = io.Copy(outFp, tarReader)
			outFp.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
