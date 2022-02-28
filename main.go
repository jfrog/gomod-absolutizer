package main

import (
	"flag"
	"os"
)

func main() {
	exitIfError(prepareAndAbsolutize())
}

func prepareAndAbsolutize() error {
	goModPath := flag.String("goModPath", "", "Path to go.mod")
	workingDir := flag.String("wd", "", "Path to working directory")
	flag.Parse()
	if *goModPath == "" || *workingDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	args, err := prepareArgs(*goModPath, *workingDir)
	if err != nil {
		return err
	}
	return Absolutize(args)
}
