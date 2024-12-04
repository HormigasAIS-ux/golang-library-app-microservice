gen:
	@protoc \
		--proto_path=interface/grpc/proto \
		interface/grpc/proto/auth.proto \
		--go_out=interface/grpc/genproto/auth --go_opt=paths=source_relative \
		--go-grpc_out=interface/grpc/genproto/auth --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=interface/grpc/proto \
		interface/grpc/proto/author.proto \
		--go_out=interface/grpc/genproto/author --go_opt=paths=source_relative \
		--go-grpc_out=interface/grpc/genproto/author --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=interface/grpc/proto \
		interface/grpc/proto/book.proto \
		--go_out=interface/grpc/genproto/book --go_opt=paths=source_relative \
		--go-grpc_out=interface/grpc/genproto/book --go-grpc_opt=paths=source_relative

