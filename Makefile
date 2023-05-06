.PHONY: deploy run ping compress

deploy: ping compress

run: ping
	cd bin&& ping.exe

ping:
	go build -trimpath -ldflags "-s -w" -o bin/ping.exe goping/cmd/ping

compress:
	upx -9 bin/ping.exe



