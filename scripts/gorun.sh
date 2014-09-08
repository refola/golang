#!/bin/bash
# Run a single-file Go program and remove generated files.

name=${1%\.go}
echo "Running go program: $name"

echo "Formatting $name.go"
gofmt -w $name.go

echo "Compiling $name.go"
8g $name.go

echo "Linking $name.8"
8l -o $name.bin $name.8

echo "Running $name.bin"
./$name.bin

echo "Deleting generated files"
rm $name.8
rm $name.bin

exit
