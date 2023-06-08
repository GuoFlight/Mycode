#!/bin/bash

#################################################################
# 脚本作用: Linux的回收站
# 作者: 郭少
# 版本: v1.1
# 最近一次更新时间:2023-06-09
#################################################################

garbage_dir="$HOME/garbage"

set -e

#若文件不存在
if [ ! -e $1 ]; then
        echo "删除失败 '$1': No such file or directory"
        exit 1
fi

get_new_filename(){
    #文件名与删除时间的分隔符
    sep_char="_"

    
        #分离出文件名
    old_filename=`echo "$1" | awk -F/ '{print $NF}'`

    len_old_filename=${#old_filename} 
    let temp=len_old_filename-1

    #若参数为目录,则去除最后的/ 
    if [ ${old_filename:temp} = "/" ]; then
        old_filename=${old_filename:0:$temp}
    fi  

    #得到新文件名
    cur_time=`date "+%Y%m%d%H%M%S"`
    new_filename="$old_filename$sep_char$cur_time"
        
    echo $new_filename
    return 0
}

#创建垃圾桶
mkdir -p $garbage_dir

new_filename=`get_new_filename $1`

mv $1 $garbage_dir/$new_filename
