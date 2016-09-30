YCE部署指南
===============


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

