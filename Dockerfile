FROM golang:1.21.6-alpine AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apk update && apk add --no-cache postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o avito_testcase ./cmd/main.go


CMD ["./avito_testcase"]
