# Tensor Programming Tutorials
Studies in gRPC following [TensorProgramming](https://www.youtube.com/watch?v=1MPWPq2N768&list=PLJbE2Yu2zumCe9cO3SIyragJ8pLmVv0z9&index=28&ab_channel=TensorProgramming)

## [BasicAPI](https://youtu.be/Y92WWaZJl24?list=PLJbE2Yu2zumCe9cO3SIyragJ8pLmVv0z9)
Simple API using *gRPC* server and *Gin*.

### Usage
In one terminal, type in:
```bash
cd BasicAPI/server/
go run ./main.go
```

In another, type in:
```bash
cd BasicAPI/client/
go run ./main.go
```

By default the server runs on <code>:8080</code> and client runs on <code>:4040</code>.

## [ChatApp](https://youtu.be/mML6GiOAM1w?list=PLJbE2Yu2zumCe9cO3SIyragJ8pLmVv0z9)
Simple Chat Application using *gRPC* and *Docker*.

### Usage
#### Server without Docker
In one terminal, type the following to run the server:
```bash
cd ChatApp/
go run ./main.go
```
#### Server with Docker
```bash
cd ChatApp/
sudo docker build -t chat_grpc_server .
sudo docker run --rm --name chat_grpc_server -it -p 8080:8080 chat_grpc_server
```

#### Client
In another terminal, type in:
```bash
cd ChatApp/client/
go run ./main.go
```

By default the server runs on <code>:8080</code> and client runs on <code>:4040</code>.