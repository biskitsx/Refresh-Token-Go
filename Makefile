restart:
	docker compose down -v 
	docker compose up db -d

run-restart: restart
	go run main.go

run:
	go run main.go
