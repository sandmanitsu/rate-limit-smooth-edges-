init:
	cd docker && docker-compose up -d

down:
	cd docker && docker-compose down

run:
	clear && go run main.go

test:
	clear &\
	bash -c "go run main.go" &\
	bash -c "k6 run k6-test.js"