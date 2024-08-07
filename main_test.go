package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/module"
)

func TestAbsolutizeUnix(t *testing.T) {
	if os.PathSeparator != '/' {
		t.Skip("Skipping test on Windows")
	}
	args := prepareGoMod(t, "unix", "/test/dummy/abs")
	defer os.Remove(args.GoModPath)
	assert.NoError(t, Absolutize(args))

	// Check results
	goMod, errs := parseGoMod(args.GoModPath)
	assert.Nil(t, errs)
	assert.Equal(t, module.Version{Path: "/test/dummy/a", Version: ""}, goMod.Replace[0].New)
	assert.Equal(t, module.Version{Path: "/test/dummy/a", Version: ""}, goMod.Replace[0].New)
	assert.Equal(t, module.Version{Path: "c", Version: "v1.0.1"}, goMod.Replace[2].New)
	assert.Equal(t, module.Version{Path: "/some/absolute/path/unix", Version: ""}, goMod.Replace[3].New)
}

func TestAbsolutizeWin(t *testing.T) {
	if os.PathSeparator != '\\' {
		t.Skip("Skipping test on Unix")
	}
	args := prepareGoMod(t, "windows", "C:\\test\\dummy\\abs")
	defer os.Remove(args.GoModPath)
	assert.NoError(t, Absolutize(args))

	// Check results
	goMod, errs := parseGoMod(args.GoModPath)
	assert.Nil(t, errs)
	assert.Equal(t, module.Version{Path: "C:\\test\\dummy\\a", Version: ""}, goMod.Replace[0].New)
	assert.Equal(t, module.Version{Path: "C:\\test\\dummy\\a", Version: ""}, goMod.Replace[0].New)
	assert.Equal(t, module.Version{Path: "c", Version: "v1.0.1"}, goMod.Replace[2].New)
	assert.Equal(t, module.Version{Path: "C:\\some\\absolute\\path\\", Version: ""}, goMod.Replace[3].New)
}

func TestErroneousGoMod(t *testing.T) {
	goMod, errs := parseGoMod(filepath.Join("testdata", "erroneous", "go.mod"))
	assert.Len(t, errs, 1)
	assert.Error(t, errs[0])
	assert.Nil(t, goMod)
}

func TestMinGoVersionGoMod(t *testing.T) {
	_, errs := parseGoMod(filepath.Join("testdata", "version", "minor", "go.mod"))
	assert.Nil(t, errs)
}

func TestMinGoVersionWithPatchGoMod(t *testing.T) {
	_, errs := parseGoMod(filepath.Join("testdata", "version", "patch", "go.mod"))
	assert.Nil(t, errs)
}

func TestBadGoModPath(t *testing.T) {
	goMod, errs := parseGoMod(filepath.Join("testdata", "nonexist", "go.mod"))
	assert.Len(t, errs, 1)
	assert.ErrorIs(t, errs[0], os.ErrNotExist)
	assert.Nil(t, goMod)
}

func prepareGoMod(t *testing.T, goModDir, workingDir string) *AbsolutizeArgs {
	goModPath, err := os.CreateTemp("", "go.mod")
	assert.NoError(t, err)
	bytesRead, err := os.ReadFile(filepath.Join("testdata", goModDir, "go.mod"))
	assert.NoError(t, err)
	err = os.WriteFile(goModPath.Name(), bytesRead, 0644)
	assert.NoError(t, err)
	return &AbsolutizeArgs{GoModPath: goModPath.Name(), WorkingDir: workingDir}
}
