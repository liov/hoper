DROP DATABASE IF EXISTS hoper;
CREATE DATABASE hoper;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`                bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`              varchar(15)         NULL DEFAULT NULL COMMENT '昵称',
    `activated_at`      datetime(3)         NULL DEFAULT NULL,
    `account`           varchar(15)         NULL DEFAULT NULL COMMENT '账号',
    `password`          varchar(63)         NULL DEFAULT NULL,
    `email`             varchar(31)         NULL DEFAULT NULL,
    `phone`             varchar(31)         NULL DEFAULT NULL,
    `gender`            tinyint(1)          NULL DEFAULT NULL,
    `birthday`          datetime            NULL DEFAULT NULL,
    `introduction`      varchar(255)        NULL DEFAULT NULL,
    `score`             bigint(20) UNSIGNED NULL DEFAULT NULL,
    `signature`         varchar(255)        NULL DEFAULT NULL,
    `role`              tinyint(1) UNSIGNED NULL DEFAULT NULL,
    `avatar_url`        varchar(255)        NULL DEFAULT NULL,
    `cover_url`         varchar(255)        NULL DEFAULT NULL,
    `address`           varchar(255)        NULL DEFAULT NULL,
    `location`          varchar(255)        NULL DEFAULT NULL,
    `updated_at`        datetime(3)         NULL DEFAULT NULL,
    `banned_at`         datetime(3)         NULL DEFAULT NULL,
    `created_at`        datetime(3)         NULL DEFAULT NULL,
    `last_activated_at` datetime(3)         NULL DEFAULT NULL,
    `status`            tinyint(1) UNSIGNED NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT '用户表';

DROP TABLE IF EXISTS `user_score_log`;
CREATE TABLE `user_score_log`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`    bigint(20) UNSIGNED NOT NULL COMMENT 'user_id',
    `change`     tinyint(1) UNSIGNED NOT NULL COMMENT '0：减少，1：增加',
    `score`      bigint(20) UNSIGNED NOT NULL COMMENT '变化的分数',
    `reason`     varchar(255)        NULL DEFAULT NULL COMMENT '变化的原因',
    `created_at` datetime(3)         NULL DEFAULT NULL COMMENT '时间',
    `remark`     varchar(255)        NULL DEFAULT NULL COMMENT '备注'

) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT '用户分数日志表';

DROP TABLE IF EXISTS `user_banned_log`;
CREATE TABLE `user_banned_log`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`    bigint(20) UNSIGNED NOT NULL COMMENT 'user_id',
    `duration`   bigint(20) UNSIGNED NOT NULL COMMENT '持续时间',
    `reason`     varchar(255)        NULL DEFAULT NULL COMMENT '变化的原因',
    `created_at` datetime(3)         NULL DEFAULT NULL COMMENT '时间',
    `remark`     varchar(255)        NULL DEFAULT NULL COMMENT '备注'

) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT '用户封禁日志表';

DROP TABLE IF EXISTS `user_action_log`;
CREATE TABLE `user_action_log`
(
    `id`         bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`    bigint(20) UNSIGNED NOT NULL COMMENT 'user_id',
    `action`     tinyint(1) UNSIGNED NOT NULL COMMENT '操作，0：减少，1：增加',
    `location`   varchar(255)        NULL DEFAULT NULL COMMENT '操作地点',
    `ip`         varchar(31)         NULL DEFAULT NULL COMMENT '操作IP',
    `device`     varchar(31)         NULL DEFAULT NULL COMMENT '操作设备',
    `created_at` datetime(3)         NULL DEFAULT NULL COMMENT '时间',
    `remark`     varchar(255)        NULL DEFAULT NULL COMMENT '备注'

) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_general_ci
    COMMENT '用户操作日志表';


#set global max_allowed_packet = 200*1024*1024
#https://gitee.com/indanimp/china_area_mysql
/*DROP TABLE IF EXISTS `cnarea_2018`;

CREATE TABLE `cnarea_2018`
(
    `id`          mediumint(7) unsigned          NOT NULL AUTO_INCREMENT,
    `level`       tinyint(1) unsigned            NOT NULL COMMENT '层级',
    `parent_code` bigint(14) unsigned            NOT NULL DEFAULT '0' COMMENT '父级行政代码',
    `area_code`   bigint(14) unsigned            NOT NULL DEFAULT '0' COMMENT '行政代码',
    `zip_code`    mediumint(6) unsigned zerofill NOT NULL DEFAULT '000000' COMMENT '邮政编码',
    `city_code`   char(6)                        NOT NULL DEFAULT '' COMMENT '区号',
    `name`        varchar(50)                    NOT NULL DEFAULT '' COMMENT '名称',
    `short_name`  varchar(50)                    NOT NULL DEFAULT '' COMMENT '简称',
    `merger_name` varchar(50)                    NOT NULL DEFAULT '' COMMENT '组合名',
    `pinyin`      varchar(30)                    NOT NULL DEFAULT '' COMMENT '拼音',
    `lng`         decimal(10, 6)                 NOT NULL DEFAULT '0.000000' COMMENT '经度',
    `lat`         decimal(10, 6)                 NOT NULL DEFAULT '0.000000' COMMENT '纬度',
    PRIMARY KEY (`id`),
    KEY `name` (`name`),
    KEY `level` (`level`),
    KEY `area_code` (`area_code`),
    KEY `parent_code` (`parent_code`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;*/

#https://raw.githubusercontent.com/adyliu/china_area/master/area_code_2019.sql.gz
/*DROP TABLE IF EXISTS `area_code_2019`;

CREATE TABLE `area_code_2019`
(
    `code`  bigint(12) unsigned NOT NULL COMMENT '区划代码',
    `name`  varchar(128)        NOT NULL DEFAULT '' COMMENT '名称',
    `level` tinyint(1)          NOT NULL COMMENT '级别1-5,省市县镇村',
    `pcode` bigint(12)                   DEFAULT NULL COMMENT '父级区划代码',
    PRIMARY KEY (`code`),
    KEY `name` (`name`),
    KEY `level` (`level`),
    KEY `pcode` (`pcode`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

CREATE VIEW area_index_2019 AS
SELECT a.code,e.name AS province,d.name AS city  ,c.name AS county,b.name AS town,a.name AS villagetr
FROM area_code_2019 a
         JOIN area_code_2019 b ON a.level=5 AND b.level=4 AND a.pcode=b.code
         JOIN area_code_2019 c ON b.pcode=c.code
         JOIN area_code_2019 d ON c.pcode=d.code
         JOIN area_code_2019 e ON d.pcode=e.code
ORDER BY a.code*/

CREATE TABLE `cnarea`
(
    `area_code`   bigint(14) unsigned            NOT NULL DEFAULT '0' COMMENT '行政代码',
    `name`        varchar(50)                    NOT NULL DEFAULT '' COMMENT '名称',
    `level`       tinyint(1) unsigned            NOT NULL COMMENT '层级',
    `parent_code` bigint(14) unsigned            NOT NULL DEFAULT '0' COMMENT '父级行政代码',
    `zip_code`    mediumint(6) unsigned zerofill NOT NULL DEFAULT '000000' COMMENT '邮政编码',
    `city_code`   char(6)                        NOT NULL DEFAULT '' COMMENT '区号',
    `short_name`  varchar(50)                    NOT NULL DEFAULT '' COMMENT '简称',
    `merger_name` varchar(50)                    NOT NULL DEFAULT '' COMMENT '组合名',
    `pinyin`      varchar(30)                    NOT NULL DEFAULT '' COMMENT '拼音',
    `lng`         decimal(10, 6)                 NOT NULL DEFAULT '0.000000' COMMENT '经度',
    `lat`         decimal(10, 6)                 NOT NULL DEFAULT '0.000000' COMMENT '纬度',
    PRIMARY KEY (`area_code`),
    KEY `name` (`name`),
    KEY `level` (`level`),
    KEY `area_code` (`area_code`),
    KEY `parent_code` (`parent_code`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

insert into cnarea (area_code, name, level, parent_code, zip_code, city_code, short_name, merger_name, pinyin, lng, lat)
select area_code,
       name,
       level,
       parent_code,
       zip_code,
       city_code,
       short_name,
       merger_name,
       pinyin,
       lng,
       lat
from cnarea_2018;