# Jack Compiler

Compiler from The Jack language to The Hack platform (Nand2Tetris Part 2)

## Requirements

You need Go to develop, build, and run this project.

## Installation

1. Clone this repo:

```
git clone https://github.com/hazemKrimi/jack-compiler
```

2. To run this against Jack files, run the following command with the path of the `.jack` file (or directory):

```
./out/jack-compiler <path_to_jack_file_or_directory>
```

The result `.vm` file will be written in the same location as the source file.

## Build

To build this project run the following command:

```
go build -o out
```

You will find the executable in the `out` directory.
