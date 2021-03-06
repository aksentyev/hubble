#!/bin/bash

function die() {
  echo $*
  exit 1
}

# Initialize profile.cov
echo "mode: count" > coverage.txt

# Initialize error tracking
ERROR=""

# Test each package and append coverage profile info to profile.cov
for pkg in `cat .testpackages`
do
    #$HOME/gopath/bin/
    go test -v -covermode=count -coverprofile=profile_tmp.cov $pkg || ERROR="Error testing $pkg"
    tail -n +2 profile_tmp.cov >> coverage.txt || die "Unable to append coverage for $pkg"
done

if [ ! -z "$ERROR" ]
then
    die "Encountered error, last error was: $ERROR"
fi
