# Where is the Sun?
A educational implementation of a service visualizing sunnyness on a map (sunnyness = inverted cloud cover so far)

## Run
### Backend
- standalone: `cd backend && ./start.sh`
- docker:
    - `docker build -t sunnyness .`
    - `docker run --name c_sunnyness --env-file ./.env.local -p 8083:8083 sunnyness`
    - `docker start c_sunnyness`
- a `launch.json` for VS Code is included
- didn't get everything around modules so far, sometimes `export GO111MODULE="auto"` helps (running app/tests in combination with IDE etc...)

#### Test
- So far, tests are mostly only used to help with fiddling during implementation. However, the codebase provides mock support through `github.com/golang/mock/gomock` and `github.com/golang/mock/mockgen` (see [Generate Mocks](#generate-mocks) below) and the test suite can be easily enhanced into a functional test suite
- run tests from cli (from within `./backend/test`): `go test -v ./...`
- run single test from cli: `go test -v -timeout 30s -run ^TestReader$ backend/test`

#### Generate mocks
- install mockgen: `go install github.com/golang/mock/mockgen@v1.6.0`
- `mockgen` is then installed to your `GOPATH`, i.e. most often `$HOME/go`
- THEREFORE: make sure to add your `GOPATH` to your `PATH` so that `mockgen` is runnable the following command can run:
- Generate mocks (from within `./backend`): `go generate ./...`

#### Notes
- at some point I had to run `go get github.com/go-kit/kit/circuitbreaker@v0.12.0`. `go mod tidy|download` would not do the trick.

### Frontend