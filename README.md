# graph-db-feeder

## gRPC
- When the file `proto/graph-db-feeder.proto` is updated, this command needs to be run to recompile the code:
```bash
$ /bin/protoc/bin/protoc --proto_path=./proto --go_out=./proto --go_opt=paths=source_relative --go-grpc_out=./proto --go-grpc_opt=paths=source_relative ./proto/graph-db-feeder.proto
```

## Settin Up Development Environment

### Step 1: Initialize Okteto
- Clone the repository:
```bash
$ git clone ssh://git@gitlab.ost.ch:45022/ins/jalapeno-api/sr-app.git
```
- Initialize okteto:
```bash
$ okteto init
```
- Replace content of okteto.yml with the following:
```yml
name: graph-db-feeder
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
  - 2348:2345
  - 8083:8080
environment:
  - ARANGO_DB=http://bongard.internet-box.ch:8529
  - ARANGO_DB_USER=root
  - ARANGO_DB_PASSWORD=jalapeno
  - ARANGO_DB_NAME=jalapeno
```

### Step 2: Initialize the Container
- Open VSCode in the root of the repository.
- Hit `cmd`  + `p` to open the command pallet.
- Enter `>` and then choose `okteto up`
- When prompted, choose your `okteto.yml` file.
- When prompted, choose `Linux` as the containers operating system.

### Step 3: Setup the Container
- In the VSCode instance from the container, install the `Go` extension, otherwise the command `go` will not work on the VSCode command line.
- Install any additional extensions you want.

#### Install the Protocol Buffer Compiler
Here is the official guide: https://grpc.io/docs/protoc-installation/  
Just run these commands:
```bash
$ apt update
$ apt install unzip
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip
$ unzip protoc-3.17.3-linux-x86_64.zip -d /bin/protoc
$ rm protoc-3.17.3-linux-x86_64.zip
```

#### Install the gRPC Library for Go
Here is the official guide: https://grpc.io/docs/languages/go/quickstart/  
Just run these commands:
```bash
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

