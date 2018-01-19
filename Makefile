LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.Version=$(shell git rev-parse HEAD)"

.PHONY: init
init:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/golang/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	@chmod +x ./hack/check.sh
	@chmod +x ./hooks/pre-commit

.PHONY: setup
setup: init
	git init
	@echo "Install pre-commit hook"
	@ln -s $(shell pwd)/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit || true
	dep init

.PHONY: check
check:
	@./hack/check.sh ${scope}

.PHONY: ci
ci: init
	@dep ensure
	@make check

.PHONY: build-master
build-master: check
	go build -ldflags '$(LDFLAGS)' -o ./binary/darwin/officerk-master ./cmd/master/

.PHONY: run-master
run-master:
	./binary/darwin/officerk-master -c ./conf/app.conf -d

.PHONY: build-node
build-node: check
	go build -ldflags '$(LDFLAGS)' -o ./binary/darwin/officerk-node ./cmd/node

.PHONY: run-node
run-node:
	./binary/darwin/officerk-node -c ./conf/app.conf -d

.PHONY: build-docker-master
build-docker-master:
	GOOS=linux go build -ldflags '$(LDFLAGS) -s -w' -o ./binary/linux/officerk-master ./cmd/master/
	docker build -t cosmtrek/officerk-master -f ./binary/Dockerfile.master .

.PHONY: build-docker-node
build-docker-node:
	GOOS=linux go build -ldflags '$(LDFLAGS) -s -w' -o ./binary/linux/officerk-node ./cmd/node
	docker build -t cosmtrek/officerk-node -f ./binary/Dockerfile.node .

.PHONY: run-docker-master
run-docker-master:
	docker run -it --rm -p 9392:9392 cosmtrek/officerk-master -d

.PHONY: run-docker-node
run-docker-node:
	docker run -it --rm cosmtrek/officerk-node -d
