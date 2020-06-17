FROM golang:1.14 AS builder

# Install migrate tool
WORKDIR /go/src
RUN git clone https://github.com/golang-migrate/migrate.git
WORKDIR migrate/cmd/migrate
RUN git tag -l
RUN git checkout v4.6.2
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'postgres' -ldflags="-X main.Version=$(git describe --tags)" -a -o /pgmigrate .

COPY go.mod go.sum /go/src/project/
WORKDIR /go/src/project
RUN go mod download

COPY . /go/src/project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /project ./cmd/srv/...

FROM alpine:3.10.2
RUN apk --no-cache add ca-certificates
COPY --from=builder /project /pgmigrate /go/src/project/docker-entrypoint.sh ./
COPY --from=builder /go/src/project/db ./db/
RUN chmod +x ./docker-entrypoint.sh
RUN chmod +x ./project
RUN chmod +x ./pgmigrate
EXPOSE 50051
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["project"]