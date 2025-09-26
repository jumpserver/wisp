FROM golang:1.24-bullseye AS stage-build
ARG TARGETARCH
ENV CGO_ENABLED=0
ENV GO111MODULE=on

WORKDIR /opt/wisp

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN make build

FROM debian:bullseye-slim

ARG TARGETARCH
ENV LANG=en_US.UTF-8
ARG APT_MIRROR=http://deb.debian.org
ARG DEPENDENCIES="                    \
        bash-completion               \
        jq                            \
        less                          \
        ca-certificates"

RUN set -ex \
    && sed -i "s@http://.*.debian.org@${APT_MIRROR}@g" /etc/apt/sources.list \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && apt-get update \
    && apt-get install -y --no-install-recommends ${DEPENDENCIES} \
    && apt-get clean all \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/wisp

COPY --from=stage-build /opt/wisp/build/wisp .

CMD ["./wisp"]