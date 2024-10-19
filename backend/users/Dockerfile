FROM golang:1.23.2-alpine AS builder

ENV CGO_ENABLED=0

RUN apk add --no-cache ca-certificates git curl openssh-keygen openssl

# install Golang-migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.13.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 migrate
RUN cp ./migrate /usr/local/bin

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

WORKDIR /opt/ylem_users

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build .

VOLUME /opt/ylem_users/config/jwt

# Generate SSH keys for JWT if don't exist
RUN printf "n" | ssh-keygen -t rsa -b 4096 -m PEM -f ./config/jwt/private.pem -q -N ""; exit 0
RUN printf "" | openssl rsa -in ./config/jwt/private.pem -pubout -passin 'pass:' -outform PEM -out ./config/jwt/public.pem; exit 0

FROM golang:1.23.2-alpine AS final

#RUN apk add --no-cache openssh-keygen openssl

COPY --from=builder /opt /opt

#USER root

EXPOSE 7333

#VOLUME /opt/ylem_users/config/jwt

WORKDIR /opt/ylem_users

COPY --from=builder /user/group /user/passwd /etc/

COPY --from=builder /usr/local/bin /usr/local/bin

#USER nobody:nobody

#CMD ["/opt/ylem_users/ylem_users", "server", "serve"]
