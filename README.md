# Hack

`hack` is a commandline tool to assist your programming content.

## Installation

```bash
go install github.com/Ladicle/hack@latest
```

## Usage

```
Hack assists your programming contest.

Usage:
  hack [command]

Available Commands:
  add         Add resource to the contest directory
  go          Print path to the directory
  help        Help about any command
  init        Create directories and download samples
  next        Print next quiz directory
  open        Open current task page
  test        Test your program

Flags:
      --config string   path to the configuration file (default "~/.config/hack")
  -h, --help            help for hack

Use "hack [command] --help" for more information about a command.
```

## Quick Started

Write configuration and save it as a `~/.config/hack` file.

```
atcoder:
  user: <username>
  pass: <password>
basedir: <path/to/directory>
```

Initialize contest directory and download samples.

```
$ hack init abc100
Initialize directory for abc100:
 ✓ Scraping task abc100_a
   ✓ Scraping sample #1
   ✓ Scraping sample #2
   ✓ Scraping sample #3
 ✓ Scraping task abc100_b
   ✓ Scraping sample #1
   ✓ Scraping sample #2
   ✓ Scraping sample #3
 ✓ Scraping task abc100_c
   ✓ Scraping sample #1
   ✓ Scraping sample #2
   ✓ Scraping sample #3
```

After writing the code, test and submit it if the program pass all test cases.

```
$ hack test
[AC] Sample #1
[AC] Sample #2
[AC] Sample #3
Copy main.py!
```
