#coding:utf-8
#作用：更换软件源
import os
address=['阿里源','英格源','中科大源','163网易源']
for i in range(1,len(address)+1):
	print(str(i)+'、'+address[i-1])
address_input=input("请输入软件源序号:")
while (not (address_input.isdigit() and int(address_input) >= 1 and int(address_input) <=len(address))):
	address_input=input("请输入正确的软件源序号:")
address_input=int(address_input)
os.system("rm -f /etc/yum.repos.d/CentOS-Base.repo")
if address_input == 1:
	os.system("curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo")
elif address_input == 2:
	os.system("curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.eagleslab.com:8889/base.repo")
elif address_input == 3:
	os.system("curl -o /etc/yum.repos.d/CentOS-Base.repo https://lug.ustc.edu.cn/wiki/_export/code/mirrors/help/centos?codeblock=3")
elif address_input == 4:
	os.system("curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.163.com/.help/CentOS6-Base-163.repo")
os.system("yum makecache")
