FROM base

EXPOSE 6666

WORKDIR ./test/sign

# if env variables TEST of BENCH were not set in docker-compose.yml, use these values
# as fallback values. These values do not match any test.
ENV TEST none
ENV BENCH none

CMD ["sh", "-c", "TOPOLOGY=/go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY_NAME}.json go test -test.run ${TEST} -test.bench ${BENCH} ./sign_test.go"]