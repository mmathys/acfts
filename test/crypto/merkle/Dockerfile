FROM base

EXPOSE 6666

WORKDIR ./test/crypto/merkle

ENV TEST none
ENV BENCH none

CMD ["sh", "-c", "go test -run ${TEST}$ -bench ${BENCH}$ ./merkle_test.go"]