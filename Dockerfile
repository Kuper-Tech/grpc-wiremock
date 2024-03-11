ARG GOLANG_IMAGE_REPO
ARG GOLANG_IMAGE_NAME
ARG GOLANG_IMAGE_TAG

ARG WIREMOCK_IMAGE_TAG
ARG WIREMOCK_IMAGE_REPO
ARG WIREMOCK_IMAGE_NAME

FROM ${GOLANG_IMAGE_REPO}/${GOLANG_IMAGE_NAME}:${GOLANG_IMAGE_TAG} as golang

ENV CGO_ENABLED=0
ENV GOBIN="/go/bin"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
  && go install github.com/githubnemo/CompileDaemon@v1.2.1 \
  && go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.8

COPY . /code
RUN make install-cli -C /code

FROM ${GOLANG_IMAGE_REPO}/${GOLANG_IMAGE_NAME}:${GOLANG_IMAGE_TAG} as gopath

ENV CGO_ENABLED=0
ENV GOPATH="/go"
ENV GOCACHE="/go/go-build"

# Warm up go mod & build cache
COPY ./example /tmp/gocache
RUN make install -C /tmp/gocache && \
    rm -rf /tmp/gocache

FROM ${WIREMOCK_IMAGE_REPO}/${WIREMOCK_IMAGE_NAME}:${WIREMOCK_IMAGE_TAG}

# Install tools
RUN apk add --no-cache \
    sudo \
    iptables \
    rsyslog \
    jq \
    gettext \
    make \
    protobuf-dev \
    lsof \
    perl \
    curl \
    nginx \
   	yq \
    tree \
    git


# Create the user
ARG USERNAME=mock
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG HOME=/home/$USERNAME

RUN addgroup -g $USER_GID $USERNAME && \
    adduser -u $USER_UID -G $USERNAME --disabled-password -h /home/$USERNAME $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

# Go is required to generate a proxy
ENV CGO_ENABLED=0
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV GOBIN="/go/bin"
ENV GOCACHE="/go/go-build"
ENV PATH="${GOBIN}:${GOROOT}/bin:/scripts:${PATH}"

## Go root
COPY --from=golang ${GOROOT} ${GOROOT}

# Go mod & build cache
COPY --chown=${USERNAME}:${USERNAME} --from=gopath ${GOPATH} ${GOPATH}
## Assuming that GOCACHE is included into GOPATH
#COPY --chown=${USERNAME}:${USERNAME} --from=gopath ${GOCACHE} ${GOCACHE}

# Go installed binaries
COPY --chown=${USERNAME}:${USERNAME} --from=golang ${GOBIN} ${GOBIN}

COPY --from=ochinchina/supervisord:latest /usr/local/bin/supervisord /usr/local/bin/supervisord

# Proxy
ARG PROXY_DIR="/var/proxy"
RUN mkdir ${PROXY_DIR} && \
    chown -R "${USER_UID}:${USER_UID}" ${PROXY_DIR}

# Logging
COPY ./etc/rsyslog.d/wiremock.conf /etc/rsyslog.d/wiremock.conf

# Scripts
COPY scripts /scripts
RUN chmod -R +x /scripts

# Nginx
ARG NGINX_DIR="/etc/nginx/http.d"
COPY etc/nginx/http.d ${NGINX_DIR}
RUN chown -R "${USER_UID}:${USER_UID}" ${NGINX_DIR}

# Supervisord
ARG SUPERVISORD_DIR="/etc/supervisord"
COPY etc/supervisord ${SUPERVISORD_DIR}
RUN chown -R "${USER_UID}:${USER_UID}" ${SUPERVISORD_DIR}


USER ${USERNAME}
WORKDIR ${HOME}

LABEL org.opencontainers.image.source="https://github.com/SberMarket-Tech/grpc-wiremock"
LABEL org.opencontainers.image.description="WireMock for multiple APIs with support of gRPC and HTTP"
LABEL org.opencontainers.image.licenses="Apache 2.0"

ENTRYPOINT ["/scripts/init.sh"]
