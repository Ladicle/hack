# Hack

`hack` is CLI tool to support programming contests.

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
