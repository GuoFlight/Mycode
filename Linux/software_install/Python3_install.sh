#本脚本用于安装Python3
echo "即将安装EPEL源："
sudo yum -y install epel-release
echo "即将安装Python3.4:"
sudo yum install -y python34
echo "即将安装pip和setuptools:"
curl -O https://bootstrap.pypa.io/get-pip.py
echo "正在安装："
sudo /usr/bin/python3.4 get-pip.py

