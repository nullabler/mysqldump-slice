##@ Develop:

sh: ## start shell in backend
	docker-compose exec mysqldump-slice sh

build: ## build slice for centos
	docker run --rm -v ${PWD}:/mysqldump-slice -w /mysqldump-slice lunny/centos-go:latest go build -o target/slice cmd/main.go
