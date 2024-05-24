
kraken:
	@echo Running krakend
	krakend run -c gateway/krakend.json

run_user:
	@echo Running user_service
	go run user_service/cmd/app/main.go

run_book:
	@echo Running book_service
	go run book_service/cmd/app/main.go

lint:
	golangci-lint run