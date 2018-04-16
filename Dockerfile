FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -d .
COPY main.go    .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as runner
# RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
EXPOSE "8080"
ENV PORT 8080
CMD ["./app"]  