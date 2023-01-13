##@ Develop:

sh: ## start shell in backend
	docker-compose exec mysqldump-slice sh

fmt: ## run gofmt
	gofmt -w .

build: ## build slice for centos
	sudo rm target/slice
	docker run --rm -v ${PWD}:/mysqldump-slice -w /mysqldump-slice lunny/centos-go:latest go build -o target/slice cmd/slice.go
	sudo chown ${USER}:${USER} target/slice

