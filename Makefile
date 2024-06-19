.PHONY: tidy
tidy:
	go mod tidy

.PHONY: style
style:
	goimports -l -w ./api
	goimports -l -w ./broker
	goimports -l -w ./client
	goimports -l -w ./runtime
	goimports -l -w ./security
	goimports -l -w ./server
	goimports -l -w ./store
	goimports -l -w ./telemetry
	goimports -l -w ./utils

.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: proto-account
proto-account:
	protoc proto/account/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-health
proto-health:
	protoc proto/health/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-rule
proto-rule:
	protoc proto/rule/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-runtime
proto-runtime:
	protoc proto/runtime/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-greeter
proto-greeter:
	protoc examples/greeter/proto/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: build-image
build-image:
	docker build -t github.com/w-h-a/pkg:v1.0.0 .

.PHONY: load-image
load-image:
	kind load docker-image github.com/w-h-a/pkg:v1.0.0

.PHONY: port-forward
port-forward:
	kubectl port-forward service/runtime 8080:8080

