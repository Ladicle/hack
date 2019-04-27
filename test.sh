#!/bin/bash -e

echo "# check and install hack command"
make check
make install

echo; echo "# Test help"
hack -h

echo; echo "# Test version"
hack v

echo; echo "# Test get command"
hack g

echo; echo "# Test set command"
echo "## no argument"
hack s || :
echo "## Invalid argument"
hack s atcoder || :
echo "## Normal setting"
hack s atcoder/abc123

echo; echo "# Test jump"
pwd
echo "## Jump from the outside quiz directory"
cd `hack j`
pwd
echo "## Jump to next quiz directory"
cd `hack j`
pwd
echo "## Jump to specified directory"
cd `hack j abc123_d`
pwd
echo "## No Jump from the last directory"
cd `hack j`
pwd

echo; echo "# Test set command2"
echo "## New setting"
hack s atcoder/abc0

echo; echo "# Test jump2"
pwd
echo "## Jump from the outside quiz directory"
cd `hack j`
pwd
rm -r ~/contest/atcoder/abc0
echo "## Jump to not exists directory"
cd `hack j`
pwd

echo; echo "# Test initialization"
echo "TBD"
