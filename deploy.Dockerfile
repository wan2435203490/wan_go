FROM golang:1.19 as build
MAINTAINER "wan 2435203490@qq.com"

# go mod Installation source, container environment variable addition will override the default variable value
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

#定义时区参数
ENV TZ=Asia/Shanghai
#安装时区数据包
RUN apk add tzdata
#设置时区
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo '$TZ' > /etc/timezone

WORKDIR /wan_go
COPY . .
ADD /config/config.yaml /wan_goconfig/

EXPOSE 8090
CMD ["./blog"]