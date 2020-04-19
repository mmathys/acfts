FROM base

EXPOSE 6666

WORKDIR ./benchmark/map

ENV TEST none
ENV BENCH none
ENV NUM_WORKERS 10

CMD ["sh", "-c", "go test -run ${TEST} -bench ${BENCH} ./map_test.go -args ${NUM_WORKERS}"]