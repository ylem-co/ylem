FROM golang:1.27.4 AS builder

RUN apt-get update && apt-get install -y ca-certificates git curl

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/ylem_taskrunner

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM golang:1.27.4 AS final

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /opt /opt

USER nobody:nobody

EXPOSE 7335

WORKDIR /opt/ylem_taskrunner

#CMD ["/opt/ylem_taskrunner/ylem_taskrunner", "loadbalancer", "start"]
