build:
	go build -gcflags="-N -l" -o hfcli main.go

fmt:
	gofmt -l -s -w .