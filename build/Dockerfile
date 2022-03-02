FROM onosproject/golang-build:v1 as build

ENV ONOS_GRPC_CLIENT_ROOT=/go/src/github.com/onosproject/grpc-client
ENV CGO_ENABLED=0

RUN mkdir -p $ONOS_GRPC_CLIENT_ROOT/

COPY . $ONOS_GRPC_CLIENT_ROOT

RUN cd $ONOS_GRPC_CLIENT_ROOT && GO111MODULE=on go build -o /go/bin/main ./cmd/


FROM alpine:3.11
RUN apk add bash
ENV HOME=/home/grpc-client
RUN mkdir $HOME
WORKDIR $HOME

COPY --from=build /go/bin/main /usr/local/bin/

CMD ["main"]