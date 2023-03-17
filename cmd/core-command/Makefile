.PHONY: build run

build:
	cd ../../; make command;

run:
	export EDGEX_SECURITY_SECRET_STORE=false; ./core-command -cp=consul.http://localhost:8500 -registry
