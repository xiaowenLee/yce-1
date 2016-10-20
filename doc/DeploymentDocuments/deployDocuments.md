YCE部署指南
===============

本文档主要包含了首次部署YCE的操作流程、YCE的更新操作流程以及相关说明。

### 首次部署 
#### 部署前准备
在首次部署开始之前需要做一些准备。这些准备包括

* 基础环境准备
* 获取部署脚本
* 私有镜像仓库搭建
* 相关镜像创建
* Git连接检查
* NodePort检查
* 数据库初始化


#####  基础环境准备
基础环境搭建, 需要若干台机器搭建Kubernetes集群, 在master节点上安装好git、go、Docker等开发环境, 建议版本为:
* git的版本为1.9.1
* go的版本为go 1.6.2 linux/amd64
* Docker的版本为docker 1.11.1
* Kubernetes的版本为1.2.0

##### 获取部署脚本 
从prod3(10.149)上拷贝包~/ycetestplace/deploy-prod.tar.gz到Kubernetes的master节点上, 并解压, 得到目录deploy/

##### 私有镜像仓库搭建
私有仓库的搭建, 用Docker直接启动, 对外开放端口15000, 域名为img.reg.3g, 相应的证书为deploy/docker/bin/domain.crt. 分别给集群里的每个节点更新证书,并重启Docker。
证书md5为: d3cf84f3dc9efb9175e09fe3625cda2c

##### 相关镜像创建
相关镜像创建(或导入), 这里的镜像包含了YCE运行所要依赖的MySQL镜像、Redis镜像等, 它们的名称为:
* MySQL镜像, 地址为img.reg.3g:15000/mysql:5.7.13
* Redis镜像, 主节点地址为img.reg.3g:15000/redis:3.0.7
* 从节点地址为img.reg.3g:15000/redis-slave:3.0.3

如果在别的镜像仓库已有这些镜像, 可以直接下载或者导出为tar文件,拷贝过来,再docker push到刚才搭建的镜像仓库。

##### Git连接检查
将master节点的公钥放入管理员Github账户的SSH-Keys里, 以便从开源仓库下拉源代码文件。

##### NodePort检查
在redis/redis-master-svc.yaml, redis/redis-slave-svc.yaml里指定了redis-master, redis-slave开放的NodePort端口为32379和32380。
在yce/yce-svc.yaml里指定了yce开放的NodePort端口为32080
要确保上面的端口未被占用,或在相应的文件里更改相应的端口。

##### 数据库初始化
数据库初始表位于mysql/sql/yce-initdata.sql, 但是在导入前需要对里面的内容进行初始化,以适配安装的集群环境。
需要改动的表有:
datacenter表: 这个表所指示的为数据中心列表, 需要将已有的Kubernetes集群名称、IP地址、端口等填入。
organization表: 这个表所指示的是组织信息, 其中包含了该组织对应的数据中心的列表。


#### 部署步骤
切换到deploy/

运行deploy.sh脚本

部署完后不建议删除deploy目录,以便后续更新使用

### 更新 
#### 更新准备
更新需要准备好可运行的YCE环境, 通过打包新的镜像, 从YCE平台的滚动升级功能进行更新
#### 更新步骤
切换到deploy/目录

运行update.sh脚本
git pull
go build
docker push




### 以下为旧版本, 无效

部署前准备
---------------
YCE目前是Yao的alpha版本, 此版本的部署需要:

安装git, go, Docker, Kubernetes等开发环境。其中:
* git的版本为1.9.1
* go的版本为go 1.6.2 linux/amd64
* Docker的版本为docker 1.11.1
* Kubernetes的版本为1.2.0

镜像包括:

- MySQL镜像, 地址为img.reg.3g:15000/mysql:5.7.13
- Redis镜像, 主节点地址为img.reg.3g:15000/redis:3.0.7, 从节点地址为img.reg.3g:15000/redis-slave:3.0.3
- YCE镜像, 地址为img.reg.3g:15000/yce-alpha:test-$(shell date +%F), 首次部署的yce镜像版本为yce-alpha:test-2016-09-26, 写在yce/yce-deployment.yaml里

YCE的正常运作需要这些程序共同运行, 启动脚本及yaml文件均放置在一个名为deploy的目录下, 目录结构及说明:

    |-- deploy
	    |-- Makefile                                   
	    |-- docker
	        |-- Dockerfile              
	        |-- bin
	            |-- domain.crt                      // 访问私有镜像仓库img.reg.3g:15000的证书
	            |-- yce
	        |-- src
	            |-- app
	                |-- frontend                    // 前端代码 
	                |-- backend                     // 后端代码 
                |-- ....                            // 第三方代码
		|-- yce
		    |-- yce-deploymnet.yaml                 
		    |-- yce-svc.yaml
		    |-- yce-up.sh
		    |-- yce-down.sh
		|-- mysql
    	    |-- mysql-deployment.yaml
		    |-- mysql-svc.yaml
    	    |-- mysql-up.sh
		    |-- mysql-down.sh
    	    |-- sql
		        |-- yce-initdata.sql                // 带有初始数据的MySQL表 
		        |-- yce.sql                         // 20160922的数据导出的表
		|-- redis
		    |-- redis-master-deployment.yaml
		    |-- redis-master-svc.yaml
		    |-- redis-slave-deployment.yaml
		    |-- redis-slave-svc.yaml
		    |-- redis-up.sh
		    |-- redis-down.sh
		|-- namespace
		    |-- limits.yaml                         // 为某命名空间下每个容器配置资源限额 
		    |-- namespace.yaml                      // 创建命名空间
		    |-- quota.yaml                          // 为某命名空分配资源配额
		|-- delete.sh                               // 顺序删除每个节点上的已有镜像
			

其中Makefile用于简化部署和删除, `make init`: 初始化, `make build`: 一键构建, `make deploy`: 一键部署,  `make clean`: 一键删除.

另外,在yce、mysql、redis下除了包含相应的部署和服务yaml文件外,还包含了相应的手动启动/停止脚本: *-up.sh / *-down.sh

部署流程
--------------

首次部署:

1. `make init`: 建立frontend和backend代码目录,初始化git并添加远程仓库地址。
1. `make build`: 从远程仓库获取最新代码; 编译为二进制可执行文件; 打成Docker镜像; 推到镜像仓库;
2. `make deploy`: 部署MySQL镜像; 导入初始数据库; 部署redis集群; 部署yce; 

更新流程
--------------
更新:

1. `make build`: 更新前端代码,并重新编译和打Docker镜像 
2. 通过YCE管理系统进行滚动升级(建议)

