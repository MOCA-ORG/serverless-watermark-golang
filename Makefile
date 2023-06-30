.PHONY: build clean deploy gomodgen

build:
	export GO111MODULE=on
	go mod tidy
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/watermark watermark/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

deploy-prod: clean build
	sls deploy --verbose --stage prod
