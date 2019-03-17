#脚本作用:更换ssh端口号
import os
def add_port(port):
    os.system("".join(["sed -i 's/PORTNUMBER/PORTNUMBER\\nPort ",str(port),"/g' /etc/ssh/sshd_config"]))
    os.system("".join(["firewall-cmd --add-port=", port, "/tcp --permanent"]))
    os.system("firewall-cmd --reload")
    os.system("".join(["semanage port -a -t ssh_port_t -p tcp ",port]))
    os.system("systemctl restart sshd")
def modify_port(old_port,new_port):
    if old_port == "22":
        os.system("".join(["sed -i 's/#Port ",str(old_port),"/Port ", new_port,"/g' /etc/ssh/sshd_config"]))
        os.system("".join(["sed -i 's/Port ",str(old_port),"/Port ", new_port,"/g' /etc/ssh/sshd_config"]))
    else:
        os.system("".join(["sed -i 's/Port ",str(old_port),"/Port ", new_port,"/g' /etc/ssh/sshd_config"]))
    os.system("".join(["firewall-cmd --add-port=", new_port, "/tcp --permanent"]))
    os.system("firewall-cmd --reload")
    os.system("".join(["semanage port -a -t ssh_port_t -p tcp ", new_port]))
    os.system("systemctl restart sshd")
num=0
while num!=1 and num!=2:
    print("1.增加ssh端口")
    print("2.修改ssh端口")
    num=input("请输入序号:")
    if num.isdigit():
        num=int(num)
if num == 1:
    port = ""
    while not port.isdigit():
        port=input("请输入要增加的端口号:")
    add_port(port)
elif num == 2:
    old_port = ""
    new_port = ""
    while not old_port.isdigit():
        old_port = input("请输入旧端口号:")
    while not new_port.isdigit():
        new_port = input("请输入新端口号:")
    modify_port(old_port,new_port)
print("注意:请检查是否成功.")

