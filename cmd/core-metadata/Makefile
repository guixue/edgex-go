.PHONY: build run

build:
	cd ../../; make metadata;

run:
	export EDGEX_SECURITY_SECRET_STORE=false; ./core-metadata -cp=consul.http://localhost:8500 -registry
