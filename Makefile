proto_src := internal/daemonpb/daemon.proto
proto_out := \
	internal/daemonpb/daemon.pb.go \
	internal/daemonpb/daemon_grpc.pb.go

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
