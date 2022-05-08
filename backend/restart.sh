# go generate ./...
go build -o bin/sunnyness main.go
export $(cat .env.local | xargs)
./bin/sunnyness
