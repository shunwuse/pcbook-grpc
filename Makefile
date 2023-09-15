run:
	go run main.go

gen:
	protoc --proto_path=proto proto/*.proto --go_out=./pb

clean:
	rm pb/*.go