FROM base

WORKDIR ./server

RUN go build ./...
RUN go install

ENV ADDRESS 04d34e04d720691c6d392cfc49d59d501e813ab6d879a1aba563a26b76c6f3109791c7aa8297fca7d399b068ff1baae0c9587f371ca04ccadded614a49e474cc25
ENV ADAPTER rpc
ENV TOPOLOGY dockerSimple
ENV INSTANCE 0

WORKDIR /go/src/github.com/mmathys/acfts
#CMD ["sh", "-c", "server -b --pprof -a ${ADDRESS} --adapter ${ADAPTER} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]
CMD ["sh", "-c", "server -b -a ${ADDRESS} --adapter ${ADAPTER} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json -i ${INSTANCE}"]