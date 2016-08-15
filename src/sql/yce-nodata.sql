-- MySQL dump 10.13  Distrib 5.5.50, for debian-linux-gnu (x86_64)
--
-- Host: 172.21.1.11    Database: yce
-- ------------------------------------------------------
-- Server version	5.7.13

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `datacenter`
--

DROP TABLE IF EXISTS `datacenter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `datacenter` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '数据中心名称',
  `host` varchar(256) NOT NULL COMMENT 'APIServer host',
  `port` int(10) NOT NULL COMMENT 'APIServer 端口',
  `secret` text COMMENT 'TLS证书',
  `status` int(10) NOT NULL COMMENT '1启用/0弃用',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间戳',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间戳',
  `modifiedOp` int(10) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='数据中心表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dcquota`
--

DROP TABLE IF EXISTS `dcquota`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dcquota` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `orgId` int(10) NOT NULL COMMENT '组织ID',
  `dcId` int(10) NOT NULL COMMENT '数据中心ID',
  `podNumLimit` int(10) NOT NULL COMMENT '限制Pod的数量',
  `podCpuMax` int(10) NOT NULL COMMENT '每个Pod使用的最大CPU（单位core）',
  `podMemMax` int(10) NOT NULL COMMENT '每个Pod使用的最大内存（单位G）',
  `podCpuMin` int(10) NOT NULL COMMENT '每个Pod使用的最小CPU（单位core）',
  `podMemMin` int(10) NOT NULL COMMENT '每个Pod使用的最小内存（单位G）',
  `rbdQuota` int(10) NOT NULL COMMENT '最大云盘多少G',
  `podRbdMax` int(10) NOT NULL COMMENT '每块云盘的最大限制（单位G）',
  `podRbdMin` int(10) NOT NULL COMMENT '每块云盘的最小限制（单位G）',
  `price` varchar(256) NOT NULL COMMENT '金额用string代替decimal(15,4)',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间戳',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间戳',
  `modifiedOp` int(10) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`),
  KEY `FK_dcquota_1` (`orgId`),
  KEY `FK_dcquota_2` (`dcId`),
  CONSTRAINT `FK_dcquota_1` FOREIGN KEY (`orgId`) REFERENCES `organization` (`id`),
  CONSTRAINT `FK_dcquota_2` FOREIGN KEY (`dcId`) REFERENCES `datacenter` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='数据中心配额表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deployment`
--

DROP TABLE IF EXISTS `deployment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deployment` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name` varchar(256) NOT NULL COMMENT '应用名',
  `actionType` int(10) NOT NULL COMMENT '操作类型（上线，回滚等）',
  `actioVerb` varchar(256) NOT NULL COMMENT 'GET/POST/DELETE',
  `actionUrl` varchar(256) NOT NULL COMMENT '操作的URL',
  `actionAt` varchar(256) NOT NULL COMMENT '操作时间戳',
  `actionOp` int(10) NOT NULL COMMENT '操作人员',
  `dcList` varchar(256) NOT NULL COMMENT '操作的数据中心（以：分割）',
  `success` int(10) NOT NULL COMMENT '操作是否成功',
  `reason` varchar(256) NOT NULL COMMENT '出错原因',
  `json` text NOT NULL COMMENT '存储Json文件内容',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`),
  KEY `FK_deployment_1` (`actionType`),
  CONSTRAINT `FK_deployment_1` FOREIGN KEY (`actionType`) REFERENCES `option` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='应用发布日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `option`
--

DROP TABLE IF EXISTS `option`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `option` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '操作类型名称',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间戳',
  `modifiedOp` varchar(256) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='操作类型属性表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `organization`
--

DROP TABLE IF EXISTS `organization`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `organization` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(256) NOT NULL COMMENT '组织名称（默认是namespace）',
  `cpuQuota` int(10) NOT NULL COMMENT '多少个核，不支持小数',
  `memQuota` int(10) NOT NULL COMMENT '多少G的内存，不支持小数',
  `budget` varchar(256) NOT NULL COMMENT '年度预算（元）用string代替decimal(15,4)',
  `balance` varchar(256) NOT NULL COMMENT '剩余余额（元）用string代替decimal(15,4)',
  `status` int(10) NOT NULL COMMENT '1启用/0弃用',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间',
  `modifiedOp` int(10) NOT NULL COMMENT '最后修改人员',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='组织表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `quota`
--

DROP TABLE IF EXISTS `quota`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quota` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '名称',
  `cpu` int(10) NOT NULL COMMENT 'CPU(单位core)',
  `mem` int(10) NOT NULL COMMENT '内存(单位G)',
  `rbd` int(10) NOT NULL COMMENT '云盘(单位G)',
  `price` varchar(256) NOT NULL COMMENT '人民币/每年用string代替decimal(15,4)',
  `status` int(10) NOT NULL COMMENT '1启用/0弃用',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间戳',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间戳',
  `modifiedOp` varchar(256) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='配额表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rbd`
--

DROP TABLE IF EXISTS `rbd`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rbd` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `name` varchar(256) NOT NULL COMMENT '云盘的名称（rbd image 名称）',
  `pool` varchar(256) NOT NULL COMMENT 'rbd image所属的pool池',
  `size` int(10) NOT NULL COMMENT '云盘大小（单位G）',
  `filesystem` varchar(256) NOT NULL COMMENT '文件系统（默认ext4，目前只支持ext4）',
  `orgId` int(10) NOT NULL COMMENT '关联的组织',
  `dcId` int(10) NOT NULL COMMENT '关联的数据中心',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间戳',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间戳',
  `modifiedOp` varchar(256) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`),
  KEY `FK_rbd_1` (`orgId`),
  KEY `FK_rbd_2` (`dcId`),
  CONSTRAINT `FK_rbd_1` FOREIGN KEY (`orgId`) REFERENCES `organization` (`id`),
  CONSTRAINT `FK_rbd_2` FOREIGN KEY (`dcId`) REFERENCES `datacenter` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='云盘表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长主键',
  `name` varchar(256) NOT NULL COMMENT '组织名称（默认是namespace）',
  `password` varchar(256) NOT NULL COMMENT '密码（加盐加密）',
  `orgId` int(10) NOT NULL COMMENT '组织ID',
  `status` int(10) NOT NULL COMMENT '1启用/0弃用',
  `createdAt` varchar(256) NOT NULL COMMENT '创建时间（2016-07-22T10:20:30Z）',
  `modifiedAt` varchar(256) NOT NULL COMMENT '最后修改时间',
  `modifedOp` int(10) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`),
  KEY `FK_user_1` (`orgId`),
  CONSTRAINT `FK_user_1` FOREIGN KEY (`orgId`) REFERENCES `organization` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-08-15 17:35:29
