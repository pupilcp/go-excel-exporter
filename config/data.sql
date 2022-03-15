CREATE TABLE `download_task` (
     `task_id` int(11) NOT NULL AUTO_INCREMENT,
     `user_id` int(10) unsigned NOT NULL DEFAULT '0',
     `file_name` varchar(32) NOT NULL DEFAULT '' COMMENT '文件名',
     `request_url` varchar(255) NOT NULL DEFAULT '' COMMENT '请求url',
     `request_params` text NOT NULL COMMENT '请求参数',
     `request_method` varchar(8) NOT NULL DEFAULT '' COMMENT '请求方法，POST/GET',
     `download_file` varchar(255) NOT NULL DEFAULT '' COMMENT '下载文件地址',
     `task_status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '任务状态，0：未处理，1：处理中，2：处理成功，3：处理失败',
     `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注信息',
     `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '任务创建时间',
     `updated_at` int(11) NOT NULL DEFAULT '0' COMMENT '任务更新时间',
     PRIMARY KEY (`task_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='下载任务表';