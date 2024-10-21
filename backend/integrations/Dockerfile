FROM public.ecr.aws/b5q6i6w4/ylem-public-images@sha256:c73b7d09874740f2c0df7003954fbbc46310ae30363e4a201d809aac5dff6afc AS builder

RUN apt-get update && apt-get install -y ca-certificates git curl

# install Golang-migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz
RUN cp ./migrate /usr/local/bin

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/ylem_integrations

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build .

FROM public.ecr.aws/b5q6i6w4/ylem-public-images@sha256:c73b7d09874740f2c0df7003954fbbc46310ae30363e4a201d809aac5dff6afc AS final

COPY --from=builder /usr/local/bin /usr/local/bin

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /opt /opt

#USER nobody:nobody

EXPOSE 7337

VOLUME /opt/ylem_integrations/config/keys

WORKDIR /opt/ylem_integrations

#CMD ["/opt/ylem_integrations/ylem_integrations", "server", "serve"]
