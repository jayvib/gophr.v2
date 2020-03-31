# Copyrighted 2020. Jayson Vibandor
#
# Variables
CLIENT_APP="gophr.client"
API_BIN="gophr.engine"
APPNAME=gophr

unit-test:
	@go test -tags=unit -covermode=atomic -short ./... | grep -v '^?'

build-api: mod
	@echo "Building ${APPNAME}"
	if [ ! -e ./bin ]; then mkdir ./bin; fi
	go build -o ./bin/${API_BIN} ./cmd/gophr/

build:
	docker build -t ${APPNAME} --file ./deployment/gophr/Dockerfile .

up: build
	docker-compose -f  docker-compose.yaml up -d

down:
	docker-compose -f  docker-compose.yaml down -d

################UTILITY#################
mod: ## To download the dependency of the app
	go mod download

clean:
	sudo docker image prune -af

# To give description of the target
targets: ## To give description of the target
	@echo '#####Makefile Targets######'
	@grep '^[^#[:space:]].*:' Makefile | sed 's/:.*//g'

lint: ## Lints the project source code
	golint $(shell go list ./... | grep -v /vendor/)

fmt: ## Format source files excluding the vendor directory
	@go fmt -x ./...

help: ## Display the available targets and its description
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#####THIRD-PARTY TOOL INSTALLATION####
install-tools:
	@echo Installing mockery
	@go get github.com/vektra/mockery/cmd/mockery
	@echo Installing golint
	@go get -u golang.org/x/lint/golint

.PHONY:
	mod targets build-client clean targets lint
	fmt help build-client-step build-client 
	build-client-docker run-client install-tools
