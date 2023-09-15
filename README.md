```shell
# install protoc via brew on mac
brew install protobuf

# install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# install protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# generate pb file from proto via protoc
make gen

# clean up generated pb file
make clean

# run program
make run
```