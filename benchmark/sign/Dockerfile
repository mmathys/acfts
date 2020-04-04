FROM base

WORKDIR ./benchmark/sign

ENV TEST BenchmarkSignNoNetwork
ENV TOPOLOGY dockerSimple

CMD ["sh", "-c", "go test -bench BenchmarkSignNoNetwork ./sign_test.go -args /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]