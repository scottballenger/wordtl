# wordtl

`wordtl` is a `tool` that allows `anyone` to `help them solve a wordle`.

Invoke `wordtl` with possibilites to help you find a list of words that meet you criteria. For example:
- What is a list of 5 letter words that have "T" as the first letter and have "R" that is not in the second letter?

## Prerequisites

Not much here. You can run wordtl on `Windows or Mac`.

If you would like to compile, test, and build the code then you will need `Golang` installed.

## Installing wordtl

To install wordtl, follow these steps:

macOS:
```
copy "wordtl" to your machine
chmod 775 wordtl # to make it executable
```

Windows:
```
copy "wordtl.exe" to your machine
```
## Using wordtl

To use wordtl, follow these steps:

macOS:
```
cd <dir that contains wordtl>
./wordtl
```

Windows:
```
cd <dir that contains wordtl.exe>
wordtl.exe
```

### Usage:
```
Usage of ./wordtl:
  -a	Auto Shoot Mode (default - Manual Shot)
  -d float
    	Detonation Radius (meters) (default 20)
  -e	English Units (default - Metric)
  -m	Real-time Target Movement (default - Pause Target During Shot Decision)
  -p	Print Shot Profile
```
## Building/Testing wordtl
`wordtl` is developed in Golang. You will need to download Golang from https://golang.org/doc/install. You can install additional developer tools such as an IDE if you would like, but it is not required.

### Golang Version
This code was compiled with `go version go1.15.7 darwin/amd64`. Run `go version` to see what you are using.

### Compile the Code and Build Executables

To build the code and create the stand-alone executable for your platform, just run the following command:

```
cd wordtl
go build
```

macOS:
This will create the executable `wordtl` that you can run.

Windows:
This will create the executable `wordtl.exe` that you can run.

#### Compiling the Code for other Platforms

For the complete list of operating systems and architectures that can be cross compiled, see https://golang.org/doc/install/source#environment

##### Compiling for Windows from macOS

If you are on a macOS platform and want to create an executable for Windows, then you would run the following:

```
cd wordtl
GOOS=windows go build
```

This will create the executable `wordtl.exe` that you can run on Windows.

##### Compiling for macOS from Windows

If you are on a Windows platform and want to create an executable for macOS, then you would run the following:

```
cd wordtl
GOOS=darwin go build
```

This will create the executable `wordtl` that you can run on macOS.

### Run Unit Tests

To run the unit tests for your platform, just run the following command:

```
cd wordtl
go test
```

Upon execution, you should see something that ends with:
```
PASS
ok      wordtl    0.319s
```

## Contributing to wordtl
To contribute to wordtl, follow these steps:

1. Fork this repository.
2. Create a branch: `git checkout -b <branch_name>`.
3. Make your changes and commit them: `git commit -m '<commit_message>'`
4. Push to the original branch: `git push origin wordtl/<location>`
5. Create the pull request.

Alternatively see the GitHub documentation on [creating a pull request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).


## License

This project uses the following license: [MIT License](https://github.com/scottballenger/wordtl/blob/main/LICENSE).