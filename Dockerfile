FROM golang:1.14

WORKDIR /go/src/github.com/mmathys/acfts
COPY . .

RUN go get ./...