FROM base

WORKDIR ./benchmark/agents

ENV TOPOLOGY dockerSimple

CMD ["sh", "-c", "go test ./agent_test.go -run TestAgentsRPCFixed -args /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json ${NUM_WORKERS}"]