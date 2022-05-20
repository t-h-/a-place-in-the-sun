# A Place in the Sun
This is a (work in progress) pet project to learn Golang, its concurrency model and Gokit, a Golang microservice framework. The Vue frontend visualizes current sunnyness of a region in a Heatmap. For this, it queries the backend for a point grid of the current displyed rectangular region. The backend then queries a public weather api for all of the requested points in parallel, does some bilinear interpolation and returns the interpolated grid of points to the frontend.
Currently sunnyness = inverse of cloud cover (a value provided in the `weatherapi.com` response).

TODOs:
- parameter finetuning of heatmap and point grid granularity
- parameter externalization for frontend
- docker-compose
- big one: realworld Kubernetes deployment on GCP

## Backend
### Run
- standalone: `cd backend && ./start.sh`
- docker:
    - `docker build -t sunnyness/backend .`
    - `docker run --name c_sunnyness_backend --env-file ./.env.local -p 8083:8083 --rm sunnyness/backend`
- a `launch.json` for VS Code is included
- didn't get everything around modules so far, sometimes `export GO111MODULE="auto"` helps (running app/tests in combination with IDE etc...)

### Generate mocks
- install mockgen: `go install github.com/golang/mock/mockgen@v1.6.0`
- `mockgen` is then installed to your `GOPATH`, i.e. most often `$HOME/go`
- THEREFORE: make sure to add your `GOPATH` to your `PATH` so that `mockgen` is runnable the following command can run:
- Generate mocks (from within `./backend`): `go generate ./...`

### Test
- So far, tests are mostly only used to help with fiddling during implementation. However, the codebase provides mock support through `github.com/golang/mock/gomock` and `github.com/golang/mock/mockgen` (see [Generate Mocks](#generate-mocks) below) and the test suite can be easily enhanced into a functional test suite
- run all tests from cli (from within `./backend/test`): `go test -v ./...`
- run single test from cli: `go test -v -timeout 30s -run ^TestReader$ backend/test`

### Debug in VSCode
- launch "Launch go backend" configuration (found in `launch.json`)

### Notes
- at some point I had to run `go get github.com/go-kit/kit/circuitbreaker@v0.12.0`. `go mod tidy|download` would not do the trick.

## Frontend
### Run
- standalone
    - `npm install`
    - `npm run serve`
- docker
    - `docker build -t sunnyness/frontend .`
    - `docker run -it -p 8080:8080 --rm --name c_sunnyness_frontend sunnyness/frontend`
    - 

### Debug in VSCode
- `npm run serve`
- launch "Launch Chrome against localhost" configuration (found in `launch.json`)