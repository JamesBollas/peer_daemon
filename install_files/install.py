import os

in_dir = os.getcwd()

print('Generating torrc.')

socks_socket = os.path.join(in_dir,"torsocks.socket")

torrc = "SocksPort unix:" + socks_socket + "\n"

hidden_service_dir = os.path.join(in_dir, "hidden_service")

torrc += "HiddenServiceDir " + hidden_service_dir + "\n"

hidden_service_socket = os.path.join(in_dir, "hidden_service.socket")

torrc += "HiddenServicePort 80 unix:" + hidden_service_socket + "\n"

with open("yukon.torrc", "w") as f:
	f.write(torrc)
