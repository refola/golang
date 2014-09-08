#!/bin/bash

# Convenience script for committing all changes to this repository in a single command

# Get script location
DELINKED=`readlink -f "$0"`
HERE="`dirname "$DELINKED"`"
OFFSET="../" # Where the repo's root is relative to the script
cd $HERE/$OFFSET

echo "Setting files to add/remove"
hg addremove

echo "Committing"
hg commit

echo "Pushing changes"
hg push

echo "Exiting script"
exit
