FROM golang:1.21 AS go

FROM go AS linter
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

FROM go AS mockery
RUN go install github.com/vektra/mockery/v2@v2.40.3

FROM go AS easyjson
RUN go install github.com/mailru/easyjson/easyjson@v0.7.7

FROM go AS protoc
ARG TARGETARCH
RUN VERSION=26.0 && \
    apt-get update && \
    apt-get install unzip && \
    if [ "$TARGETARCH" = "arm64" ]; then ARCH="aarch_64"; else ARCH="x86_64"; fi; \
    wget -O x.zip https://github.com/protocolbuffers/protobuf/releases/download/v$VERSION/protoc-$VERSION-linux-$ARCH.zip && \
    unzip x.zip -d protoc

FROM go AS protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0

FROM go AS grpc
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

FROM go AS sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.25.0

FROM go AS validate
RUN go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

FROM go AS dispatcher
WORKDIR /app
ADD ./pkg/protoc-gen-dispatcher .
RUN go install .


FROM go AS dev
ENV INSIDE_DEV_CONTAINER 1
COPY --from=linter     /go/bin/golangci-lint        /usr/bin/
COPY --from=mockery    /go/bin/mockery               /usr/bin/
COPY --from=easyjson   /go/bin/easyjson              /usr/bin/
COPY --from=protoc     /go/protoc/                   /usr/
COPY --from=protobuf   /go/bin/protoc-gen-go         /usr/bin/
COPY --from=grpc       /go/bin/protoc-gen-go-grpc    /usr/bin/
COPY --from=dispatcher /go/bin/protoc-gen-dispatcher /usr/bin/
COPY --from=sqlc       /go/bin/sqlc                  /usr/bin/
COPY --from=validate   /go/bin/protoc-gen-validate   /usr/bin/
WORKDIR /app
