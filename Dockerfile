FROM golang:buster AS builder
WORKDIR /src

COPY go.* ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg
RUN go build -o server ./cmd/server

FROM debian:buster
ENV ASSETS_PATH=/usr/share/sso

COPY --from=builder /src/server /usr/local/bin/
COPY assets "${ASSETS_PATH}"
COPY .docker/docker-entrypoint.sh /

ENTRYPOINT [ "/docker-entrypoint.sh" ]
CMD [ "server" ]
