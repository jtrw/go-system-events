.PHONY: install
install: ##up-services ## install up environment
	docker-compose up -d --build

.PHONY: start
start: ##up-services ## spin up environment
	docker-compose up -d

.PHONY: stop
stop: ## stop environment
	docker-compose stop
