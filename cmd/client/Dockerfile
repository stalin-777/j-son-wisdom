ARG GO_VER=1.23
ARG ALPINE_VER=3.19
FROM golang:${GO_VER}-alpine${ALPINE_VER} AS builder

WORKDIR /src
COPY ./ ./

RUN go mod download
RUN go build -o app cmd/client/main.go

FROM alpine:${ALPINE_VER}

RUN addgroup -g 101 app && \
    adduser -H -u 101 -G app -s /bin/sh -D app

WORKDIR /app/

EXPOSE 8080

COPY --from=builder --chown=app:app /src/app .

ENV NETWORK=tcp
ENV ADDRESS=localhost:8080

USER app

SHELL ["/bin/sh", "-c"]
ENTRYPOINT /app/app --network "$NETWORK" --address "$ADDRESS"

