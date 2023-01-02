
precommit: ensure format generate test check addlicense
	@echo "ready to commit"

ensure:
	go mod verify
	go mod vendor

format:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	@find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	@find . -type f -name '*.go' -not -path './vendor/*' -exec goimports -w "{}" +

generate:
	rm -rf mocks
	go generate ./...

test:
	go test -p=1 -cover -race $(shell go list ./... | grep -v /vendor/)

check: lint vet errcheck

vet:
	@go vet $(shell go list ./... | grep -v /vendor/)

lint:
	GO111MODULE=off go get golang.org/x/lint/golint
	@golint -min_confidence 1 $(shell go list ./... | grep -v /vendor/)

errcheck:
	GO111MODULE=off go get github.com/kisielk/errcheck
	@errcheck -ignore '(Close|Write|Fprint)' $(shell go list ./... | grep -v /vendor/)

addlicense:
	@addlicense -c "Benjamin Borbe" -y 2023 -l bsd ./*.go
