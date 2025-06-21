all: gensql genproto build

genproto:
	@protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		pb/users/v1/users.proto

gensql:
	@sqlc generate


build:
	@go build -o dist/run ./cmd/run
