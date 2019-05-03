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
  set         Switch contest current contest
  get         get contests
  init        Initialize workspace for the contest
  jump        Get current quiz directory
  test        Test main program
  copy        Copy main program to clipboard
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

## Quick Started

Set the next contest to work on.

```
$ hack set atcoder/abc100
🤖 < OK! I set "atcoder/abc100" for the next contest
```

Jump to the contest root directory.

```
$ cd (hack jump)
```

Initialize the current contest.

```
$ hack init
🤖 < Sure! I'll setup environment for "abc100" contest.

 ✓ Scraping abc100 quizzes 🔎
 ✓ Creating 4 quiz directories 📦
 ✓ Scraping abc100_a quizzes 📥
 ✓ Scraping sample #1 📝
 ✓ Scraping sample #2 📝
 ✓ Scraping sample #3 📝
 ✓ Scraping sample #3 📝
 ✓ Scraping abc100_b quizzes 📥
 ✓ Scraping sample #1 📝
 ✓ Scraping sample #2 📝
 ✓ Scraping sample #3 📝
 ✓ Scraping sample #3 📝
 ✓ Scraping abc100_c quizzes 📥
 ✓ Scraping sample #1 📝
 ✓ Scraping sample #2 📝
 ✓ Scraping sample #3 📝
 ✓ Scraping sample #4 📝
 ✓ Scraping sample #4 📝
 ✓ Scraping abc100_d quizzes 📥
 ✓ Scraping sample #1 📝
 ✓ Scraping sample #2 📝
 ✓ Scraping sample #3 📝
 ✓ Scraping sample #4 📝
```

Go to the first quiz directory.

```
$ cd (hack j)
```

Let's programming!

```
$ emacs main.cpp
```

After writing the code, test it.

```
$ hack test
[AC] Sample #1
[AC] Sample #2
[AC] Sample #3
```

There is no problem, copy and submit it :)

```
$ hack copy
```
