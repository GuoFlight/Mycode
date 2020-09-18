# -*- coding: UTF-8 -*-
#######################################################
# 脚本作用：下载指定Group中的所有项目
# 作者：flightguofei
# 安装requests命令：python3 -m pip install requests
#######################################################

import requests
import json
import os

workspace = '/Users/didi/Desktop/test-dir'      #Projects要下载到这个目录下
token = "xxxxxxxxxxxx"
group_name = "OP"               #想要拉取的Group名，区分大小写

api_groups = 'http://xxx.com/api/v4/groups?private_token={token}'
api_projects = 'http://xxx.com/api/v4/groups/{group_id}/projects?private_token={token}'
api_groups = api_groups.replace('{token}',token)
api_projects= api_projects.replace('{token}',token)

#创建目录
if not os.path.isdir(workspace):		#判断path是否为一个存在的文件夹
    os.makedirs(workspace)
os.chdir(workspace)	#转移工作目录

r_groups = requests.get(api_groups)
groups = json.loads(r_groups.text)     #得到字典类型的列表
for i in groups:
    if i['name'] == group_name:
        group_id = i['id']
        api_projects = api_projects.replace('{group_id}',str(group_id))
        r_projects = requests.get(api_projects)
        projects =  json.loads(r_projects.text)     #得到指定Group中的所有Projects，字典类型的列表

        #拉取所有Projects
        for i in projects:
            print(i['ssh_url_to_repo']," :")
            os.system('git clone '+i['ssh_url_to_repo'])  # 运行命令,返回0则执行成功
            print()
        break
