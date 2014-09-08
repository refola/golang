#!/bin/bash
# Compile all Go files in the list of directories into one tidy binary, treating each directory as a separate package.
# Also formats all .go files in the given directories.

# Make sure the package name is passed.
if [ -z "$1" ]
then
	echo "Usage: `basename $0` relative/path/to/buildOrderFile"
	exit 1
fi

echo "Compiling packages listed in $1"

fname=`echo $1 | sed "s/.*\///"` # Get the filename without the path.
path=`echo $1 | sed "s/\/[^/]*$/\//"` # Get the path without the filename.
cd $path

# Read in package list
pkgs=`cat $fname | grep -v ^#` # Skip comment lines starting with #.
echo "Packages to compile: `echo $pkgs | sed "s/^//g"`"

package(){
	pkgname=`echo $1 | sed "s/.*\///"`
	echo "Formatting source in package $1"
	echo "gofmt -w ./$1/*.go"
	gofmt -w ./$1/*.go
	echo "Compiling package $1"
	# Get only non-test .go files, get rid of newlines in the list, and add the relative path to each file.
	files=`ls -A -1 $1 | grep \.go$ | grep -v _test\.go$ | sed "s/^/ /g" | sed "s#\s# ./$1/#g"` 
	echo "8g -o $pkgname.8 $files"
	8g -o $pkgname.8 $files # Compile the whole package into one file.
}

# Process all packages in the list.
for pkg in $pkgs
do
	package $pkg
done

# pkg is the name of the last (main) package in the dependency chain
echo "<8l -o result.bin $pkg.8>"
8l -o result.bin $pkg.8
echo "rm ./*.8"
rm ./*.8 # Clean up intermediate results

exit
