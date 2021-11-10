import requests
import os
from dotenv import dotenv_values
import requests_unixsocket

def convert_unix_socket(socket):
    return localsocket.replace("/","%2F")

requests_unixsocket.monkeypatch()

config = dotenv_values(".env")

localsocket = config['LOCAL_SOCKET']

localsocket = convert_unix_socket(localsocket)

with open(config['HOSTNAME_PATH'],'r') as f:
    hostname = f.read()

hostname = "http://" + hostname[:-1]


headers = {'service':"chat"}

url = "http+unix://" + localsocket+"/connect"

print(url)

r = requests.post(url, headers=headers)
print(r.json())


headers = {'address':hostname}

url = "http+unix://" + localsocket+"/sendmessage"

print(url)

r = requests.post(url, data=b"chat\nHello World!", headers=headers)

print(r)