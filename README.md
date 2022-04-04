# Zero-UI No Auth Proxy
Simple reverse proxy which injects the token and thus disables auth.
see https://github.com/dec0dOS/zero-ui

Available at https://hub.docker.com/r/maltegrosse/zero-ui-noauth


### Environment variables
```
PROTOCOL (http|https)
AUTH_PATH (/auth/login)
EXPOSE_PORT (9999)
CONNECT_HOST=192.168.2.1
CONNECT_PORT=4000
USER (admin)
PASSWORD (zero-ui)
```

### Run Docker
Add environment variables..
```
docker run -ti -e ... maltegrosse/zero-ui-noauth:1.0.0
```

### Run Docker-Compose
```
version: "3.7"

services:
  zt-ui-noauth:
    image: maltegrosse/zero-ui-noauth:1.0.0
    container_name: zu-ui-noauth
    restart: unless-stopped
    expose:
      - "9999"
    ports:
      - "9999:9999"
    depends_on:
      - zero-ui
    environment:
      - PROTOCOL=http
      - AUTH_PATH=/auth/login
      - EXPOSE_PORT=9999
      - CONNECT_HOST=zu-main
      - CONNECT_PORT=4000
      - USER=admin
      - PASSWORD=zero-ui
    
    ...

```

### Build Docker
```
docker build -t maltegrosse/zero-ui-noauth:latest .
```

### Additional Notes
Please dont abuse this proxy server for malicious js injection
