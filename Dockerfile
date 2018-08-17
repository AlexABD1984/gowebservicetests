FROM golang:latest 
ENV NATS_URI=nats://nats:4222
WORKDIR /go/src/app
RUN go get -v -u github.com/xeipuuv/gojsonschema
RUN go get -v -u github.com/nats-io/go-nats
RUN go get -v -u github.com/buaazp/fasthttprouter
RUN go get -v -u github.com/valyala/fasthttp
WORKDIR /go/src/app
COPY . .
RUN go build -o main .
CMD ["/go/src/app/main"]
