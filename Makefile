#
.PHONY: default
default: build-server

setup:
	cd ui && npm install

build-ui:
	cd ui && npm run build

run-dev-ui:
	cd ui && npm run dev

run-dev: build-ui
	go build . && ./jwixtac ui -o 

build-server: setup build-ui
	go build .
