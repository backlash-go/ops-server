FROM golang:1.20 AS build-env

ENV WORKSPACE=/workspace

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

RUN mkdir $WORKSPACE

WORKDIR $WORKSPACE
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./cmd

RUN echo 'new --'

FROM alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache libc6-compat tzdata curl gcompat  \
&& ln -s /lib/libc.so.6 /usr/lib/libresolv.so.2 \
&& echo "Asia/Shanghai" > /etc/timezone \
&& ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /app

COPY --from=build-env /workspace/cmd /app/

WORKDIR /app

ENTRYPOINT ./cmd
