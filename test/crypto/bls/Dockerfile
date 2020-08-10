FROM base

EXPOSE 6666

WORKDIR ./test/crypto/bls

# if env variables TEST of BENCH were not set in docker-compose.yml, use these values
# as fallback values. These values do not match any test.
ENV TEST none
ENV BENCH none

CMD ["sh", "-c", "go test -test.run ${TEST} -test.bench ${BENCH} ./bls_test.go"]