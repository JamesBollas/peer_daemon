build: main.go
	go build -o ./install_files/peer_daemon *.go
windows: main.go
	GOOS=windows GOARCH=amd64 go build -o install_files/peer_daemon_win.exe *.go
