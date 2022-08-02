postgres:
	docker start dev-postgres

test: 
	go test -v -cover ./...

.PHONY: postgres test