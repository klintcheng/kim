# KIM

King IM CLoud

## 简介

    kim 是一个高性能分式式通信架构。 


## 环境准备

1. 安装docker
2. 安装docker-compose
3. 启动环境
   - > docker-compose -f "docker-compose.yml" up -d --build
4. 进入Mysql，修改访问权限：
   1. > docker exec -it kim_mysql /bin/sh
   2. > mysql -uroot -p123456
   3. > GRANT ALL ON *.* TO 'root'@'%';
   4. > flush privileges;
5. 创建数据库
   1. > create database kim_base default character set utf8mb4 collate utf8mb4_unicode_ci;
   2. > create database kim_message default character set utf8mb4 collate utf8mb4_unicode_ci;

## 通信层代码演示

首先进入examples目录：

1. **启用服务端**

```cmd
go run main.go mock_srv -p ws
INFO[0000] started                                       id=srv1 listen=":8000" module=ws.server
```

2. **启用客户端**

```cmd
$ go run main.go mock_cli -p ws
WARN[0000] 1uWbA9ajf86A44J8t4k2AtsadQG receive message [hello from server ] 
WARN[0001] 1uWbA9ajf86A44J8t4k2AtsadQG receive message [hello from server ] 
WARN[0002] 1uWbA9ajf86A44J8t4k2AtsadQG receive message [hello from server ] 
WARN[0003] 1uWbA9ajf86A44J8t4k2AtsadQG receive message [hello from server ] 
...
```