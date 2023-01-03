##@ System:

up: ## start project
	docker-compose up --build -d

down: ## stop project
	docker-compose down

watch: ## watch project
	docker-compose up --build 

state: ## show state
	docker-compose ps

logs: ## show last 100 lines of logs
	docker-compose logs --tail=100 $(ARGS)


