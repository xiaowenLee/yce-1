<img src="http://kubernetes.io/kubernetes/img/warning.png" alt="WARNING" width="25" height="25"> 

####修改请谨慎

YCE部署指南
==============

作者: [maxwell92](https://github.com/maxwell92)

最后修订: 2016-11-06

目录
--------------
###目的
文档主要包含了首次部署YCE的操作流程、YCE的更新操作流程以及相关说明。

如果确定部署环境已搭建好, 请直接按[首次部署步骤](https://github.com/maxwell92/yce/blob/master/doc/DeploymentDocuments/deployDocuments.md#首次部署步骤)进行部署。如果不确定部署环境, 请先阅读[首次部署须知](https://github.com/maxwell92/yce/blob/master/doc/DeploymentDocuments/deployDocuments.md#首次部署须知)
如果非首次部署,即进行更新, 请直接按[更新步骤](https://github.com/maxwell92/yce/blob/master/doc/DeploymentDocuments/deployDocuments.md#更新步骤)进行更新

###首次部署
#### 首次部署须知 
在首次部署开始之前需要做一些准备。这些准备包括

* 切换到root用户
* 基础环境准备
* 获取部署脚本
* 私有镜像仓库搭建
* 相关镜像创建
* Git连接检查
* NodePort检查
* 数据库初始化

#####  切换到root用户
需要切换到root用户进行后续操作, *注: 为减少安全隐患, 请在部署完成后退出root*

#####  基础环境准备
基础环境搭建, 需要若干台机器搭建Kubernetes集群, 在master节点上安装好git、go、Docker等开发环境, 建议版本为:
* git的版本为1.9.1
* go的版本为go 1.6.2 linux/amd64
* Docker的版本为docker 1.11.1
* Kubernetes的版本为1.2.0
* 提供导入初始数据功能的MySQL

注: 用于导入初始数据的MySQL

##### 获取部署脚本 
从prod3(10.149)上拷贝包~/archive/yce/deploy.tar.gz到Kubernetes的master节点上, 并解压得到目录deploy/

##### 私有镜像仓库搭建
私有仓库的搭建, 用Docker直接启动, 对外开放端口15000, 域名为img.reg.3g, 相应的证书为deploy/docker/bin/domain.crt. 分别给集群里的每个节点更新证书,并重启Docker。
证书md5为: d3cf84f3dc9efb9175e09fe3625cda2c
更新方法: 为Kubernetes的每个节点的/etc/docker/certs.d目录下,添加img.reg.3g:15000目录, 将该证书放置于此目录下, 最后重启docker

##### 相关镜像创建
相关镜像创建(或导入), 这里的镜像包含了YCE运行所要依赖的MySQL镜像、Redis镜像等, 它们的名称为:
* MySQL镜像, 地址为img.reg.3g:15000/mysql:5.7.13, 校验: 1195b21c3a45 
* Redis镜像, 主节点地址为img.reg.3g:15000/redis:3.0.7, 校验: bab6d3ebc78c
* 从节点地址为img.reg.3g:15000/redis-slave:3.0.3, 校验: d1ae45cdf710
* Ubuntu-base:v3, img.reg.3g:15000/ubuntu-base:v3, 校验: 9bce8c1d0877

如果在别的镜像仓库已有这些镜像, 可以直接下载或者导出为tar文件,拷贝过来,再docker push到刚才搭建的镜像仓库。
这些镜像同时存放在prod3(10.149):~/archive/yce/baseimg.tar.gz里

##### Git连接检查
将master节点的公钥放入YCE管理员Github账户的SSH-Keys里, 以便从Github下拉源代码文件。
放入公钥后需要在服务器上执行: `ssh -T git@github.com`, 看到 `Hi maxwell92! You've successfully authenticated, but GitHub does not provide shell access.` 类似文字即表示认证通过。

##### NodePort检查
在deploy/yce/yce-svc.yaml里指定了yce开放的NodePort端口为32080. 
要确保上面的端口未被占用,否则需要释放该端口。
*注: 这个端口已经讨论确定*

##### 数据库初始化
数据库初始表位于deploy/mysql/sql/yce-initdata.sql, 里面的内容默认初始化适配待部署的集群环境。

如果对表里的初始内容不确定, 需要检查的表有:
datacenter表: 这个表所指示的为数据中心列表, 检查待部署的Kubernetes集群名称、IP地址、端口等。 *注: 名称初始化为湖南*
organization表: 这个表所指示的是组织信息, 其中包含了该组织对应的数据中心的列表, 初始化为yce
user表: 这个表所指示的是用户表, 初始化为仅有admin 
nodeport表: 这个表所指示的是nodePort使用, 初始化为32080, 供yce使用

管理员首次登录后要确认数据中心表信息是否正确, 这样nodePort才会被创建到nodePort表里, 提供更改名字功能。


#### 首次部署步骤 
切换到deploy/

运行deploy.sh脚本: ./deploy.sh

部署完后不建议删除deploy目录,以便后续更新使用

*注: 部署过程中如果出错, 请查看脚本信息或联系[maxwell](github.com/maxwell92))*

### 更新 
#### 更新准备
更新需要准备好可运行的YCE环境, 通过打包新的镜像, 从YCE平台的滚动升级功能进行更新
#### 更新步骤
切换到deploy/目录

运行update.sh脚本: ./update.sh

使用admin用户登录yce, 在应用管理页面上使用`升级`功能进行滚动升级, 数秒后刷新页面即可。

*注: 更新过程中如果出错, 请查看脚本信息或联系[maxwell](github.com/maxwell92))*


### 以下为旧版本, 无效///////////////////////////////////////////////////

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

