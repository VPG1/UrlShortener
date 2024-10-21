# URL SHORTENER RESTFUL API SERVICE
## Description
Service for creating and managing URL aliases. With the ability to create an account and temporary links.
## Startup instruction
Clone the repo
``` console
git clone https://github.com/VPG1/UrlShortener.git
```
Download all dependencies
``` console
go mod download
```
Create if not exist config file in config directory. </br>
For example:
``` yaml 
# config/config.yaml
env: "local" # (local, prod)
alias_len: 8
postgres_server:
  address: "localhost"
  port: "5432"
  user: "user"
  password: "passwd"
  db_name: "url-storage"
http_server:
  address: "localhost"
  port: "8080"
  timeout: 4s
  idle_timeout: 60s
```
Start PostgreSQL container
``` console
docker compose up -d
```
Run the service
``` console
go run ./cmd/url-shortener/main.go
```

## Usage:

You can use curl, postman or swagger documentation to create a user and get a JWT token. Next, you can create a shortened link, get your shortened links and delete them.
