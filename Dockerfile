FROM golang:1.14.2-alpine3.11

RUN apk add --update --no-cache \
    --repository http://dl-3.alpinelinux.org/alpine/edge/community \
    --repository http://dl-3.alpinelinux.org/alpine/edge/main \
    vips-dev
RUN apk add build-base

WORKDIR /go/src/project/

COPY . /go/src/project/
RUN go build -o /bin/go-app

WORKDIR /app
CMD ["/bin/go-app"]
