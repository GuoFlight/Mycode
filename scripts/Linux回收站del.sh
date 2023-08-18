#!/bin/bash

#################################################################
# 脚本作用: Linux的回收站
# 作者: 郭少
# 版本: v2.0
# 最近一次更新时间:2023-08-19
#################################################################

garbage_dir="$HOME/garbage"
target="$1"
set -e

function init(){
        #若文件不存在
        if [ ! -e "$target" ]; then
                echo "删除失败 '$target': No such file or directory"
                exit 1
        fi
        #创建垃圾桶
        mkdir -p $garbage_dir
}

get_new_filename(){
    #文件名与删除时间的分隔符
    sep_char="_"
    
    #分离出文件名
    old_filename=`basename $target`

    #得到新文件名
    cur_time=`date "+%Y%m%d%H%M%S"`
    new_filename="$old_filename$sep_char$cur_time"
        
    echo $new_filename  # 返回值
    return 0
}

################################## mian ##################################

init

new_filename=`get_new_filename`

mv $target $garbage_dir/$new_filename
