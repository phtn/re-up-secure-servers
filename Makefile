run:
	go run cmd/fast/main.go

b:
	go build -o ./build/fast cmd/fast/main.go

tidy:
	go mod tidy

clean:
	rm -rf ./build/fast
