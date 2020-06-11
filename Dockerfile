FROM golang:1.14 AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o app .

# Install migrate tool
WORKDIR ./
RUN git clone https://github.com/golang-migrate/migrate.git
WORKDIR migrate/cmd/migrate
RUN git tag -l
RUN git checkout v4.6.2
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'postgres' -ldflags="-X main.Version=$(git describe --tags)" -a -o /pgmigrate .

FROM alpine:3.10.2
RUN apk --no-cache add ca-certificates
COPY --from=builder /app /pgmigrate ./
COPY --from=builder /app/docker-entrypoint.sh ./
RUN chmod +x ./docker-entrypoint.sh
RUN chmod +x ./app
RUN chmod +x ./pgmigrate
EXPOSE 50051
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["app"]