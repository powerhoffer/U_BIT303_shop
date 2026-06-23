# YUTANK 商城管理系统使用说明

[English version](README.md)

## 环境要求

本项目是 GoFrame 单体后端 + Vue 管理端前端。以下环境已用于本地开发验证：

- 操作系统：WSL Ubuntu 20.04
- Go：1.25.5
- GoFrame CLI：2.10.0
- GoFrame 框架：2.9.5
- MySQL：8.0.x
- Node.js：22.22.2
- npm：10.9.7

默认本地数据库配置：

    Host: 127.0.0.1
    Port: 3306
    User: root
    Password: 123456
    Database: yutank_shop

说明：yutank_shop 是默认本地开发数据库名。

## 启动项目

### 1. 初始化数据库

进入项目目录：

    cd /path/to/yutank-shop

验证 MySQL 连接：

    mysql -uroot -p123456
    SELECT 1;

导入初始化 SQL：

    mysql -uroot -p123456 < hack/init.sql

### 2. 启动后端

先启动后端服务。默认接口地址：

    http://127.0.0.1:8000

启动命令：

    cd /path/to/yutank-shop
    /usr/local/go/bin/go run main.go

健康检查：

    http://127.0.0.1:8000/health

Swagger：

    http://127.0.0.1:8000/swagger/

### 3. 启动管理端前端

再启动管理端前端。默认访问地址：

    http://localhost:8080

启动命令：

    cd /path/to/yutank-shop/frontend_manage
    npm install
    npm run dev

### 4. 登录管理端

在浏览器中打开管理端登录页：

    http://localhost:8080/#/login

测试账号：

    账号：root
    密码：123456

登录成功后会进入管理端工作台。

## 当前功能

- 工作台：查看当前员工信息、我的积分和管理模块概览。
- 员工列表：新增员工、编辑员工信息、启用或禁用员工、重置员工密码。
- 我的积分：查看当前登录员工的积分余额和积分流水。
- 积分操作：选择员工后进行加积分、扣积分，并查看员工积分流水。
- 分类列表：查看当前商品分类信息。
- 商品管理：新增商品、编辑商品、查看商品详情、上架或下架商品。
- 员工端接口：已提供商品浏览和购物车后端接口，后续可继续开发员工端页面。
- 右上角用户菜单：修改当前登录账号密码，或退出登录。

## 注意事项

- 管理端前端接口统一请求后端 /backend 路由。
- 员工端接口统一使用 /frontend 路由。
- 如果无法登录，请先确认目标环境中已经存在 root 员工账号，并且密码为 123456。
