
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