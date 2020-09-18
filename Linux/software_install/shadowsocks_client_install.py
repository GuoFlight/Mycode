#   适用环境:CentOS7
#   解释器:Python3
#   作用:安装shadowsocks
#   作者:郭老师

import os

print("安装pip：")
os.system("sudo yum -y install epel-release")
os.system("sudo yum -y install python-pip")

print("安装shadowsocks客户端:")
os.system("sudo pip install shadowsocks")

print("配置shadowsocks客户端:")
os.system("sudo mkdir -p /etc/shadowsocks")
shadowsocks="""{
    "server":"SERVER_IP",
    "server_port":SERVER_PORT,
    "local_address": "127.0.0.1",
    "local_port":LOCAL_PORT,
    "password":"PASSWD",
    "timeout":300,
    "method":"METHOD",
    "fast_open": false,
    "workers": 1
}"""
SERVER_IP=input("输入代理服务器ip:")
SERVER_PORT=input("输入代理服务器端口:")
LOCAL_PORT=input("输入本地代理端口[默认1080]:")
if LOCAL_PORT=="":
    LOCAL_PORT="1080"
PASSWD=input("输入代理服务器代理密码:")
method="""加密方式:
    1.aes-256-cfb
"""
print(method)
temp=input("请选择：")
if temp == '1':
    METHOD="aes-256-cfb"
shadowsocks=shadowsocks.replace("SERVER_IP",SERVER_IP)
shadowsocks=shadowsocks.replace("SERVER_PORT",SERVER_PORT)
shadowsocks=shadowsocks.replace("LOCAL_PORT",LOCAL_PORT)
shadowsocks=shadowsocks.replace("PASSWD",PASSWD)
shadowsocks=shadowsocks.replace("METHOD",METHOD)

with open("/etc/shadowsocks/shadowsocks.json","w") as f:
    f.write(shadowsocks)
    f.flush()


shadowsocks_systemctl="""[Unit]
Description=Shadowsocks
[Service]
TimeoutStartSec=0
ExecStart=/usr/bin/sslocal -c /etc/shadowsocks/shadowsocks.json
[Install]
WantedBy=multi-user.target
"""
with open("/etc/systemd/system/shadowsocks.service","w") as f:
    f.write(shadowsocks_systemctl)

print("设置shadowsocks自启:")
os.system("systemctl enable shadowsocks.service")
print("开启shadowsocks:")
os.system("systemctl start shadowsocks.service")

print("安装privoxy:")
os.system("sudo yum -y install privoxy")
os.system("systemctl enable privoxy")
os.system("systemctl start privoxy")

print("配置privoxy:")
privoxy_config=""
with open("/etc/privoxy/config","r") as f:
    privoxy_config=f.read()
privoxy_config=privoxy_config.replace("#        forward-socks5t","forward-socks5t")
privoxy_config=privoxy_config.replace("127.0.0.1:9050","127.0.0.1:"+LOCAL_PORT)
with open("/etc/privoxy/config","w") as f:
    f.write(privoxy_config)
    f.flush()

print("配置profile:")
profile_text1="export http_proxy=http://127.0.0.1:8118\n"
profile_text2="export https_proxy=http://127.0.0.1:8118\n"
is_enable_profile=input("是否开机其启动[y/n]:")
if is_enable_profile != 'y':
    profile_text1="".join(["#",profile_text1])
    profile_text2 = "".join(["#", profile_text2])
with open("/etc/profile","a") as f:
    f.write(profile_text1)
    f.write(profile_text2)
    f.flush()

print("若要启动客户端,运行以下命令并测试:")
print("export http_proxy=http://127.0.0.1:8118")
print("export https_proxy=http://127.0.0.1:8118")
