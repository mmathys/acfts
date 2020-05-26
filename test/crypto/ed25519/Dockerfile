FROM base

EXPOSE 6666

WORKDIR ./test/crypto/ed25519

ENV TEST none
ENV BENCH none
ENV NUM_WORKERS 10

CMD ["sh", "-c", "go test -run ${TEST}$ -bench ${BENCH}$ ./batch_verify_test.go -args ${NUM_WORKERS}"]