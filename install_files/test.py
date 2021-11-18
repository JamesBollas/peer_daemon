import requests
import os
from dotenv import dotenv_values
import requests_unixsocket
import socket
import asyncio

def convert_unix_socket(socket):
    return localsocket.replace("/","%2F")

requests_unixsocket.monkeypatch()

config = dotenv_values(".env")

localsocket = config['LOCAL_SOCKET']

localsocket = convert_unix_socket(localsocket)

with open(config['HOSTNAME_PATH'],'r') as f:
    hostname = f.read()

hostname = "http://" + hostname[:-1]


headers = {'service':"chat", "id":b"1"}

url = "http+unix://" + localsocket+"/getmessage"

print(url)

r = requests.post(url, headers=headers)
print(r.content)


# headers = {'address':hostname}

# url = "http+unix://" + localsocket+"/sendmessage"

# print(url)

# r = requests.post(url, data=b"chat\nHello World!", headers=headers)

# print(str(r.content))