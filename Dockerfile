FROM golang:1.23.2 AS build-stage
WORKDIR /app
COPY . .
# You can optimize build time by creating a custom image with these dependencies installed
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
RUN go install github.com/a-h/templ/cmd/templ@v0.3.819
RUN ./bin/tailwindcss-extra-linux-x64 -i ./static/css/input.css -o ./static/css/output.css --minify
RUN sqlc generate
RUN templ generate
RUN go mod tidy
RUN go mod vendor
# CGO_ENABLED=0 issue with github.com/mattn/go-sqlite3 occurs because the go-sqlite3 package requires the use of cgo to work
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags prod -o bin/main .

FROM gcr.io/distroless/base-debian12 AS run-stage
COPY --from=build-stage /app/sql/migrations /sql/migrations
COPY --from=build-stage /app/data /data
COPY --from=build-stage /app/bin/main /bin/main
CMD ["/bin/main"]
