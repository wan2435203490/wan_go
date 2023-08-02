
INSERT INTO `user` (`id`, `user_name`, `password`, `phone_number`, `email`, `user_status`, `gender`, `open_id`, `avatar`, `introduction`, `user_type`, `update_by`, `admire`) VALUES(1, 'Wan', '47bce5c74f589f4867dbd57e9ca9f808', '', '', 1, 1, '', '', '', 0, 'Wan', '');

INSERT INTO `web_info` (`id`, `web_name`, `web_title`, `notices`, `footer`, `background_image`, `avatar`, `random_avatar`, `random_name`, `random_cover`, `waifu_json`, `status`) VALUES(1, 'Wan', '小世界', '[]', '云想衣裳花想容， 春风拂槛露华浓。', '', '', '[]', '[]', '[]', '{}', 1);

INSERT INTO `im_chat_group` (`id`, `group_name`, `master_user_id`, `introduction`, `notice`, `in_type`) VALUES(-1, '公共聊天室', 1, '公共聊天室', '欢迎光临！', 0);

INSERT INTO `im_chat_group_user` (`id`, `group_id`, `user_id`, `admin_flag`, `user_status`) VALUES(1, -1, 1, 1, 1);

INSERT INTO `poetize`.`family`(`id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `bg_cover`, `man_cover`, `woman_cover`, `man_name`, `woman_name`, `timing`, `countdown_title`, `countdown_time`, `status`, `family_info`, `like_count`) VALUES (1, '2023-07-02 21:44:14.000', NULL, NULL, 1, 'https://api.7585.net.cn/bing/api.php?rand=1', 'https://www.2fun.top/love/avatar/37.jpeg', 'https://www.2fun.top/love/avatar/38.jpg', '亚瑟', '安琪拉', '2022-5-16', '下个活动', '2023-8-8', 1, 'qwq', 99);