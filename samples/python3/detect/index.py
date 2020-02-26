import os

def handler(event):
    funcName = os.getenv("funcName")
    if funcName == None:
        return "hello world from aliyun"
    else:
        return "hello world from HCloud" + funcName