# Hack

`hack` is CLI tool to support programming contests. It
can create sample input/output file easily and test if
your code is correct.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Hack](#hack)
    - [Supported contests](#supported-contests)
    - [Installation](#installation)
    - [Quick Start](#quick-start)

<!-- markdown-toc end -->

## Supported contests

- AtCorder
- Google Code Jam

## Installation

```
$ git clone https://github.com/Ladicle/hack.git
$ make install
```

## Quick Start

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

That's all preparation before contest begin. After opened
the quizzes, you probably read a quiz sentence, and create
a sample input/output files for testing. `sample` command
supports such action as follows.

```
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
