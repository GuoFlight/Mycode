#!/usr/bin/env python
########################
# 脚本作用：批量更新所有git项目
# 作者：郭飞
########################

import argparse
import os
#解析参数
parser = argparse.ArgumentParser(description='此命令可以批量更新所有git项目')
parser.add_argument('workDir', type=str, help='存放git Projects的目录')
args = parser.parse_args()
workDir = os.path.abspath(args.workDir)     #去除最后的"/"，并得到绝对路径


projects = os.listdir(workDir)      #得到所有项目
for project in projects:
    if os.path.isdir(workDir + "/" + project):
        os.chdir(workDir + "/" + project)	#转移工作目录
        ret = os.popen("git pull origin master")  # 返回一个管道对象
        print(ret.read())  # 得到命令的输出结果
        ret.close()  # 使用完关闭

