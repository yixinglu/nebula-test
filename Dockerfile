FROM golang:1.13.2-alpine as builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
ENV NEBULA_TEST /home/nebula-test

COPY . ${NEBULA_TEST}

WORKDIR ${NEBULA_TEST}

RUN go build -o target/nebula-test . \
  && mkdir -p /usr/local/nebula/bin/ \
  && cp target/nebula-test /usr/local/nebula/bin/nebula-test

FROM golang:1.13.2-alpine

ENV NEBULA_HOME /usr/local/nebula/

COPY --from=builder ${NEBULA_HOME} ${NEBULA_HOME}

WORKDIR ${NEBULA_HOME}

ENV PATH ${NEBULA_HOME}/bin:${PATH}

ENTRYPOINT ["nebula-test"]