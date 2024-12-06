gen:
	@protoc \
		--proto_path=proto \
		proto/auth.proto \
		--go_out=genproto/auth --go_opt=paths=source_relative \
		--go-grpc_out=genproto/auth --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/author.proto \
		--go_out=genproto/author --go_opt=paths=source_relative \
		--go-grpc_out=genproto/author --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/book.proto \
		--go_out=genproto/book --go_opt=paths=source_relative \
		--go-grpc_out=genproto/book --go-grpc_opt=paths=source_relative

