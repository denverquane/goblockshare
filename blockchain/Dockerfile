FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN https_proxy=wwwproxy.sandia.gov:80 go get -d .
COPY main/blockchainServer.go    .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as runner
# RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
EXPOSE "5002"
ENV PORT 5002
CMD ["./app"]  