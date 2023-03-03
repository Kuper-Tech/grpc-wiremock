server {
    server_name {{ .Domain }} {{ .Domain }}.*;

    listen      80;
    listen [::]:80;
    listen 443 ssl;

    ssl_certificate /etc/ssl/mock/mock.crt;
    ssl_certificate_key /etc/ssl/mock/mock.key;

    access_log  /var/log/nginx/{{ .Domain }}.access.log  main;

    location / {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_pass http://localhost:{{ .Port }};
    }
}
