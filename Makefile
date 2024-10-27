run-exercise:
	@go run exercise-service/cmd/main.go

run-kitchen:
	@go run services/kitchen/*.go

gen-exercise:
	@protoc \
		--proto_path=proto "proto/exercises/exercise.proto" \
		--go_out=exercise-service/internal/genproto/ --go_opt=paths=source_relative \
  	--go-grpc_out=exercise-service/internal/genproto/ --go-grpc_opt=paths=source_relative