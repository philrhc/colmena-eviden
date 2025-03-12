#!/bin/sh
go get github.com/swaggo/swag/cmd/swag@latest \
  && go get github.com/swaggo/http-swagger \
  && go install github.com/swaggo/swag/cmd/swag@latest \
  && go install github.com/swaggo/http-swagger \
  && PATH=$(go env GOPATH)/bin:$PATH/ \
  && git add . \
  && git commit -m "before" \
  && swag init -g cmd/main.go -o docs --parseDependency \
  && git diff --exit-code \
  && go mod tidy