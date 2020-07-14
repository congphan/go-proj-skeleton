SRC_PATH:= ${PWD}

test:
	go test ./... -v
	
mock-repo:	
	charlatan -dir=${SRC_PATH}/app/domain/repo -output=${SRC_PATH}/app/domain/repo/mock/mock.go -package=mock UserRepo AccountRepo TransactionRepo

docker-dev:
	docker-compose -p go-prj-skeleton -f ${SRC_PATH}/deployment/docker-compose-dev.yml up -d
	
build:
	go build -o project ${SRC_PATH}/cmd/srv/...