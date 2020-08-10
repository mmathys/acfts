FROM base

EXPOSE 6666

WORKDIR ./test/crypto/merkle

ENV TEST none
ENV BENCH none
ENV POOL_SIZE 32
ENV BENCHTIME 10s

CMD ["sh", "-c", "go test -run ${TEST}$ -bench ${BENCH}$ ./merkle_test.go -benchtime ${BENCHTIME} -args ${POOL_SIZE}"]