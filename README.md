# Hack

`hack` is a commandline tool to assist your programming content.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Hack](#hack)
    - [Installation](#installation)
    - [Usage](#usage)

<!-- markdown-toc end -->

## Installation

```bash
go get -u github.com/Ladicle/hack 
```

## Usage

```
Usage:
  hack [command]

Available Commands:
  copy        Copy main program to clipboard
  help        Help about any command
  init        Initialize workspace for the contest
  jump        Get current quiz directory
  list        list contests
  sample      Create sample files
  switch      Switch contest current contest
  test        Test main program
  version     Show this command version

Flags:
  -c, --config string   config file (default /home/ladicle/.hack.yaml)
  -h, --help            help for hack

Use "hack [command] --help" for more information about a command.
```
