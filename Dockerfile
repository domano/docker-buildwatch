FROM golang

COPY . /app
WORKDIR /app

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app ./
CMD ["./app"]  

