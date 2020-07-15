SRC_PATH:= ${PWD}

test:
	go test ./... -v
	
mock-repo:	
	charlatan -dir=${SRC_PATH}/app/domain/repo -output=${SRC_PATH}/app/domain/repo/mock/mock.go -package=mock UserRepo AccountRepo TransactionRepo

docker-dev:
	docker-compose -p go-prj-skeleton -f ${SRC_PATH}/deployment/docker-compose-dev.yml up -d

migrate-tool:
	git clone --recursive https://github.com/golang-migrate/migrate.git ${GOPATH}/src/github.com/golang-migrate/migrate
	cd ${GOPATH}/src/github.com/golang-migrate/migrate/cmd/migrate;\
	git checkout v4.11.0;\
	CGO_ENABLED=0 go build -tags 'mysql' -ldflags="-X main.Version=$(git describe --tags)" -a -o ${GOPATH}/bin/mysqlmigrate .

migrate:
	mysqlmigrate -database mysql://admin:moneyforward@123@\(localhost:3306\)/asset_accounting?multiStatements=true -path db/migrations up
	
build:
	go build -o project ${SRC_PATH}/cmd/srv/...