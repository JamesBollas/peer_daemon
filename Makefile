build:
	go build -o ./install_files/peer_daemon ./src/*.go
windows:
	GOOS=windows GOARCH=amd64 go build -o install_files/peer_daemon_win.exe ./src/*.go
