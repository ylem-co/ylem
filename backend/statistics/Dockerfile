FROM golang:1.23.2-alpine AS builder

RUN apk add --no-cache ca-certificates git wget

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/ylem_statistics

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM golang:1.23.2-alpine AS final

RUN apk add --no-cache wget

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /opt /opt

#USER root

EXPOSE 7332

WORKDIR /opt/ylem_statistics

#CMD ["/opt/ylem_statistics/ylem_statistics", "server", "serve"]
