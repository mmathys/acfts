FROM base

RUN cd ./server && go get
RUN go install ./server

ENV ADDRESS 0x04d34e04d720691c6d392cfc49d59d501e813ab6d879a1aba563a26b76c6f3109791c7aa8297fca7d399b068ff1baae0c9587f371ca04ccadded614a49e474cc25
ENV ADAPTER rpc
ENV TOPOLOGY localSimple

CMD ["sh", "-c", "server -b -a ${ADDRESS} --adapter ${ADAPTER} --topology ./topologies/${TOPOLOGY}.json"]