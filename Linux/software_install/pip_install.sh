#!/bin/bash
echo "下载.py文件:"
wget https://bootstrap.pypa.io/get-pip.py
echo "开始安装:"
python get-pip.py
echo "pip版本:"
pip -V
