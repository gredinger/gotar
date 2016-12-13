package main

import (
	"flag"
	"github.com/caleblloyd/gotar"
	"io"
	"log"
	"os"
)

func main() {
	createPtr := flag.Bool("c", false, "create")
	extractPtr := flag.Bool("x", false, "extract")
	gzipPtr := flag.Bool("z", false, "compress with gzip")
	verbosePtr := flag.Bool("v", false, "verbose")
	workDirPtr := flag.String("C", "", "change directory")

	flag.Parse()
	create := *createPtr
	extract := *extractPtr
	gzip := *gzipPtr
	verbose := *verbosePtr
	workDir := *workDirPtr

	if !create && !extract {
		log.Fatal("Either -c (create) or -x (extract) flag must be set")
	}

	if create {
		if flag.NArg() < 2 {
			log.Fatal("gotar -c [tarPath] [inPath]...")
		}
		tarPath := flag.Arg(0)
		var w io.Writer
		var v bool
		if tarPath == "-" {
			w = os.Stdout
			v = false
		} else {
			outFp, err := os.Create(tarPath)
			if err != nil {
				log.Fatal(err)
			}
			defer outFp.Close()
			w = outFp
			v = verbose
		}
		opts := &gotar.TarOpts{
			Gzip:    gzip,
			InPaths: flag.Args()[1:],
			Verbose: v,
			WorkDir: workDir,
			Writer:  w,
		}
		err := gotar.Tar(opts)
		if err != nil {
			log.Fatal(err)
		}
	} else if extract {
		if flag.NArg() < 1 {
			log.Fatal("gotar -x [tarPath]")
		}
		tarPath := flag.Arg(0)
		var r io.Reader
		var v bool
		if tarPath == "-" {
			r = os.Stdin
			v = false
		} else {
			inFp, err := os.Open(tarPath)
			if err != nil {
				log.Fatal(err)
			}
			defer inFp.Close()
			r = inFp
			v = verbose
		}
		opts := &gotar.UntarOpts{
			Gzip:    gzip,
			Reader:  r,
			Verbose: v,
			WorkDir: workDir,
		}
		err := gotar.Untar(opts)
		if err != nil {
			log.Fatal(err)
		}
	}
}
