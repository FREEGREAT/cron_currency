migrate-up:
	goose -dir ./migration postgres "postgres://user:user@localhost:5430/subscribe_db" up

migrate-down:
	goose -dir ./migration postgres "postgres://user:user@localhost:5430/subscribe_db" down