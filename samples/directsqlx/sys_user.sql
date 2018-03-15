/*
MySQL Data Transfer

Source Server         : MysqlLocal
Source Server Version : 50712
Source Host           : localhost:3306
Source Database       : faygo

Target Server Type    : MYSQL
Target Server Version : 50712
File Encoding         : 65001

Date: 2016-12-31 14:51:59
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
  `id` char(36) NOT NULL,
  `code` varchar(50) NOT NULL,
  `cnname` varchar(100) NOT NULL,
  `enname` varchar(100) DEFAULT NULL,
  `nick` varchar(50) DEFAULT NULL,
  `pwd` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `UQE_sys_user_id` (`id`),
  KEY `IDX_sys_user_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO `sys_user` VALUES ('001', 'admin', 'admin', null, '我是畅雨！！！',  '123456');
INSERT INTO `sys_user` VALUES ('0ef51ef5-82a3-4d2c-8a85-e7fb1306991a', 'zs', '張三', null, '張三',  '123');
INSERT INTO `sys_user` VALUES ('191ccbc0-3732-4985-a469-ac4a4f7ed024', 'll', '李立', null, '李立', '123');
INSERT INTO `sys_user` VALUES ('191cc3c0-3732-4985-a469-ac4a4f7ed024', 'ww', '王五', null, '王五', '123');