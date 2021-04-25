#!/bin/bash
# @Desc 此脚本用于获取一个指定区间且未被占用的随机端口号
# @Author 郭飞 <guoofei@outlook.com>

PORT=0
#判断当前端口是否被占用，没被占用返回0，反之1
function isListening {
   TCPListeningnum=`netstat -an | grep ":$1 " | awk '$1 == "tcp" && $NF == "LISTEN" {print $0}' | wc -l`
   UDPListeningnum=`netstat -an | grep ":$1 " | awk '$1 == "udp" && $NF == "0.0.0.0:*" {print $0}' | wc -l`
   (( Listeningnum = TCPListeningnum + UDPListeningnum ))
   if [ $Listeningnum == 0 ]; then
       echo "0"
   else
       echo "1"
   fi
}

#得到区间随机数
function getRandomNum {
   shuf -i $1-$2 -n1
}

#得到随机端口
function getRandomPort {
   templ=0
   while [ $PORT == 0 ]; do
       temp1=`getRandomNum $1 $2`
       if [ `isListening $temp1` == 0 ] ; then
              PORT=$temp1
       fi
   done
   echo "$PORT"         #输出此port
}
getRandomPort 1024 10000; #这里指定了1024~10000区间，从中任取一个未占用端口号