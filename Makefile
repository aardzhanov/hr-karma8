build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o serverB ./serverB/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o serverA ./serverA/
clean:
	rm ./serverA/serverA
	rm ./serverB/serverB
