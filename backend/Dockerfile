FROM golang:1.17-alpine

WORKDIR /app

COPY . ./

RUN go mod download
RUN go build -o /sunnyness

EXPOSE 8083

CMD [ "/sunnyness" ]