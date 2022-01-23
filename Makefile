build:
	go build -o ./install_files/peer_daemon `find . -wholename './src/*.go' -a ! -wholename './src/*env.go'` ./src/dev_env.go
windows:
	GOOS=windows GOARCH=amd64 go build -o install_files/peer_daemon_win.exe ./src/*.go
linux_64:
	GOOS=linux GOARCH=amd64 go build -o ./install_files/peerd `find . -wholename './src/*.go' -a ! -wholename './src/*env.go'` ./src/linux_env.go