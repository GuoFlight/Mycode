########################
# 脚本作用：计算发票总额
# 使用方法：指定发票所在路径，并将发票金额设置在文件名开头。如188.12元-xxx.pdf
# 作者：郭飞
########################
import os
import re
workspace="/Users/didi/Desktop/0814团建发票"      #发票所在路径

sum=0.0
count=0
files=os.listdir(workspace)
for i in range(len(files)):
    money=re.findall("^[\d]+.[\d]{2}元",files[i])
    if len(money)!=0:           #只计算符合规则的文件
        money=float(money[0][0:-1])
        sum+=money
        count+=1
print("总金额：",sum)
print("文件数量：",count)
