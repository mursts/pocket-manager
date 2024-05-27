FROM golang:1.22.3-alpine as builder

WORKDIR /app

COPY cmd/ ./cmd/
COPY app.go ./
COPY pocket.go ./
COPY go.mod ./

RUN go mod tidy
RUN go build -o pocket ./cmd/pocket/main.go

FROM alpine:3.19.1

COPY --from=builder /app/pocket /bin/pocket
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /

ENV ZONEINFO=/zoneinfo.zip

RUN chmod +x /bin/pocket
CMD /bin/pocket

