#!/bin/bash

#################################################################
# 脚本作用: Linux的回收站
# 作者: 京城郭少
# 版本: v3.0
# 最近一次更新时间: 2025-04-24
#################################################################
set -e
garbage_dir="$HOME/.garbage"

function init(){
        # 创建回收站
        mkdir -p $garbage_dir
}

function get_new_filename(){
    # 文件名与删除时间的分隔符
    sep_char="_"

    # 分离出文件名
    old_filename=`basename $target`

    # 得到新文件名
    cur_time=`date "+%Y%m%d%H%M%S${sep_char}${RANDOM}"`
    new_filename="$old_filename$sep_char$cur_time"

    # 返回结果       
    echo $new_filename
    return 0
}
function del_target(){
    target=$1
    [ "$target" == "" ] && echo "程序异常" && exit 1
    new_filename=`get_new_filename`
    mv $target $garbage_dir/$new_filename
}

################################## mian ##################################

init
for target in "$@"; do 
    del_target $target
done