server_bin: server/server.go
	cd server; go build server.go; cd ..

main_bin: main.go
	go build main.go

run: main_bin server_bin
	./main

clean: 
	rm main server/server 

