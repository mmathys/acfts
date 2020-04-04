FROM base

WORKDIR ./benchmark/sign

ENV TOPOLOGY localSimple

CMD ["sh", "-c", "go test ./sign_test.go -run BenchmarkSignNoNetwork -args /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json"]