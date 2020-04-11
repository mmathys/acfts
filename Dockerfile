FROM golang:1.14

WORKDIR /go/src/github.com/mmathys/acfts

COPY go.mod go.sum ./

RUN go mod download

COPY . .