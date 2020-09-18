#!/bin/bash
#脚本作用:安装docker-compose
#最新版安装网址：https://docs.docker.com/compose/install/#install-compose
echo "下载Docker-Compose："
sudo curl -L "https://github.com/docker/compose/releases/download/1.23.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
echo "安装版本："
docker-compose --version
