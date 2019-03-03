relay_runner: main.go
	GOOS=linux GOARCH=arm GOARM=5 go build

.PHONY: copy
copy: relay_runner
	ssh lime sudo systemctl stop relay_runner
	scp relay_runner lime:/home/pi/
	ssh lime sudo systemctl start relay_runner
