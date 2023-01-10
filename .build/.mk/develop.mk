##@ Develop:

sh: ## start shell in backend
	docker-compose exec mysqldump-slice sh

build: ## build
	docker run --rm -v ${PWD}:/mysqldump-slice -w /mysqldump-slice lunny/centos-go:latest go build -o cmd/slice cmd/main.go
