FROM golang AS builder

WORKDIR /build

COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build

FROM alpine AS runtime

COPY --from=builder /build/mxget /usr/local/bin/

CMD ["mxget", "serve"]

EXPOSE 8080
