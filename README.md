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
  get         get contests
  help        Help about any command
  init        Initialize workspace for the contest
  jump        Get current quiz directory
  set         Switch contest current contest
  test        Test main program
  version     Show this command version

Flags:
      --alsologtostderr                  log to standard error as well as files
  -c, --config string                    config file (default ~/.hack.yaml)
  -h, --help                             help for hack
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

Use "hack [command] --help" for more information about a command.
```
