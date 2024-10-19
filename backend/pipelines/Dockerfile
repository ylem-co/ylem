FROM golang:1.23.2-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache ca-certificates git curl

# install Golang-migrate tool
# install Golang-migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz
RUN cp ./migrate /usr/local/bin

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/ylem_pipelines

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

EXPOSE 7336

WORKDIR /opt/ylem_pipelines

#CMD ["/opt/ylem_pipelines/ylem_pipelines", "server", "serve"]
