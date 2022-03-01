# wordtl

`wordtl` is a `tool` that allows `anyone` to `help them solve a wordle`.

Invoke `wordtl` with possibilites to help you find a list of words that meet your criteria. For example:
- What is a list of 5 letter words that have:
  - "T" as the first letter, and
  - "R" is not in postion 2,
  - "IES" are excluded.

## Prerequisites

Not much here. You can run wordtl on `Windows or Mac`.

If you would like to compile, test, and build the code then you will need `Golang` installed.

### Word List
You must provide a word file for wordtl to read. Please see "Installing wordtl" for details.

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
### Word List
wordtl reads `CSW21.txt` in the same directory as wordtl (for macOS) or wordtl.exe (for Wndows) by default. It can be downloaded from https://ia903406.us.archive.org/31/items/csw21/CSW21.txt. 
- Ensure that a word file (an ASCII text file with one word per line) is downloaded and avaiable for wordctl to read. 
- CSW21.txt is the same as CSW22.txt from https://www.dropbox.com/s/gagbzhzbe2900ua/CSW22.txt and is described by the Collins Coalition here: https://www.cocoscrabble.org/lexicon.

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
  -f string
        Word File: Name/Path of ASCII text file containing one word per line. (default "CSW21.txt")
  -p string
        Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '-' 1) must include all letters specified by the -w flag and 2) can be any other letter that is not excluded by the -x flag. Example value of 't----' would lookup words with a 't' in the beginning of a 5 letter word.
  -s    Print statistics of letter distribution for each letter position.
  -w string
        Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '-'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '-' in the -p flag.
  -x string
        Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.
  -(1-9) string
      	Letters that don't belong in this position (each position, 1 through 9, has their own flag): Letters that appear in the word, but not in postion #(1-9) Example value of 'ies' means that 'i', 'e', or 's' cannot appear in position #(1-9).
```
## Example
### Example Input
What is a list of 5 letter words that have:
  - "T" as the first letter, and
  - "R" is not in postion 2,
  - "IES" are excluded.

Would be specified by `./wordtl -p t---- -w r -2 r -x ies`

### Example Output
The Example input has the following output:
```
Word length: 5
Word pattern: 't----'
Wild Card letters: 'tr'
Excluded letters: 'ies'
Can't use letters in postion #2: 'r'

Possible matching words (37):
tabor talar tardo tardy targa taroc tarok tarot tarry tarty tatar taxor tayra
tharm thorn thoro thorp thraw throb throw thrum thurl tolar torah toran torch
torot torta tuart tubar tugra tumor turbo turfy turnt tutor tyran 

Try these letters (17):
a=27 o=22 h=11 u=10 y=6 n=4 b=4 l=3 m=3 w=2 d=2 c=2 g=2 f=1 k=1 x=1 p=1 
```

#### Interpreting the Output

##### Possible matching words
This is a list of words that match the input criteria. The answer is in here!

##### Try these letters
This is a list of letters in the `Possible matching words` in the order of their occurrances (greatest to least). In this case, One may want to try `./wordtl -p t---- -w roah -2 r` giving `torah` as the next guess to narrow down the possibilities further.

###### I didn't get any results?
You specified to many required items and nothing matched your query. Simply remove some of the constraints to open the query to more results.

## Building/Testing wordtl
`wordtl` is developed in Golang. You will need to download Golang from https://golang.org/doc/install. You can install additional developer tools such as an IDE if you would like, but it is not required.

`TLDR;` you can run `build-all` to build all executables and run unit tests.

### Golang Version
This code was compiled with `go version go1.16.2 darwin/amd64`. Run `go version` to see what you are using.

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