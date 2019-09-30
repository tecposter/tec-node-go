#!/bin/bash

set -ex

OUTFILE="fdsafdaterfdsa.out"

if [ "$1" != "" ]; then
	go test $1 -v -coverprofile=$OUTFILE fmt && \
		go tool cover -func=$OUTFILE && \
		cat $OUTFILE && \
		rm $OUTFILE
else
	echo "hello"
fi
