# Copyrighted 2020. Jayson Vibandor
#
# Variables
CLIENT_APP="gophr.client"
CLIENT_DOCKERFILE="dockerfile.client"
API_APP="gophr.app"

mod: ## To download the dependency of the app
	go mod download

clean:
	sudo docker image prune -af

# To give description of the target
targets: ## To give description of the target
	@echo '#####Makefile Targets######'
	@grep '^[^#[:space:]].*:' Makefile | sed 's/:.*//g'

build-client-step: build-client-docker

build-client: mod ## Building executable file for the gophr client app
	@echo "Building ${CLIENT_APP}"
	if [ ! -e ./bin ]; then mkdir ./bin; fi

	go build -o ./bin/${CLIENT_APP} ./gophr.client/ 

build-client-docker:
	sudo docker build -t ${CLIENT_APP} --file ./deployment/dockerfile.client .

run-client:
	docker-compose up -d --no-deps -f ./deployment/docker-compose-client.yaml

.PHONY:
	mod targets build-client
