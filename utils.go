package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type AbsolutizeArgs struct {
	// Path to go.mod file
	GoModPath string
	// The working directory which will be concatenated to the relative path in the go.mod file
	WorkingDir string
}

// Create AbsolutizeArgs struct from the input AbsolutizeArgs
// GoModPath  - Path to go.mod file
// WorkingDir - The working directory which will be concatenated to the relative path in the go.mod file
func prepareArgs(goModPath, workingDir string) (*AbsolutizeArgs, error) {
	absWorkingDir, err := filepath.Abs(workingDir)
	if err != nil {
		return nil, errors.New("Couldn't absolutize working directory " + workingDir)
	}
	return &AbsolutizeArgs{
		GoModPath:  goModPath,
		WorkingDir: absWorkingDir,
	}, nil
}

// Make relative path in go.mod absolute
// AbsolutizeArgs - The input arguments
func Absolutize(args *AbsolutizeArgs) error {
	file, errs := parseGoMod(args.GoModPath)
	if errs != nil {
		for _, err := range errs {
			fmt.Printf("Error: %v\n", err)
		}
		return errors.New("Couldn't parse go.mod file " + args.GoModPath)
	}

	if err := replaceRelativePaths(file, args.WorkingDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		return errors.New("Couldn't replace path to absolute in " + args.GoModPath)
	}

	if err := saveGoMod(file); err != nil {
		for _, err := range errs {
			fmt.Printf("Error: %v\n", err)
		}
		return errors.New("Couldn't save go.mod in " + args.GoModPath)
	}
	return nil
}

// Parse go.mod file and return structurized object
// GoModPath - Path to go.mod file
func parseGoMod(goModPath string) (result *modfile.File, errs []error) {
	file, err := os.Open(filepath.Clean(goModPath))
	if err != nil {
		return nil, append(errs, err)
	}
	// #nosec G307
	defer func() {
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}()

	stat, err := file.Stat()
	if err != nil {
		return nil, append(errs, err)
	}
	goModContent := make([]byte, stat.Size())
	if _, err = file.Read(goModContent); err != nil && err != io.EOF {
		return nil, append(errs, err)
	}

	result, err = modfile.Parse(goModPath, goModContent, nil)
	if err != nil {
		return nil, append(errs, err)
	}
	return
}

// Replace relative paths in go.mod to absolute paths
// file       - Parsed go.mod file
// WorkingDir - The working directory which will be concatenated to the relative path in the go.mod file
func replaceRelativePaths(file *modfile.File, workingDir string) error {
	for _, replace := range file.Replace {
		currentPath := replace.New.Path
		if replace.New.Version == "" && !filepath.IsAbs(currentPath) {
			absPath := filepath.Clean(filepath.Join(workingDir, currentPath))
			fmt.Println("Replacing " + currentPath + " with " + absPath)
			if err := file.AddReplace(replace.Old.Path, "", absPath, ""); err != nil {
				return err
			}
		}
	}
	return nil
}

// Save changes in go.mod file
// file - Parsed go.mod file
func saveGoMod(file *modfile.File) (errs []error) {
	content, err := file.Format()
	if err != nil {
		return append(errs, err)
	}

	goModfile, err := os.Create(file.Syntax.Name)
	if err != nil {
		return append(errs, err)
	}
	// #nosec G307
	defer func() {
		if err = goModfile.Close(); err != nil {
			errs = append(errs, err)
		}
	}()
	_, err = goModfile.Write(content)
	if err != nil {
		errs = append(errs, err)
	}
	return
}

// Exit with error code 1 if the input error is not nil
// err - The error to check
func exitIfError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err.Error())
	os.Exit(1)
}
