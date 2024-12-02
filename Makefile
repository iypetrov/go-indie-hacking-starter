build-prod:
	@sqlc generate
	@./bin/tailwindcss-extra-macos-x64 -i ./static/css/input.css -o ./static/css/output.css --minify
	@templ generate
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags prod -o bin/main .

run-local-mac:
	@docker volume inspect go_indie_hacking_starter_data > /dev/null 2>&1 || docker volume create go_indie_hacking_starter_data
	@docker ps -a --filter "name=^/go_indie_hacking_starter_postgres$$" --format "{{.Names}}" | grep -w go_indie_hacking_starter_postgres > /dev/null 2>&1 || \
	docker run -d --name go_indie_hacking_starter_postgres -e POSTGRES_DB=go-indie-hacking-starter -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -v go_indie_hacking_starter_data:/var/lib/postgresql/data -p 5432:5432 postgres:15
	@sqlc generate
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false & \
	air -c .air.toml & \
	./bin/tailwindcss-extra-macos-x64 -i ./static/css/input.css -o ./static/css/output.css --watch

run-local-linux:
	@docker volume inspect go_indie_hacking_starter_data > /dev/null 2>&1 || docker volume create go_indie_hacking_starter_data
	@docker ps -a --filter "name=^/go_indie_hacking_starter_postgres$$" --format "{{.Names}}" | grep -w go_indie_hacking_starter_postgres > /dev/null 2>&1 || \
	docker run -d --name go_indie_hacking_starter_postgres -e POSTGRES_DB=go-indie-hacking-starter -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -v go_indie_hacking_starter_data:/var/lib/postgresql/data -p 5432:5432 postgres:15
	@sqlc generate
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false & \
	air -c .air.toml & \
	./bin/tailwindcss-extra-linux-x64 -i ./static/css/input.css -o ./static/css/output.css --watch

fmt:
	@go fmt ./...
	@goimports -l -w .
	@templ fmt .
	@find . -name '*.sql' -exec pg_format -i {} +

update-deps:
	@curl -sL https://github.com/dobicinaitis/tailwind-cli-extra/releases/latest/download/tailwindcss-extra-linux-x64 -o bin/tailwindcss-extra-linux-x64
	@chmod +x bin/tailwindcss-extra-linux-x64 
	@curl -sL https://github.com/dobicinaitis/tailwind-cli-extra/releases/latest/download/tailwindcss-extra-macos-x64 -o bin/tailwindcss-extra-macos-x64 
	@chmod +x bin/tailwindcss-extra-macos-x64 
	@curl -sL https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js -o static/js/htmx.min.js
	@curl -sL https://cdn.jsdelivr.net/npm/alpinejs@3.14.3/dist/cdn.min.js -o static/js/alpine.min.js
