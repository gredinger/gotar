package main

import (
	"github.com/caleblloyd/gotar"
	"flag"
	"log"
	"io"
	"os"
)

func main(){
	createPtr := flag.Bool("c", false, "create")
	extractPtr := flag.Bool("x", false, "extract")
	//compressPtr := flag.Bool("z", false, "compress with gzip")
	//verbosePtr := flag.Bool("v", false, "verbose")
	workDirPtr := flag.String("C", "", "change directory")

	flag.Parse()
	create := *createPtr
	extract := *extractPtr
	//compress := *compressPtr
	//verbose := *verbosePtr
	workDir := *workDirPtr

	if (!create && !extract){
		log.Fatal("Either -c (create) or -x (extract) flag must be set")
	}

	if (create){
		if (flag.NArg() < 2){
			log.Fatal("gotar -c [outPath] [inPath]...")
		}
		outPath := flag.Arg(0)
		var w io.Writer;
		if (outPath == "-"){
			w = os.Stdout;
		} else {
			outFp, err := os.Create(outPath)
			if err != nil {
				log.Fatal(err)
			}
			defer outFp.Close()
			w = outFp
		}
		to := &gotar.TarOpts{
			InPaths: flag.Args()[1:],
			WorkDir: workDir,
			Writer: w,
		}
		err := gotar.Tar(to)
		if (err != nil){
			log.Fatal(err)
		}
	}
}