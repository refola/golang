#!/bin/bash
# Update the Go installation.

echo "Updating Go!"
cd $GOROOT
echo "Getting release info"
hg pull
echo "Updating to latest release"
hg update release
cd src
echo "Rebuilding source"
time ./all.bash
exit
