# Where is the Sun?


## Run
### Backend
`cd backend && ./restart.sh

#### Generate mocks
- install mockgen: `go install github.com/golang/mock/mockgen@v1.6.0`
- `mockgen` is then installed to your `GOPATH`, i.e. most often `$HOME/go`
- THEREFORE: make sure to add your `GOPATH` to your `PATH` so that `mockgen` is runnable the following command can run:
- Generate mocks (from within `./backend`): `go generate ./...`