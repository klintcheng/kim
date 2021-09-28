# KIM

King IM CLoud

## 简介

**kim 是一个高性能分布式即时通信系统。**

![structure.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/2633b07fd1a144d685ceed9be5f64911~tplv-k3u1fbpfcp-watermark.image)

同时，也是`《分布式IM原理与实战: 从0到1打造即时通讯云》`一书的实战项目。你可以点击下图，观看详细介绍。

[![book](./book.png)](https://juejin.cn/book/6963277002044342311)

- Web SDK： [Typescript SDK](https://github.com/klintcheng/kim_web_sdk)
- Flutter SDK： [Flutter SDK](https://github.com/szhua/KimSdk)
  - 由[@szhua](https://github.com/szhua)小友提供

## 环境准备

1. 安装docker
2. 安装docker-compose
3. 启动环境
   -  docker-compose -f "docker-compose.yml" up -d --build
4. 进入Mysql，修改访问权限：
   1. docker exec -it kim_mysql /bin/sh
   2. mysql -uroot -p123456
   3. GRANT ALL ON *.* TO 'root'@'%';
   4. flush privileges;
5. 创建数据库
   1. create database kim_base default character set utf8mb4 collate utf8mb4_unicode_ci;
   2. create database kim_message default character set utf8mb4 collate utf8mb4_unicode_ci;

6. docker部署服务kim单节点服务
   - docker-compose -f "docker-compose-kim.yml" up -d --build