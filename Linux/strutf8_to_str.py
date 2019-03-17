#脚本作用:将str类型的utf-8码转换成字符串
import re
s="%E4%BA%A7%E5%93%814"		#测试用


def strutf8_to_str(s):
    ret=re.findall("%..%..%..",s)
    for i in ret:
       j=i.split("%")
       tmp1=int(j[1],16)
       tmp2=int(j[2],16)
       tmp3=int(j[3],16)
       tmp = bytes((tmp1, tmp2, tmp3))
       tmp=str(tmp, "utf8")
       s=s.replace(i,tmp)
    return s
print(strutf8_to_str(s))
