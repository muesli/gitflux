FROM alpine
COPY gitflux_*.apk /tmp/
RUN apk add --allow-untrusted /tmp/gitflux_*.apk
ENTRYPOINT ["/usr/local/bin/gitflux"]
