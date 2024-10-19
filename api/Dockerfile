FROM golang:1.23.2-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache ca-certificates git curl

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/ylem_api

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM golang:1.23.2-alpine AS final

COPY --from=builder /usr/local/bin /usr/local/bin

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /opt /opt

#USER root

EXPOSE 7339

WORKDIR /opt/ylem_api

#CMD ["/opt/ylem_api/ylem_api", "server", "serve"]
