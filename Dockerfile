FROM alpine:edge

RUN apk update && apk add \
    bash \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY crypto-ticks-downloader /bin/app
COPY config.yaml /etc/crypto.yaml


CMD ["/bin/app", "--config=/etc/crypto.yaml"]
