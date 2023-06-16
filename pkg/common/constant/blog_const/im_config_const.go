package blog_const

import "time"

var (
	/**
	 * 默认群组ID，注册时自动加入
	 */
	DEFAULT_GROUP_ID int32 = -1

	/**
	 * 系统消息ID
	 */
	DEFAULT_SYSTEM_MESSAGE_ID = -1

	/**
	 * 协议名字
	 */
	PROTOCOL_NAME = "protocol_poetize"

	CHARSET = "UTF-8"

	/**
	 * 监听端口
	 */
	SERVER_PORT = 9324

	/**
	 * 心跳超时时间，单位：毫秒
	 */
	HEARTBEAT_TIMEOUT = 1000 * 600

	/**
	 * IP数据监控统计，时间段
	 */
	DURATION_DEFAULT  = int64(time.Minute * 5)
	IP_STAT_DURATIONS = []int64{int64(time.Minute * 5)}

	/**
	 * 默认群聊
	 */
	GROUP_DEFAULT = "group_default"

	/**
	 * 加入群
	 * <p>
	 * 0：无需验证
	 * 1：需要群主或管理员同意
	 */
	IN_TYPE_FALSE = false
	IN_TYPE_TRUE  = true

	/**
	 * 用户状态[-1:审核不通过或者踢出群聊，0:未审核，1:审核通过，2:禁言]
	 */
	GROUP_USER_STATUS_BAN        int8 = -1
	GROUP_USER_STATUS_NOT_VERIFY int8 = 0
	GROUP_USER_STATUS_PASS       int8 = 1
	GROUP_USER_STATUS_SILENCE    int8 = 2

	/**
	 * 朋友状态[-1:审核不通过或者删除好友，0:未审核，1:审核通过]
	 */
	FRIEND_STATUS_BAN        int8 = -1
	FRIEND_STATUS_NOT_VERIFY int8 = 0
	FRIEND_STATUS_PASS       int8 = 1

	/**
	 * 是否已读
	 * <p>
	 * 0：未读
	 * 1：已读
	 */
	USER_MESSAGE_STATUS_FALSE = false
	USER_MESSAGE_STATUS_TRUE  = true

	/**
	 * 是否是群组管理员
	 * <p>
	 * 0：否
	 * 1：是
	 */
	ADMIN_FLAG_FALSE = false
	ADMIN_FLAG_TRUE  = true

	/**
	 * 群类型[1:聊天群，2:话题]
	 */
	GROUP_COMMON = 1
	GROUP_TOPIC  = 2
)
