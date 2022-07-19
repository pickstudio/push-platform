# push-platform

``` markdown
other_service => [PushPlatform]http_server => SQS => [PushPlatform]worker => SNS for Push
```

# getting started

### Settings up
#### 1. create `.env`
```
# required environment
ENV=development
DEBUG=true

LOCALHOST_HTTP_DSN="0.0.0.0:50100"
LOCALHOST_HTTP_TIMEOUT="2s"

# AWS Secrets keys only use local machine, not remote machines
AWS_REGION=ap-northeast-2
AWS_ACCESS_KEY_ID=###
AWS_SECRET_ACCESS_KEY=###
```

#### 2. direnv allow
- [direnv](https://www.44bits.io/ko/post/direnv_for_managing_directory_environment) is good for backend engineers
``` bash
~/Workspace/pickstudio/push-platform   main ✚ ● ?  direnv allow
# direnv: export +AWS_REGION +DEBUG +ENV +HTTP_SERVER_DSN ....
```

#### 3. go run server
``` bash
go run cmd/http_server/main.go
go run cmd/worker/main.go
```

### Developments

#### run unit test
``` bash
go test ./... -v
```

#### run integration test

``` bash
go run cmd/http_server/main.go
# on another terminal
go run cmd/e2e_test/main.go
```

#### generate codes
- openapi code following `api/oapi/oapi.yaml`

``` bash
oapi-codegen --config=api/oapi/v1/oapi-codegen-config.yaml api/oapi/v1/v1.yaml > api/oapi/v1/v1.oapi.go
go generate ./...
```