# [Installing - gomplate documentation](https://docs.gomplate.ca/installing/)
ifeq ($(OS),Windows_NT)
	# OSNAME = WIN32
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		# OSNAME = LINUX
		export REDIS_MASTER_IP=$(shell ip -4 addr show eth0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')
	endif
	ifeq ($(UNAME_S),Darwin)
		# OSNAME = OSX
		export REDIS_MASTER_IP=$(shell ipconfig getifaddr en0)
	endif
endif

up:
	mkdir -p /tmp/data/sentinel
	cd docker; \
	docker-compose up -d
down:
	cd docker; \
	docker-compose down
cp:
	cd docker/sentinel; \
	gomplate -f sentinel.conf.tpl -o sentinel1.conf; \
	gomplate -f sentinel.conf.tpl -o sentinel2.conf; \
	gomplate -f sentinel.conf.tpl -o sentinel3.conf; 
try:
	@echo test.$(REDIS_MASTER_IP) 
	@echo 'Hello, {{ .Env.USER }}' | gomplate
	@echo 'Hello, {{ .Env.REDIS_MASTER_IP }}' | gomplate
	gomplate -f docker/sentinel/sentinel1.conf -o docker/sentinel/sentinel1.conf
test:
	go clean -testcache
	go test -v ./...
uml:
	mkdir -p doc
	goplantuml -recursive \
	  --show-implementations \
	  -aggregate-private-members -show-aggregations -show-compositions \
	  -show-connection-labels -show-aliases ./ \
	  > doc/class.puml
	cp -p doc/class.puml doc/class.md
	gsed -i '1i\```plantuml' doc/class.md
	echo '```' >> doc/class.md