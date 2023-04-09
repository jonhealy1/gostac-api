.PHONY: database
database:
	docker-compose up database elasticsearch

.PHONY: api
api:
	docker-compose up api

.PHONY: down
down:
	docker-compose down

psql-shell:		## Enter psql shell
	docker-compose run --rm database \
	psql -h database -U username -d postgis

.PHONY: test
test:
	go clean -testcache
	go test github.com/jonhealy1/goapi-stac/tests

