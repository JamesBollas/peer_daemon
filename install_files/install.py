import os

in_dir = os.getcwd()

print('Generating torrc.')

socks_socket = os.path.join(in_dir,"torsocks.socket")

torrc = "SocksPort unix:" + socks_socket + "\n"

hidden_service_dir = os.path.join(in_dir, "hidden_service")

torrc += "HiddenServiceDir " + hidden_service_dir + "\n"

hidden_service_socket = os.path.join(in_dir, "hidden_service.socket")

torrc += "HiddenServicePort 80 unix:" + hidden_service_socket + "\n"

with open("peer_daemon.torrc", "w") as f:
	f.write(torrc)

print('Generating env')

local_socket = os.path.join(in_dir,"local.socket")
proxy_socket = os.path.join(in_dir,"torsocks.socket")
hostname = os.path.join(hidden_service_dir, "hostname")
proxy_executable = "/usr/bin/tor"
proxy_config = os.path.join(in_dir, "peer_daemon.torrc")

env = 'LOCAL_SOCKET="' + local_socket + '"\n'
env += 'PROXY_SOCKET="' + proxy_socket + '"\n'
env += 'HIDDEN_SERVICE_SOCKET="' + hidden_service_socket + '"\n'
env += 'HOSTNAME_PATH="' + hostname + '"\n'
env += 'PROXY_EXECUTABLE="' + proxy_executable + '"\n'
env += 'PROXY_CONFIG="' + proxy_config + '"\n'

with open(".env", "w") as f:
	f.write(env)