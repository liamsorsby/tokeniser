FROM golang:1.12 as builder

LABEL maintainer="Liam Sorsby <liam.sorsby@sky.com>"

WORKDIR $GOPATH/src/github.com/liamsorsby/tokeniser

COPY . .

RUN go get -d -v ./...

RUN go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

RUN go test && gofmt && golangci-lint run --color always

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o sts

FROM scratch

COPY --from=builder /go/src/github.com/liamsorsby/tokeniser /app/

WORKDIR /app

CMD ["./sts"]
