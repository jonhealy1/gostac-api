.PHONY: database
database:
	docker-compose up database elasticsearch

.PHONY: api
api:
	docker-compose up pg-api es-api

.PHONY: msg
msg:
	docker-compose up kafka zookeeper
	
.PHONY: down
down:
	docker-compose down

psql-shell:		## Enter psql shell
	docker-compose run --rm database \
	psql -h database -U username -d postgis

.PHONY: test
test:
	go clean -testcache
	cd pg-api && go test github.com/jonhealy1/goapi-stac/pg-api/tests

