FROM base

WORKDIR ./server

RUN go build ./...
RUN go install

ENV ADDRESS none
ENV ADAPTER rpc
ENV TOPOLOGY dockerSimple
ENV INSTANCE 0

WORKDIR /go/src/github.com/mmathys/acfts
#CMD ["sh", "-c", "server -b --pprof -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]
CMD ["sh", "-c", "server -b -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json -i ${INSTANCE}"]