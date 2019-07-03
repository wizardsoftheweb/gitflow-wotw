#!/bin/bash
#
# Heavily inspired by this script:
# https://gist.github.com/peterjmit/3864743

REPEATS=100
COMMAND='echo rad'

run_once() {
    /usr/bin/time -f "%E,%U,%S" ${COMMAND}
}

run_fully() {
    for (( index=0; index < $REPEATS; index++ ))
    do
        run_once
    done
}

run_fully
