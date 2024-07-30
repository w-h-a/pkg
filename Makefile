.PHONY: tidy
tidy:
	go mod tidy

.PHONY: style
style:
	goimports -l -w ./api
	goimports -l -w ./client
	goimports -l -w ./runtime
	goimports -l -w ./security
	goimports -l -w ./server
	goimports -l -w ./store
	goimports -l -w ./streams
	goimports -l -w ./telemetry
	goimports -l -w ./utils

.PHONY: test
test:
	go clean -testcache && go test -v -race ./...

.PHONY: proto-health
proto-health:
	protoc proto/health/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-account
proto-account:
	protoc proto/account/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-rule
proto-rule:
	protoc proto/rule/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-snapshot
proto-snapshot:
	protoc proto/snapshot/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-restore
proto-restore:
	protoc proto/restore/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-runtime
proto-runtime:
	protoc proto/runtime/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-streams
proto-streams:
	protoc proto/streams/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-cron
proto-cron:
	protoc proto/cron/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-ticket
proto-ticket:
	protoc proto/ticket/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-order
proto-order:
	protoc proto/order/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-payment
proto-payment:
	protoc proto/payment/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-search
proto-search:
	protoc proto/search/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-rss
proto-rss:
	protoc proto/rss/*.proto --go_out=paths=source_relative:. --proto_path=.

.PHONY: proto-news
proto-news:
	protoc proto/news/*.proto --go_out=paths=source_relative:. --proto_path=.

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

