# cache-service

## Setting Up Development Environment
Make sure you have setup the [global development environment](https://gitlab.ost.ch/ins/jalapeno-api/request-service/-/wikis/Development-Environment) first.

## Initialize Okteto
- Clone the repository:
```bash
$ git clone ssh://git@gitlab.ost.ch:45022/ins/jalapeno-api/cache-service.git
```
- Initialize okteto:
```bash
$ okteto init
```
- Replace content of okteto.yml with the following:
```yml
name: cache-service
autocreate: true
image: okteto/golang:1
command: bash
namespace: jagw-dev-<namespace-name>
securityContext:
  capabilities:
    add:
      - SYS_PTRACE
volumes:
  - /go/pkg/
  - /root/.cache/go-build/
  - /root/.vscode-server
  - /go/bin/
  - /bin/protoc/
sync:
  - .:/usr/src/app
forward:
  - 2349:2345
  - 8084:8080
environment:
  - ARANGO_DB=http://10.20.1.24:30852
  - ARANGO_DB_USER=root
  - ARANGO_DB_PASSWORD=jalapeno
  - ARANGO_DB_NAME=jalapeno
  - REDIS_PASSWORD=a-very-complex-password-here
  - SENTINEL_ADDRESS=sentinel.jagw-dev-michel.svc.cluster.local:5000
  - SENTINEL_MASTER=mymaster
  - KAFKA_ADDRESS=10.20.1.24:30092
  - LSNODE_KAFKA_TOPIC=gobmp.parsed.ls_node_events
  - LSLINK_KAFKA_TOPIC=gobmp.parsed.ls_link_events
```
