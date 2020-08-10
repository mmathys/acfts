FROM base

EXPOSE 6666

WORKDIR ./test/crypto/merkle

ENV TEST none
ENV BENCH none
ENV POOL_SIZE 32

CMD ["sh", "-c", "go test -run ${TEST}$ -bench ${BENCH}$ ./merkle_test.go -timeout 30s -args ${POOL_SIZE}"]