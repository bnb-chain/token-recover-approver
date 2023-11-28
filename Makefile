.PHONY: tools run generate build build-image
VERSION=$(shell git describe --tags --always)
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_COMMIT_DATE=$(shell git log -n1 --pretty='format:%cd' --date=format:'%Y%m%d')
REPO=github.com/bnb-chain/token-recover-approver

ldflags = -X $(REPO)/internal/version.AppVersion=$(VERSION) \
          -X $(REPO)/internal/version.GitCommit=$(GIT_COMMIT) \
          -X $(REPO)/internal/version.GitCommitDate=$(GIT_COMMIT_DATE)

tools:
	@go install github.com/google/wire/cmd/wire@v0.5.0

generate: 
	@go generate main.go

run:
	@go run -ldflags="$(ldflags)" main.go --config configs/default.config.yaml

build: generate
	go build -ldflags="$(ldflags)" -o ./build/bin/approver main.go

build-image:
	@read -p "Enter Image Name: " IMAGE_NAME; \
	docker build . -f ./Dockerfile -t "$$IMAGE_NAME"
