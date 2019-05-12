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

echo; echo "# Test set command2"
echo "## New setting"
rm -r ~/contest/atcoder/abc000 || :
hack s atcoder/abc000
echo; echo "## Test jump2"
pwd
echo "## Jump from the outside quiz directory"
cd `hack j`
pwd
echo "## Jump to not exists directory"
rm -r ~/contest/atcoder/abc000
cd `hack j`
pwd
echo "## Failed initialization"
hack i || :

echo; echo "# Test set command3"
echo "## New setting"
rm -r ~/contest/atcoder/abc100 || :
hack s atcoder/abc100
echo "## Jump to the contest directory"
cd `hack j`
pwd
echo "## AtCoder initialization"
hack i
tree
echo "## Jump to the quiz directory"
cd `hack j`
pwd
echo "## Quiz initialization"
rm * || :
hack i
ls
echo "## Submit"
cat <<EOF > main.go
package main
import "fmt"
func main() {fmt.Println("one\ntwo\nthree")}
EOF
hack sub

echo "## Jump to specified directory"
cd `hack j abc100_d`
pwd
echo "## No Jump from the last directory"
cd `hack j`
pwd

cd $HOME
rm -r ~/contest/atcoder/abc100

echo; echo "# Test codejam"
rm -r ~/contest/codejam/2000round1a || :
hack s codejam/2000round1a
echo "## Jump to the contest directory"
cd `hack j`
pwd
echo "## CodeJam initialization"
hack i || :
ls
echo "## Jump to the quiz directory"
cd `hack j`
pwd
cd $HOME
rm -r ~/contest/codejam/2000round1a

echo; echo "# Test free"
rm -r ~/contest/free/test || :
hack s free/test
echo "## Jump to the contest directory"
cd `hack j`
pwd
echo "## CodeJam initialization"
hack i -n 4
ls
echo "## Jump to the quiz directory"
cd `hack j`
pwd
cd $HOME
rm -r ~/contest/free/test
