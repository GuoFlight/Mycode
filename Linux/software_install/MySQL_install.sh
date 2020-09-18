#!/bin/bash
#作用：从官网下载安装MySQL
#平台：CentOS-7
echo "下载rpm包:"
wget "https://dev.mysql.com/get/mysql80-community-release-el7-1.noarch.rpm"
sudo rpm -Uvh mysql*.rpm
echo "安装MySQL:"
sudo yum install -y mysql-community-server
