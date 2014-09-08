#!/bin/bash
killall serve

real=`readlink -f $0` # dereference any symbolic links used to run this
path=`dirname $real` # get directory name
refola=$HOME/.config/refola


start(){
	$path/bin/serve -datapath=$refola/server/data 
	# TODO: rewrite program to get path automatically
}
start # &
exit
