init:
	cd docker && docker-compose up -d

down:
	cd docker && docker-compose down

run:
	clear && go run main.go