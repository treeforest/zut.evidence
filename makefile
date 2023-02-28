gen:
	protoc --proto_path=api/proto api/proto/*.proto --go_out=:api/pb --go-grpc_out=:api/pb --grpc-gateway_out=:api/pb --openapiv2_out=:api/swagger

gateway:
	go build -o ./bin/gateway.exe ./cmd/gateway/main.go \
	&& cd ./bin \
	&& gateway.exe -port=8080 -etcdUrl="http://localhost:12379"

file:
	go build -o ./bin/file.exe ./cmd/server/file/main.go \
	&& cd ./bin \
	&& file.exe -httpPort=8081 -rpcPort=20006 -etcdUrl="http://localhost:12379"

account:
	go build -o ./bin/account.exe ./cmd/server/account/main.go \
	&& cd ./bin \
	&& account.exe -port=20001 -etcdUrl="http://localhost:12379"

wallet:
	go build -o ./bin/wallet.exe ./cmd/server/wallet/main.go \
	&& cd ./bin \
	&& wallet.exe -port=20002 -etcdUrl="http://localhost:12379"

didResolver:
	go build -o ./bin/did_resolver.exe ./cmd/server/did_resolver/main.go \
	&& cd ./bin \
	&& did_resolver.exe -port=20003 -etcdUrl="http://localhost:12379"

comet:
	go build -o ./bin/comet.exe ./cmd/server/comet/main.go \
	&& cd ./bin \
	&& comet.exe -port=20004 -wsPort=8082 -etcdUrl="http://localhost:12379"

logic:
	go run cmd/server/logic/main.go -port=20005 -etcdUrl="http://localhost:12379"

# 构建 linux 环境下的可执行文件
build:
	go build -o ./bin/gateway ./cmd/gateway/main.go \
	&& go build -o ./bin/account ./cmd/server/account/main.go \
	&& go build -o ./bin/comet ./cmd/server/comet/main.go \
	&& go build -o ./bin/did_resolver ./cmd/server/did_resolver/main.go \
	&& go build -o ./bin/file ./cmd/server/file/main.go \
	&& go build -o ./bin/logic ./cmd/server/logic/main.go \
	&& go build -o ./bin/wallet ./cmd/server/wallet/main.go \

