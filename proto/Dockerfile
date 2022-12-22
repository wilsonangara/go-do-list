ARG GO_VERSION
FROM golang:${GO_VERSION}-stretch

RUN apt-get update; \
        apt-get install unzip

# install protoc
ARG PROTOC_VERSION
RUN curl -o protoc.zip -L "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-aarch_64.zip" && \
        mkdir -p /root/cache/proto && \
        unzip protoc.zip -d /usr/local/proto && \
        chmod -R 755 /usr/local/proto && \
        chmod +x /usr/local/proto/bin/protoc && \
        rm -f protoc.zip
ENV PATH $PATH:"/usr/local/proto/bin"
RUN chmod -R 755 /usr/local/bin

# install protoc-gen-go
ARG PROTOC_GEN_GO_VERSION
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}

# install protoc-gen-twirp
ARG PROTOC_GEN_TWIRP_VERSION
RUN go install github.com/twitchtv/twirp/protoc-gen-twirp@v${PROTOC_GEN_TWIRP_VERSION}
