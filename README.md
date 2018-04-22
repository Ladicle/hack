# Hack

`hack` is CLI tool to support programming contests. It
can create sample input/output file easily and test if
your code is correct.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Hack](#hack)
    - [Supported contests](#supported-contests)
    - [Installation](#installation)
    - [Usage](#usage)
    - [Quick Start](#quick-start)
        - [Preparation](#preparation)
        - [Create sample input/output files](#create-sample-inputoutput-files)
        - [Check your answer using samples](#check-your-answer-using-samples)

<!-- markdown-toc end -->

## Supported contests

- AtCorder
- Google Code Jam

## Installation

```
$ git clone https://github.com/Ladicle/hack.git
$ make install
```

## Usage

```
❯❯❯ hack --help
Usage: hack [OPTIONS] COMMAND

Options:
  -c --config          Configuration path (default: ~/.hack)
  -o --output          Output directory (default: ~/contest)
  -h --help            Show this help message

Commands:
  version              Show this command version
  set <PATH>           Set contest information
  info                 Info shows information
  open                 Open shows contest page in browser
  jump [QUIZ]          Jump move to quiz directory. if not specified quiz option, move to next quiz directory.
  sample [-i] [NUMBER] Sample creates input/output sample
  test [NUMBER]        Test tests your code with all samples if you don't specified the number
  copy                 Corpy copies your code to clipboard
```

## Quick Start

### Preparation

First of all, you need to set current contest information
as follows. If you do not already have hack command, please 
refer the [Installation](#installation) section.

*atcoder*
```
$ hack set atcoder/abc/93
```

*codejam*
```
$ hack set codejam/2018/practice
```

Next, you want to see contest information, you can print it
in terminal, or open website in browser.

```
$ hack info
Contest:
  Name: atcoder
  Path: /Users/aigarash/contest/atcoder/abc/093
  Quizzes:
  - a
  - b
  - c
  - d
  URL: https://abc093.contest.atcoder.jp/
CurrentQuizz: a

$ hack open
# open 'https://abc093.contest.atcoder.jp/' in browser.
```

### Create sample input/output files

That's all preparation before contest begin. After opened
the quizzes, you probably read a quiz sentence, and create
a sample input/output files for testing. `sample` command
supports such action as follows.

```
$ cd (hack jump)
$ hack sample
"1.in" is already exists. Skip it? (y/n): n

1.in:
2
1.000000
1.414213

1.out:
Case #1:
0.5 0 0
0 0.5 0
0 0 0.5
Case #2:
0.3535533905932738 0.3535533905932738 0
-0.3535533905932738 0.3535533905932738 0
0 0 0.5

Continue to create sample? (y/n): n

$ ls
1.in  1.out
```

### Check your answer using samples

After you finished to write a code, you can check it using
`test` command. The results are displayed for  each sample
inputs, and if your code is correct, test command prints
`AC`. Otherwise, it prints debugging log.

```
$ hack test
[AC] input #1
[AC] input #2
[AC] input #3
```

*Status Codes*

* AC: Answer is correct
* WA: Wrong answer
* LT: Long time
