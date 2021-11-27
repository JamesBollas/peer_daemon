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

def my_hostname():
    return hostname

def getmessageids():
    headers = {'service':"chat"}

    url = "http+unix://" + localsocket+"/getmessageids"

    #print(url)

    r = requests.post(url, headers=headers)
    return r.content.split(b"\n")[:-1]

def getmessage(i_d):
    headers = {'service':"chat", "id":i_d}

    url = "http+unix://" + localsocket+"/getmessage"

    #print(url)

    r = requests.post(url, headers=headers)
    return r.content


def sendmessage(address, service, message):
    headers = {'address':address}

    url = "http+unix://" + localsocket+"/sendmessage"

    data = bytes(service, encoding="utf-8") + b"\n" + message

    r = requests.post(url, data=data, headers=headers)




##### message manager #####

message_cache = {}

def update_messages():
    ids = getmessageids()
    for i_d in ids:
        if i_d not in message_cache:
            message_cache[i_d] = getmessage(i_d)
    return message_cache