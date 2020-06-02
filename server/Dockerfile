FROM base

WORKDIR ./server

RUN go build ./...
RUN go install

ENV ADDRESS a77000e6e7e77df500e39144e34e31f1715c98a5ecaaf34b8863f2d907c35ab4
ENV ADAPTER rpc
ENV TOPOLOGY dockerSimple
ENV INSTANCE 0

WORKDIR /go/src/github.com/mmathys/acfts
#CMD ["sh", "-c", "server -b --pprof -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]
CMD ["sh", "-c", "server -b -a ${ADDRESS} --topology /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json -i ${INSTANCE}"]