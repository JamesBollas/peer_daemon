import pdci
import time
import threading

def print_thread():
    printed = set()
    i = 0
    while True:
        i +=1 
        messages = pdci.update_messages()
        for key, message in messages.items():
            if key not in printed:
                print(message.decode("utf-8"))
                printed.add(key)
        time.sleep(.1)
        
print_thread = threading.Thread(target=print_thread)
print_thread.daemon = True
print_thread.start()

while True:
    print("\t\t\t",end="")
    message = input()
    pdci.sendmessage(pdci.my_hostname(), "chat", bytes(message,encoding="utf-8"))