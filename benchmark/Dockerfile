FROM base

WORKDIR ./benchmark

ENV TOPOLOGY localSimple
ENV NUM_WORKERS 11

CMD ["sh", "-c", "go test ./agent_test.go -run TestAgentsRPCFixed -args ../topologies/${TOPOLOGY}.json ${NUM_WORKERS}"]