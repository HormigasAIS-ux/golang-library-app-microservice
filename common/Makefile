gen:
	# auth_service
	@protoc \
		--proto_path=proto \
		proto/auth.proto \
		--go_out=../auth_service/interface/grpc/genproto/auth --go_opt=paths=source_relative \
		--go-grpc_out=../auth_service/interface/grpc/genproto/auth --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/author.proto \
		--go_out=../auth_service/interface/grpc/genproto/author --go_opt=paths=source_relative \
		--go-grpc_out=../auth_service/interface/grpc/genproto/author --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/book.proto \
		--go_out=../auth_service/interface/grpc/genproto/book --go_opt=paths=source_relative \
		--go-grpc_out=../auth_service/interface/grpc/genproto/book --go-grpc_opt=paths=source_relative

	# author_service
	@protoc \
		--proto_path=proto \
		proto/auth.proto \
		--go_out=../author_service/interface/grpc/genproto/auth --go_opt=paths=source_relative \
		--go-grpc_out=../author_service/interface/grpc/genproto/auth --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/author.proto \
		--go_out=../author_service/interface/grpc/genproto/author --go_opt=paths=source_relative \
		--go-grpc_out=../author_service/interface/grpc/genproto/author --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/book.proto \
		--go_out=../author_service/interface/grpc/genproto/book --go_opt=paths=source_relative \
		--go-grpc_out=../author_service/interface/grpc/genproto/book --go-grpc_opt=paths=source_relative

	# book_service
	@protoc \
		--proto_path=proto \
		proto/auth.proto \
		--go_out=../book_service/interface/grpc/genproto/auth --go_opt=paths=source_relative \
		--go-grpc_out=../book_service/interface/grpc/genproto/auth --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/author.proto \
		--go_out=../book_service/interface/grpc/genproto/author --go_opt=paths=source_relative \
		--go-grpc_out=../book_service/interface/grpc/genproto/author --go-grpc_opt=paths=source_relative

	@protoc \
		--proto_path=proto \
		proto/book.proto \
		--go_out=../book_service/interface/grpc/genproto/book --go_opt=paths=source_relative \
		--go-grpc_out=../book_service/interface/grpc/genproto/book --go-grpc_opt=paths=source_relative
