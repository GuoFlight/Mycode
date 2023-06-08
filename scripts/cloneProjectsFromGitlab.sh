#!/bin/bash

#####################
# 脚本作用：拉取gitlab某目录中的所有项目
# 使用方法：将想要拉取的项目名写在文件中，并配置好文件路径
# 作者：郭飞
#####################

file="/Users/didi/Desktop/router_inrouter.txt"		#需要拉取的git地址(需要绝对路径)
dir="/Users/didi/Desktop/router_inrouter"	#项目保存的目录

if [ ! -e $file ]; then
	echo "文件不存在"
	exit 1
fi
mkdir -p  $dir
cd $dir
for i in `cat $file`; do
	git clone $i
done




