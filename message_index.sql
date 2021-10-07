CREATE TABLE `t_message_index` (
  	`shard_id` INT NOT NULL  COMMENT '分片标识', 
	`account_a` VARCHAR(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '队列唯一标识', 
	`account_b` VARCHAR(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '另一方', 
	`direction` tinyint UNSIGNED NOT NULL DEFAULT '0' COMMENT '1表示AccountA为发送者', 
	`message_id` BIGINT NOT NULL COMMENT '关联消息内容表中的ID', 
	`group` VARCHAR(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '群ID，单聊情况为空', 
	`send_time` BIGINT NOT NULL COMMENT '消息发送时间', 
	KEY `idx_t_message_index_account_a` (`account_a`), 
	KEY `idx_t_message_index_send_time` (`send_time`)
) 
ENGINE = INNODB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci
PARTITION BY HASH(shard_id)
PARTITIONS 50;