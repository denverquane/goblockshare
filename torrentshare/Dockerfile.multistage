FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -d ./main
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main

FROM alpine:latest as runner
# RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app .
EXPOSE "7000"
ENV PORT 7000
CMD ["./app"]  