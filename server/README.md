# YLEM SERVER CONTAINER

<a href="https://github.com/ylem-co/ylem?tab=Apache-2.0-1-ov-file">![Static Badge](https://img.shields.io/badge/license-Apache%202.0-black)</a>
<a href="https://ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/website-ylem.co-black)</a>
<a href="https://docs.ylem.co" target="_blank">![Static Badge](https://img.shields.io/badge/documentation-docs.ylem.co-black)</a>
<a href="https://join.slack.com/t/ylem-co/shared_invite/zt-2nawzl6h0-qqJ0j7Vx_AEHfnB45xJg2Q" target="_blank">![Static Badge](https://img.shields.io/badge/community-join%20Slack-black)</a>

Nginx docker container placed in front of all the microservice APIs allowing to avoid CORS issue on the UI side.

# Configuration example:

```
server {
    listen 7331 default_server;
    listen [::]:7331 default_server ipv6only=on;

    add_header              X-Request-Id       $request_id;
    proxy_set_header        X-Request-Id       $request_id;

    server_name localhost;
    root /var/www/public;
    index index.html index.htm;
    
    location /integration-api/private/ {
        deny all;
    }

    location /user-api/private/ {
        deny all;
    }

    location /user-api/ {
        proxy_pass http://ylem_users:7333/;
    }

    location /stats-api/ {
         proxy_pass http://ylem_statistics:7332/;
    }

    location /pipeline-api/ {
        client_body_buffer_size     100M;
        client_max_body_size        100M;
        proxy_buffers 16 16k;  
        proxy_buffer_size 16k;
        proxy_pass http://ylem_pipelines:7336/;
    }

    location /integration-api/ {
         proxy_pass http://ylem_integrations:7337/;
    }

    location /oauth-api/ {
         proxy_pass http://ylem_api:7339/oauth-api/;
    }

    location ~ /\.ht {
        deny all;
    }

    location /.well-known/acme-challenge/ {
        root /var/www/letsencrypt/;
        log_not_found off;
    }
}
```

With this configuration, all microservices are now available under the [http://127.0.0.1:7331](http://127.0.0.1:7331)

For example:
* User service API: http://127.0.0.1:7331/user-api/
* Integration service API: http://127.0.0.1:7331/integration-api/
* Etc.
