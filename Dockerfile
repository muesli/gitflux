FROM golang:1.15-alpine
WORKDIR /go/src/github.com/muesli/gitflux

# Download dependencies on a different layer
COPY go.* ./
RUN go mod download

# Copy and compile code
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /gitflux .

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /gitflux /gitflux

USER 1000:1000
ENTRYPOINT [ "/gitflux" ]
CMD [ "--help" ]
