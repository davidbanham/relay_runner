relay_runner: *.go frontend frontend/*
	GO111MODULE=on GOOS=linux GOARCH=arm GOARM=5 packr2 build

frontend: frontend/src/*
	cd frontend && elm make src/*.elm && cd ..

.PHONY: copy
copy: relay_runner
	ssh lime sudo systemctl stop relay_runner
	scp relay_runner lime:/home/pi/
	ssh lime sudo systemctl start relay_runner
