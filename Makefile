PROTO_PATH = common/proto
PROTO_FILES = auth.proto author.proto book.proto
SERVICES = auth_service author_service book_service

gen:
	@for service in $(SERVICES); do \
		for proto in $(PROTO_FILES); do \
			base=$$(basename $$proto .proto); \
			protoc \
				--proto_path=$(PROTO_PATH) \
				$(PROTO_PATH)/$$proto \
				--go_out=$$service/interface/grpc/genproto/$$base --go_opt=paths=source_relative \
				--go-grpc_out=$$service/interface/grpc/genproto/$$base --go-grpc_opt=paths=source_relative; \
		done; \
	done
