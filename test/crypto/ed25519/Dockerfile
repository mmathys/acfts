FROM base

EXPOSE 6666

WORKDIR ./test/crypto/ed25519

ENV TEST none
ENV BENCH none

CMD ["sh", "-c", "go test -run ${TEST}$ -bench ${BENCH}$ ./ed25519_test.go"]