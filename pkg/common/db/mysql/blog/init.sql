
INSERT INTO `user` (`id`, `username`, `password`, `phone_number`, `email`, `user_status`, `gender`, `open_id`, `avatar`, `introduction`, `user_type`, `update_by`, `deleted`, `admire`) VALUES(1, 'Sara', '47bce5c74f589f4867dbd57e9ca9f808', '', '', 1, 1, '', '', '', 0, 'Sara', 0, '');

INSERT INTO `web_info` (`id`, `web_name`, `web_title`, `notices`, `footer`, `background_image`, `avatar`, `random_avatar`, `random_name`, `random_cover`, `waifu_json`, `status`) VALUES(1, 'Sara', '寻国记', '[]', '云想衣裳花想容， 春风拂槛露华浓。', '', '', '[]', '[]', '[]', '{}', 1);

INSERT INTO `im_chat_group` (`id`, `group_name`, `master_user_id`, `introduction`, `notice`, `in_type`) VALUES(-1, '公共聊天室', 1, '公共聊天室', '欢迎光临！', 0);

INSERT INTO `im_chat_group_user` (`id`, `group_id`, `user_id`, `admin_flag`, `user_status`) VALUES(1, -1, 1, 1, 1);
