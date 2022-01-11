# gRPC

### Proto
Create your ```.proto``` file and then go inside the package and enter in the terminal:
```bash
cd route_guide/
```
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/route_guide.proto
```
This will generate the interfaces for your proto code.