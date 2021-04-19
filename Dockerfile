FROM golang:1.16.3-alpine3.13 as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/gorux/app/

COPY go.mod .
COPY config /dist/config
COPY docker-entrypoint.sh /dist/docker-entrypoint.sh

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' -a \
  -o /dist/bootstrap  .

FROM alpine

RUN apk update && apk add --no-cache ca-certificates bash curl tzdata && update-ca-certificates
ENV TZ=Asia/Chongqing
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /dist/bootstrap /dist/bootstrap
COPY --from=builder /dist/config /dist/config
COPY --from=builder /dist/docker-entrypoint.sh /dist/docker-entrypoint.sh

EXPOSE 8181

RUN chmod +x /dist/docker-entrypoint.sh

VOLUME ["/dist/config", "/dist/data"]

WORKDIR /dist/

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["run"]

