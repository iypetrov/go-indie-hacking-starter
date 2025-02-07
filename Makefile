run-local-mac:
	@sqlc generate
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false & \
	air -c .air.toml & \
	./bin/tailwindcss-extra-macos-x64 -i ./static/css/input.css -o ./static/css/output.css --watch

run-local-linux:
	@sqlc generate
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false & \
	air -c .air.toml & \
	./bin/tailwindcss-extra-linux-x64 -i ./static/css/input.css -o ./static/css/output.css --watch

fmt:
	@go fmt ./...
	@goimports -l -w .
	@templ fmt .

update-deps:
	@curl -sL https://github.com/dobicinaitis/tailwind-cli-extra/releases/latest/download/tailwindcss-extra-linux-x64 -o bin/tailwindcss-extra-linux-x64
	@chmod +x bin/tailwindcss-extra-linux-x64
	@curl -sL https://github.com/dobicinaitis/tailwind-cli-extra/releases/latest/download/tailwindcss-extra-macos-x64 -o bin/tailwindcss-extra-macos-x64
	@chmod +x bin/tailwindcss-extra-macos-x64
	@curl -sL https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js -o static/js/htmx.min.js
	@curl -sL https://cdn.jsdelivr.net/npm/alpinejs@3.14.3/dist/cdn.min.js -o static/js/alpine.min.js
	@go get -u
	@go mod tidy
	@go mod vendor

vuln-check:
	@govulncheck -tags=prod -show verbose ./...


