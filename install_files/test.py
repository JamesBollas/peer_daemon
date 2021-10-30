import requests
import os
from dotenv import dotenv_values
import requests_unixsocket

requests_unixsocket.monkeypatch()

config = dotenv_values(".env")

localsocket = config['LOCAL_SOCKET']

localsocket = localsocket.replace("/","%2F")

with open(config['HOSTNAME_PATH'],'r') as f:
    hostname = f.read()

hostname = "http://" + hostname[:-1]

headers = {'address':hostname}

url = "http+unix://" + localsocket+"/sendmessage"

print(url)

r = requests.post(url, data={'hi':"hi"}, headers=headers)