PKGS_PASS := $(shell go list ./... | grep -v 'example/\|connector/grpc/*')
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path ".git/*" -print0 | xargs -0)

checkLocal: vet lint misspell staticcheck cyclo const test

checkLocalWithTest: overalls vet lint misspell staticcheck cyclo const

overalls:
	@echo "overalls"
	overalls -project=github.com/aberic/lilydb -covermode=count -ignore='.git,_vendor'
	go tool cover -func=overalls.coverprofile

vet:
	@echo "vet"
	go vet $(PKGS_PASS)

lint:
	@echo "golint"
	golint -set_exit_status $(PKGS_PASS)

misspell:
	@echo "misspell"
	misspell -source=text -error $(GO_FILES)

staticcheck:
	@echo "staticcheck"
	staticcheck $(PKGS_PASS)

cyclo:
	@echo "gocyclo"
	gocyclo -over 15 $(GO_FILES)
	gocyclo -top 10 $(GO_FILES)

const:
	@echo "goconst"
	goconst $(PKGS_PASS)

veralls:
	@echo "goveralls"
	goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $(COVERALLS_TOKEN)

test:
	@echo "test"
	go test -v -cover $(PKGS_PASS)

