FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 1
RUN apk add --no-cache gcc musl-dev
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

#ADD go.mod .
#ADD go.sum .
#RUN go mod download
COPY . .
COPY soot-analysis-1.0-SNAPSHOT.jar /app/
RUN --mount=type=cache,target=/root/.cache/go-build,id=go_build_cache,sharing=shared go build -ldflags='-s -w -extldflags "-static"' -o /app/microsegement  main.go


FROM alpine:3.18.0
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/soot-analysis-1.0-SNAPSHOT.jar /app/soot-analysis-1.0-SNAPSHOT.jar
COPY --from=builder /app/microsegement /app/microsegement

CMD ["/bin/sh"]