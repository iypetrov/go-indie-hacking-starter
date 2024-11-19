build-prod:
	@templ generate
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags prod -o bin/main main.go static_prod.go

run-local-mac:
	@./bin/tailwindcss-extra-macos-x64 -i ./static/css/input.css -o ./static/css/output.css
	@templ generate
	@air -c .air.toml

run-local-linux:
	@./bin/tailwindcss-extra-linux-x64 -i ./static/css/input.css -o ./static/css/output.css
	@templ generate
	@air -c .air.toml

fmt:
	@go fmt ./...
	@goimports -l -w .
	@templ fmt .
	@find . -name '*.sql' -exec pg_format -i {} +

update-deps:
	@go mod tidy
	@curl -sL https://github.com/dobicinaitis/tailwind-cli-extra/releases/latest/download/tailwindcss-extra-linux-x64 -o bin/tailwindcss-extra-linux-x64
	@curl -sL https://github.com/dobicinaitis/tailwind-cli-extra/releases/latest/download/tailwindcss-extra-macos-x64 -o bin/tailwindcss-extra-macos-x64 
	@curl -sL https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js -o static/js/htmx.min.js
	@curl -sL https://cdn.jsdelivr.net/npm/alpinejs@3.14.3/dist/cdn.min.js -o static/js/alpine.min.js
