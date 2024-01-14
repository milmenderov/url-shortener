.PHONY: re down run all
re:
	docker-compose -f ./.docker/docker-compose.yml up --build -d
down:
	docker-compose -f ./.docker/docker-compose.yml down -v
run:
	docker-compose -f ./.docker/docker-compose.yml up -d
all:	re
