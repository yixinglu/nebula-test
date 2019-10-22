FROM golang:1.13.2-alpine as builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
ENV NEBULA_TEST /home/nebula-test

COPY . ${NEBULA_TEST}

WORKDIR ${NEBULA_TEST}

RUN go build -o target/nebula-test . \
  && cp target/nebula-test /usr/local/nebula-test

FROM alpine

COPY --from=builder /usr/local/nebula-test /usr/local/bin/nebula-test

ENTRYPOINT ["nebula-test"]
