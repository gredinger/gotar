#!/bin/bash
cd $(dirname $0)

function compare_dirs {
    diff -r -q $1 $2
    local rc=$?
    if [ $rc -ne 0 ]; then
        >&2 echo "Failed: compare_dirs $1 $2"
        exit $rc
    fi
    return $rc
}

rm -rf tests && mkdir -p tests

./cli -c tests/test1.tar testdir
mkdir -p tests/test1 && ./cli -x -C tests/test1 tests/test1.tar
compare_dirs testdir tests/test1/testdir

./cli -c tests/test2.tar testdir/* testdir/.a.txt
mkdir -p tests/test2 && ./cli -x -C tests/test2 tests/test2.tar
compare_dirs testdir tests/test2/testdir

cd ../
./cli/cli -C cli/ -c cli/tests/test3.tar testdir
cd cli
mkdir -p tests/test3 && ./cli -x -C tests/test3 tests/test3.tar
compare_dirs testdir tests/test3/testdir

./cli -c -v -z tests/test4.tar.gz testdir
mkdir -p tests/test4 && ./cli -x -v -z -C tests/test4 tests/test4.tar.gz
compare_dirs testdir tests/test4/testdir

# stdout
mkdir -p tests/test5 && ./cli -c - testdir | ./cli -x -C tests/test5 -
compare_dirs testdir tests/test5/testdir

rm -rf tests
