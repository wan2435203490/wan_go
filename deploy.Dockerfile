FROM golang:1.19 as build
MAINTAINER "wan 2435203490@qq.com"

# go mod Installation source, container environment variable addition will override the default variable value
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

#RUN apt-get update && apt-get install apt-transport-https && apt-get install procps\
#&&apt-get install net-tools
##Non-interactive operation
#ENV DEBIAN_FRONTEND=noninteractive
#RUN apt-get install -y vim curl tzdata gawk
##Time zone adjusted to East eighth District
#RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

# 设置固定的项目路径
ENV WORKDIR /wan_go
ENV CONFIG_NAME $WORKDIR

WORKDIR /wan_go
COPY . .
ADD /config/config.yaml /wan_go/config/

EXPOSE 8090
RUN chmod +x *.sh
CMD ["./blog"]