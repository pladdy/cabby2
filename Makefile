.PHONY: all build build/debian/usr/bin/cabby build/debian/usr/bin/cabby-cli build/debian/var/cabby/schema.sql
.PHONY: clean clean-cli cmd/cabby-cli/cabby-cli config cover cover-html db/cabby.db reportcard run run-log
.PHONY: test test-failures test test-run

BUILD_TAGS=-tags json1
BUILD_PATH=build/cabby
CLI_FILES=$(shell find cmd/cabby-cli/*.go -name '*go' | grep -v test)
PACKAGES=./ sqlite/... http/... cmd/cabby-cli/...

all: config cert dependencies

build: build/debian/usr/bin/cabby build/debian/usr/bin/cabby-cli

build/debian/etc/cabby/:
	mkdir -p $@

build/debian/etc/cabby/cabby.json: build/debian/etc/cabby/
	cp config/cabby.json build/debian/etc/cabby/

build/debian/etc/systemd/:
	mkdir -p $@

build/debian/lib/systemd/system/cabby.service.d/:
	mkdir -p $@

build/debian/usr/bin/:
	mkdir -p $@

build/debian/usr/bin/cabby-cli: build/debian/usr/bin/
	go build -o $@ $(CLI_FILES)

build/debian/usr/bin/cabby: build/debian/usr/bin/
	go build $(BUILD_TAGS) -o $@ cmd/cabby/main.go

build/debian/var/cabby/:
	mkdir -p $@

build/debian/var/cabby/schema.sql: build/debian/var/cabby/
	cp sqlite/schema.sql $@

build-debian: config build/debian/etc/cabby/cabby.json build/debian/var/cabby/schema.sql
	vagrant up
	@echo Magic has happend to make a debian...
	vagrant destroy -f

clean:
	rm -rf db/
	rm -f server.key server.crt *.log cover.out config/cabby.json
	rm -f build/debian/usr/bin/cabby build/debian/usr/bin/cabby-cli

clean-cli:
	rm -f build/debian/usr/bin/cabby-cli

cert:
	openssl req -x509 -newkey rsa:4096 -nodes -keyout server.key -out server.crt -days 365 -subj "/C=US/O=Cabby TAXII 2.0/CN=pladdy"
	chmod 600 server.key

cmd/cabby-cli/cabby-cli:
	go build -o $@ $(CLI_FILES)

config:
	@for file in $(shell find config/*example.json -type f | sed 's/.example.json//'); do \
		cp $${file}.example.json $${file}.json; \
	done
	@echo Configs available in config/

cover:
ifdef pkg
	go test $(BUILD_TAGS) -i ./$(pkg)
	go test $(BUILD_TAGS) -v -coverprofile=$(pkg).out ./$(pkg)
	go tool cover -func=$(pkg).out
	rm $(pkg).out
else
	@for package in $(PACKAGES); do \
	  go test $(BUILD_TAGS) -i ./$${package}; \
		go test $(BUILD_TAGS) -v -coverprofile=$${package}.out ./$${package}; \
		go tool cover -func=$${package}.out; \
		rm $${package}.out; \
	done
endif

cover-html:
ifdef pkg
	go test $(BUILD_TAGS) -i ./$(pkg)
	go test $(BUILD_TAGS) -v -coverprofile=$(pkg).out ./$(pkg)
	go tool cover -func=$(pkg).out
	go tool cover -html=$(pkg).out
	rm $(pkg).out
else
	@for package in $(PACKAGES); do \
	  go test $(BUILD_TAGS) -i ./$${package}; \
		go test $(BUILD_TAGS) -v -coverprofile=$${package}.out ./$${package}; \
		go tool cover -func=$${package}.out; \
		go tool cover -html=$${package}.out; \
		rm $${package}.out; \
	done
endif

cover-cabby.txt:
	go test -v $(BUILD_TAGS) -coverprofile=$@ -covermode=atomic ./

cover-http.txt:
	go test -v $(BUILD_TAGS) -coverprofile=$@ -covermode=atomic ./http/...

cover-sqlite.txt:
	go test -v $(BUILD_TAGS) -coverprofile=$@ -covermode=atomic ./sqlite/...

coverage.txt: cover-cabby.txt cover-http.txt cover-sqlite.txt
	@cat cover-cabby.txt cover-http.txt cover-sqlite.txt > $@
	@rm -f cover-cabby.txt cover-http.txt cover-sqlite.txt

db/cabby.db: cmd/cabby-cli/cabby-cli
	scripts/setup-cabby

dependencies:
	go get -t -v  ./...
	go get github.com/fzipp/gocyclo
	go get github.com/golang/lint

dev-db: db/cabby.db

fmt:
	go fmt -x

reportcard: fmt
	gocyclo -over 10 .
	golint
	go vet

run:
	go run $(BUILD_TAGS) cmd/cabby/main.go -config config/cabby.json

run-cli:
	go run $(CLI_FILES)

run-log:
	go run $(BUILD_TAGS) cmd/cabby/main.go 2>&1 | tee cabby.log

test:
ifdef pkg
	go test $(BUILD_TAGS) -i ./$(pkg)
	go test $(BUILD_TAGS) -v -cover ./$(pkg)
else
	go test $(BUILD_TAGS) -i ./...
	go test $(BUILD_TAGS) -v -cover ./...
endif

test-failures:
	go test $(BUILD_TAGS) -v ./... 2>&1 | grep -A 1 FAIL

test-run:
ifdef test
	go test $(BUILD_TAGS) -i ./...
	go test $(BUILD_TAGS) -v ./... -run $(test)
else
	@echo Syntax is 'make $@ test=<test name>'
endif
