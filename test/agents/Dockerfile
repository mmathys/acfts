FROM base

WORKDIR ./test/agents

ENV TOPOLOGY dockerSimple
ENV NUM_WORKERS 11

CMD ["sh", "-c", "go test ./agent_test.go -run TestAgentsRPCFixed -args /go/src/github.com/mmathys/acfts/topologies/${TOPOLOGY}.json ${NUM_WORKERS}"]