SRC_PATH:= ${PWD}

test:
	go test ./... -v
	
mock-repo:	
	charlatan -dir=${SRC_PATH}/appprj/domain/repo -output=${SRC_PATH}/appprj/domain/repo/mock/mock.go -package=mock UserRepo AccountRepo TransactionRepo