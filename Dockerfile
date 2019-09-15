FROM golang:1.13
LABEL maintainer="Tyler Smith <tylersmith.me@gmail.com>"
WORKDIR /go/src/github.com/tyler-smith/sync-check
COPY . .
RUN CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o /bin/sync-check .

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /bin/sync-check /bin/sync-check
ENTRYPOINT ["/bin/sync-check"]