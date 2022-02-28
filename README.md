# gomod-absolutizer

**gomod-absolutizer** is a Go library, used to search and replace relative paths in go.mod files to absolute paths.

The library is used by both the [JFrog Idea Plugin](https://github.com/jfrog/jfrog-idea-plugin) and the [JFrog VS Code Extension](https://github.com/jfrog/jfrog-vscode-extension).

## Table of Contents

- [Usage](#usage)
    - [As script](#as-script)
    - [As executable](#as-executable)
    - [As library](#as-library)
- [Example](#example)
- [Tests](#tests)
- [Release Notes](#release-notes)
- [Code contributions](#code-contributions)

## Usage
The program expects two flags:

| Flag        | Description                                                                                        |
|-------------|----------------------------------------------------------------------------------------------------|
| `goModPath` | Path to a go.mod.                                                                                  |
| `wd`        | Path to the working directory, which will be concatenated to the relative path in the go.mod file. |


You may use it in multiple ways:

### As script
```sh
go run . -goModPath=/path/to/go.mod -wd=/path/to/wd
```

### As executable
```sh
go build
./gomod-absolutizer -goModPath=/path/to/go.mod -wd=/path/to/wd
```

### As library
```go
import (
    absolutizer "github.com/jfrog/gomod-absolutizer"
)

func main() {
    args := &absolutizer.AbsolutizeArgs{
        GoModPath:  "/path/to/go.mod",
        WorkingDir: "/path/to/wd",
    }
    err := absolutizer.Absolutize(args)
    // Handle error
}

```

## Example
Given the following go.mod before running this program:
```
replace github.com/jfrog/jfrog-client-go v1.2.3 => github.com/jfrog/jfrog-client-go v1.2.4
replace github.com/jfrog/jfrog-cli-core => ../jfrog-cli-core
```

Running the following command:

```
go run . -goModPath=/Users/frogger/code/jfrog-cli/go.mod -wd=/Users/frogger/code/jfrog-cli
```

Will modify the original go.mod to:
```
replace github.com/jfrog/jfrog-client-go v1.2.3 => github.com/jfrog/jfrog-client-go v1.2.4
replace github.com/jfrog/jfrog-cli-core => /Users/frogger/code/jfrog-cli-core
```

## Tests
To run the tests, execute the following command from within the root directory of the project:

```sh
go test -v ./...
```

# Release Notes
The release notes are available [here](RELEASE.md#release-notes).

# Code Contributions
We welcome community contribution through pull requests.