##@ System:

up: ## up project [make up] [make up build=1 watch=1]
ifdef build
	$(eval OPTS=${OPTS} --build)
endif
ifdef watch
else
	$(eval OPTS=${OPTS} -d)
endif
	docker-compose up ${OPTS} --remove-orphans

down: ## stop project [make down]
	docker-compose down

ps: ## show state [make ps]
	docker-compose ps

logs: ## show last 100 lines of logs [make logs]
	docker-compose logs --tail=100 $(ARGS)


