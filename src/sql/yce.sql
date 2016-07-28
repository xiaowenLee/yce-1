-- MySQL dump 10.13  Distrib 5.5.50, for debian-linux-gnu (x86_64)
--
-- Host: 192.168.3.216    Database: yce
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
-- Table structure for table `dcquotas`
--

DROP TABLE IF EXISTS `dcquotas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dcquotas` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `org_id` int(10) NOT NULL COMMENT '组织ID',
  `dc_id` int(10) NOT NULL COMMENT '数据中心ID',
  `pod_num_limit` int(10) NOT NULL COMMENT '限制Pod的数量',
  `pod_cpu_max` int(10) NOT NULL COMMENT '每个Pod使用的最大CPU（单位core）',
  `pod_mem_max` int(10) NOT NULL COMMENT '每个Pod使用的最大内存（单位G）',
  `pod_cpu_min` int(10) NOT NULL COMMENT '每个Pod使用的最小CPU（单位core）',
  `pod_mem_min` int(10) NOT NULL COMMENT '每个Pod使用的最小内存（单位G）',
  `rbd_quota` int(10) NOT NULL COMMENT '最大云盘多少G',
  `pod_rbd_max` int(10) NOT NULL COMMENT '每块云盘的最大限制（单位G）',
  `pod_rbd_min` int(10) NOT NULL COMMENT '每块云盘的最小限制（单位G）',
  `price` double(4,2) NOT NULL COMMENT '金额',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='数据中心配额表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `dcs`
--

DROP TABLE IF EXISTS `dcs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `dcs` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '数据中心名称',
  `host` varchar(256) NOT NULL COMMENT 'APIServer host',
  `port` int(10) NOT NULL COMMENT 'APIServer 端口',
  `secret` text COMMENT 'TLS证书',
  `created_ts` datetime NOT NULL COMMENT '创建时间戳',
  `last_modified_ts` datetime NOT NULL COMMENT '最后修改时间戳',
  `last_modified_op` int(10) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='数据中心表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `deployments`
--

DROP TABLE IF EXISTS `deployments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `deployments` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL,
  `action_type` int(10) NOT NULL,
  `action_verb` varchar(256) NOT NULL,
  `action_url` varchar(256) NOT NULL,
  `action_ts` datetime NOT NULL,
  `action_op` int(10) NOT NULL,
  `dc_list` varchar(256) NOT NULL,
  `success` int(10) NOT NULL COMMENT '操作是否成功，1为成功，0为失败',
  `reason` varchar(256) NOT NULL,
  `json` text NOT NULL,
  `comment` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='应用发布日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `options`
--

DROP TABLE IF EXISTS `options`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `options` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `created_ts` datetime NOT NULL,
  `last_modified_ts` datetime NOT NULL,
  `last_modified_op` int(10) NOT NULL,
  `comment` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='操作类型属性表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orgs`
--

DROP TABLE IF EXISTS `orgs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `orgs` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '组织名称',
  `cpu_quota` int(10) NOT NULL COMMENT '多少个核，不支持小数',
  `mem_quota` int(10) NOT NULL COMMENT '多少G的内存，不支持小数',
  `budget` int(10) NOT NULL COMMENT '年度预算（元）',
  `balance` int(10) NOT NULL COMMENT '剩余余额（元）',
  `created_ts` datetime NOT NULL COMMENT '创建时间',
  `last_modified_ts` datetime NOT NULL COMMENT '最后修改时间',
  `last_modified_op` int(10) NOT NULL COMMENT '最后修改人员',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='组织表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `quotas`
--

DROP TABLE IF EXISTS `quotas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quotas` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '标准配额名称',
  `cpu` int(10) NOT NULL COMMENT 'CPU(单位core)',
  `mem` int(10) NOT NULL COMMENT '内存(单位G)',
  `rbd` int(10) NOT NULL COMMENT '云盘(单位G)',
  `price` float(4,2) NOT NULL COMMENT '人民币/每年',
  `created_ts` datetime NOT NULL COMMENT '创建时间戳',
  `last_modified_ts` datetime NOT NULL COMMENT '最后修改时间戳',
  `last_modified_op` int(10) DEFAULT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='标准配额表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `rbds`
--

DROP TABLE IF EXISTS `rbds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rbds` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长ID',
  `name` varchar(256) NOT NULL COMMENT '云盘的名称（rbd image 名称）',
  `pool` varchar(256) NOT NULL COMMENT 'rbd image所属的pool池',
  `size` int(10) NOT NULL COMMENT '云盘大小（单位G）',
  `filesystem` varchar(256) NOT NULL COMMENT '文件系统（默认ext4，目前只支持ext4）',
  `org_id` int(10) NOT NULL COMMENT '关联的组织',
  `dc_id` int(10) NOT NULL COMMENT '关联的数据中心',
  `created_ts` datetime NOT NULL COMMENT '创建时间戳',
  `last_modified_ts` datetime NOT NULL COMMENT '最后修改时间戳',
  `last_modified_op` int(10) NOT NULL COMMENT '最后修改操作人',
  `comments` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COMMENT='云盘表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '自增长主键',
  `name` varchar(256) NOT NULL COMMENT '用户名',
  `password` varchar(256) NOT NULL COMMENT '密码',
  `org_id` int(10) NOT NULL COMMENT '组织ID',
  `created_ts` datetime NOT NULL COMMENT '创建时间',
  `last_modified_ts` datetime NOT NULL COMMENT '最后修改时间',
  `last_modifed_op` int(10) NOT NULL COMMENT '最后修改操作人',
  `comment` varchar(256) DEFAULT NULL COMMENT '说明',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-07-26 13:54:12
