#�ű����ã�����������
#Ӧ��ƽ̨��windows
import os
dir="E:\MyFiles\�½��ļ���\images\photos\large"
files=os.listdir(dir)
for i in range(1,len(files)+1):
    src="\\".join([dir,files[i-1]])
    dst="\\".join([dir,str(i)+".jpeg"])
    os.rename(src,dst)