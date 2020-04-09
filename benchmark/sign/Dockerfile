FROM base

WORKDIR ./benchmark/sign

ENV TEST BenchmarkSignNoNetwork
ENV TOPOLOGY dockerSimple
ENV NUM_WORKERS 10

CMD ["sh", "-c", "go test -bench ${TEST} ./sign_test.go -args /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]