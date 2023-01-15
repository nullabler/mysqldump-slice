##@ Develop:

sh: ## start shell in backend
	docker-compose exec mysqldump-slice sh

fmt: ## run gofmt
	gofmt -w .

test: ## run go test
	go test ./repository/ -v

clear: 
	sudo rm -r target/*

build: ## build slice 
	make clear
	go build -o target/slice cmd/slice.go

build-centos: ## build slice for centos
	make clear
	docker run --rm -v ${PWD}:/mysqldump-slice -w /mysqldump-slice lunny/centos-go:latest go build -o target/slice cmd/slice.go
	sudo chown ${USER}:${USER} target/slice

run: ## run slice
	make build
	target/slice ./conf.yaml
