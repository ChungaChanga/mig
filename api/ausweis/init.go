package ausweis

//go:generate find . -name *.proto -exec protoc -I pb {} --go_out=pb --go-grpc_out=pb ;
