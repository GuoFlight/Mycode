#脚本作用：批量重命名
#应用平台：windows
import os
dir="E:\MyFiles\新建文件夹\images\photos\large"
files=os.listdir(dir)
for i in range(1,len(files)+1):
    src="\\".join([dir,files[i-1]])
    dst="\\".join([dir,str(i)+".jpeg"])
    os.rename(src,dst)