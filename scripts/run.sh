#!/usr/bin/env bash

GITROOT=`git rev-parse --show-toplevel`

cd $GITROOT/src
go run *.go
cd -
