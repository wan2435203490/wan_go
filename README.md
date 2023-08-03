# wan_go

集成go项目

# 打包
写在Makefile的

# 部署 blog
```shell
docker stop blog
docker rm blog
docker rmi blog
cd /usr/wan_go
docker build -t blog -f deploy/blog/deploy.Dockerfile .
docker-compose -f deploy/blog/docker-compose.yaml up -d
ng restart

#重启容器
docker restart blog

#更改可执行文件
docker cp /usr/wan_go/bin/blog blog:/wan_go/bin
```


# 部署 landlord
```shell
docker stop landlord
docker rm landlord
docker rmi landlord
cd /usr/wan_go
docker build -t landlord -f deploy/landlord/deploy.Dockerfile .
docker-compose -f deploy/landlord/docker-compose.yaml up -d
ng restart
```