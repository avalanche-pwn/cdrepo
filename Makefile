proto_src := internal/daemon_pb/daemon.proto
proto_out := \
	internal/daemon_pb/daemon.pb.go \
	internal/daemon_pb/daemon_grpc.pb.go

build: generate
	go build ./cmd/cdrepo/cdrepo.go

generate: $(proto_out)

clean:
	rm ./cdrepo

$(proto_out): $(proto_src)
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		$(proto_src)
