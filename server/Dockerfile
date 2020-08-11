FROM base

WORKDIR ./server

RUN go build ./...
RUN go install

ENV ADDRESS none
ENV ADAPTER rpc
ENV TOPOLOGY none
ENV INSTANCE 0
ENV POOLING false
ENV POOL_SIZE -1

WORKDIR /go/src/github.com/mmathys/acfts
#CMD ["sh", "-c", "server -b --pprof -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]
#CMD ["sh", "-c", "server -b -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json -i ${INSTANCE}"]
CMD ["sh", "-c", "server -b -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json --merkle-pooling=${POOLING} --merkle-pool-size=${POOL_SIZE}"]