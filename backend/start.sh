# go generate ./...
# go mod install
go build -o bin/sunnyness main.go
export $(cat .env.local | xargs)
./bin/sunnyness
