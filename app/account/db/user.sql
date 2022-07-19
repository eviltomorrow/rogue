-- create table user
drop table if exists `rogue_account`.`user`;
create table `rogue_account`.`user` (
    `id` BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `uuid` CHAR(36) NOT NULL COMMENT 'uuid',
    `nick_name` VARCHAR(32) COMMENT '昵称',
    `email` VARCHAR(32) NOT NULL COMMENT '邮箱',
    `phone` VARCHAR(32) NOT NULL COMMENT '手机号码',
    `disabled` TINYINT NOT NULL COMMENT '是否禁用',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间',
     UNIQUE KEY `uq_email` (`email`), UNIQUE KEY `uq_phone` (`phone`)
);
