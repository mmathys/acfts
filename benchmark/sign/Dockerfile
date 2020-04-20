FROM base

EXPOSE 6666

WORKDIR ./benchmark/sign

ENV TEST none
ENV BENCH none
ENV TOPOLOGY dockerSimple
ENV NUM_WORKERS 10

CMD ["sh", "-c", "go -test.run ${TEST} -test.bench ${BENCH} ./sign_test.go -args /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json ${NUM_WORKERS}"]