FROM golang:1.23-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN go get -t ./...

RUN go install -v ./...


## ========== Production ==========
FROM alpine:3.19

RUN apk --no-cache add  libc6-compat tzdata gcompat ca-certificates && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" >  /etc/timezone && \
    apk del tzdata

# Copy app
WORKDIR /opt/owl/

COPY --from=builder /go/bin/ .

ENTRYPOINT ["./api"]
