FROM golang:buster as builder
COPY . /source
WORKDIR /source
# speedup in china
#RUN go env -w GOSUMDB="sum.golang.google.cn" GOPROXY="https://goproxy.cn,direct"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main -ldflags "-X 'wechat-mp-server/config.Version=$(git show -s --format=%h)'"

FROM alpine:latest
WORKDIR /root
COPY --from=builder /source/main ./
EXPOSE 10151
ENTRYPOINT ["./main"]