# Go PS-3-T
Cara run?
1. `go run main.go`
2. kalau mau tes--misalnya--. Buka `localhost:port/swagger/index.html`

# Installation
```
# clone the repository
git clone [this_git_url]

# set up environment
cp .env.example .env

# go run
go run main.go
```

# Configuration
Create a `.env` file in the root directory with the following variables:
```
DB_HOST=localhost
DB_USER=user
DB_PASSWORD=pw
DB_NAME=ps3t
DB_PORT=5432 #Example: postgresql
PORT=3000
MODE=DEBUG
PROD_HOST=#Your production host
DEBUG_HOST=0.0.0.0
GRPC_HOST=localhost
GRPC_PORT=50051
EXTERNAL_GRPC_HOST=localhost # external host
EXTERNAL_GRPC_PORT=50052
```

# gRPC
### Define & Re-compile the gRPC Code (.proto)
Anytime you add or update the protobuf definition in `.proto` file, 
you need to run this command to re-compile the source code from 
protobuf.

```
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
user/user.proto
```
On this command, you are regenerating source code from 
the `./user/user.proto` file.

After the protobuf has been compiled, then you can implement the 
services in Go.

For example, if you just updated the `./user/user.proto` and 
recompile it, you can create the gRPC handler to satisfy the 
protobuf updates in `./handler/user_grpc_handler.go`. Then you 
can create a new service or even re-use your existing service 
(eg. `RegisterUser()` in `./service/user/user_service.go`) 
to be put in the handler.

### Receiving gRPC Request from External
Since the gRPC server has been setup along with the HTTP API server, 
your gRPC server will also start once you run `go run main.go` or run 
the executable binary.

If you only need to serve gRPC, then you can surely remove the HTTP 
API setup and the gRPC should remains running.

To see how your client can consume your gRPC API, you can see a 
sample of gRPC client setup under the `./grpc/client/main.go`.

You can also run that gRPC client to test calling your gRPC server by 
running `go run ./grpc/client/main.go`.

### Sending gRPC Request to External Server
We also have already set a sample for our service to call an external 
gRPC API (basically, we act as a client to send request to another 
server). We basically replicate what we've done with 
`./grpc/client/main.go`, but it is now integrated with our code 
structure.

** Note that this example require the HTTP API server to run.

To test the request, 
firstly ensure the `EXTERNAL_GRPC_HOST` & `EXTERNAL_GRPC_PORT` value 
has been set correctly to the target server address, then run:
```
go run main.go
```

Since we're simulating a HTTP API that depends to external service 
via gRPC, we need to trigger the API endpoint. 
You may call this `curl` below or try in your API plaform like Postman.
```
curl -X 'POST' \
  'http://${EXTERNAL_GRPC_HOST}:${EXTERNAL_GRPC_PORT}/api/users/register/grpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": "sample_id",
  "username": "sample_username"
  "email": "sample@email.com",
  "password": "sample_plain_password",
}'
```

# Running the App

In Go, there are two ways to run the app

## Build

```go
# For build, run this command
go build -o .build/<name-of-build.extension>

# NOTE: it is important to put the build inside of the .build folder
# to ensure the gitignore caught up with the files

# After build go application
.build/<name-of-build.extension>
```

## Go Run
```
# just do this
go run main.go

# then your operating system asking for firewall permission
```

## On Docker

```
docker-compose up -d
# Dengan flag -d untuk menjalankan container di background (detached mode).
## Dengan menggunakan -d, terminal akan langsung kembali ke prompt tanpa menampilkan log container di terminal.
## Container akan terus berjalan di background setelah perintah ini dijalankan.

# Kalau mau tambah flag --build jika ada perubahan pada Dockerfile
docker-compose up -d --build
## Flag --build digunakan untuk memaksa Docker Compose membangun ulang image sebelum menjalankan container.
## Biasanya digunakan ketika ada perubahan pada Dockerfile atau file yang terkait dengan image.
## Perintah ini akan melakukan build ulang image dan kemudian menjalankan container di background.

```

# API Documentation
One the application is running, you can access the Swagger API documentation at:
```
http://localhost:3000/swagger/index.html
```
