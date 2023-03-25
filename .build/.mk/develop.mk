##@ Develop:

sh: ## start shell in backend [make sh]
	docker-compose exec mysqldump-slice sh

fmt: ## run gofmt [make fmt]
	gofmt -w .

test: ## run go test [make test]
	go test ./repository/ ./entity/ -v

dev: ## run development watching [make dev]
	watch go test ./service/ -v

clear: ## clear target direct [make clear] 
	sudo rm -r target/*

build: ## build slice [make build]
	make clear
	go build -o target/slice cmd/slice.go

build-centos: ## build slice for centos [make build-centos]
	make clear
	docker run --rm -v ${PWD}:/mysqldump-slice -w /mysqldump-slice lunny/centos-go:latest go build -o target/slice cmd/slice.go
	sudo chown ${USER}:${USER} target/slice

run: ## run slice [make run]
	make build
	target/slice ./conf.yaml
