export GO111MODULE=on
BUILD_DATE=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_TAG=1.0

.PHONY: build
build: clean
	go build -ldflags "-X main.GitTag=$(GIT_TAG) -X main.BuildTime=$(BUILD_DATE)" ./check_api.go && \
  mkdir -p build/bin && mv check_api build/bin/kcheckApi && \
  mkdir -p build/conf && cp conf/* build/conf/ && \
  cp deploy_kcheck.sh build/ && chmod +x build/deploy_kcheck.sh

.PHONY: clean
clean:
	rm -rf build;