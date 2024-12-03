gen:
	@protoc \
		--proto_path=interface/grpc \
		interface/grpc/auth_service.proto \
		--go_out=interface/grpc/genproto --go_opt=paths=source_relative \
		--go-grpc_out=interface/grpc/genproto --go-grpc_opt=paths=source_relative
