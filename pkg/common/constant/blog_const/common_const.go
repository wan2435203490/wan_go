package blog_const

import "time"

const (
	ADMIN_USER_ID        = 1
	USER_TOKEN           = "user_token_"
	ADMIN_TOKEN          = "admin_token_"
	USER_TOKEN_INTERVAL  = "user_token_interval_"
	ADMIN_TOKEN_INTERVAL = "admin_token_interval_"
	USER_ACCESS_TOKEN    = "user_access_token_"
	ADMIN_ACCESS_TOKEN   = "admin_access_token_"
	TOKEN_HEADER         = "Authorization"
	TOKEN_EXPIRE         = time.Hour * 24
	TOKEN_INTERVAL       = time.Hour
	ADMIN                = "admin"
	ADMIN_FAMILY         = "adminFamily"
	FAMILY_LIST          = "familyList"
	COMMENT_IM_MAIL      = "comment_im_mail_"

	COMMENT_IM_MAIL_COUNT int32 = 1

	/**
	 * 验证码
	 */
	USER_CODE = "user_code_"

	/**
	 * 忘记密码时获取验证码用于找回密码
	 */
	FORGET_PASSWORD = "forget_password_"

	/**
	 * 网站信息
	 */
	WEB_INFO = "webInfo"

	/**
	 * 分类信息
	 */
	SORT_INFO = "sortInfo"

	/**
	 * 赞赏
	 */
	ADMIRE = "admire"

	/**
	 * 根据用户ID获取用户信息
	 */
	USER_CACHE = "user_"

	/**
	 * 根据文章ID获取评论数量
	 */
	COMMENT_COUNT_CACHE = "comment_count_"

	/**
	 * 根据用户ID获取该用户所有文章ID
	 */
	USER_ARTICLE_LIST = "user_article_list_"

	/**
	 * 默认缓存过期时间
	 */
	EXPIRE = 1800

	/**
	 * 树洞一次最多查询条数
	 */
	TREE_HOLE_COUNT = 200

	/**
	家庭信息默认查询条数
	*/
	FAMILY_COUNT = 10

	/**
	 * 顶层评论ID
	 */
	FIRST_COMMENT = 0

	/**
	 * 文章摘要默认字数
	 */
	SUMMARY = 80

	/**
	 * 留言的源
	 */
	TREE_HOLE_COMMENT_SOURCE = 0

	/**
	 * 资源类型
	 */
	PATH_TYPE_GRAFFITI             = "graffiti"
	PATH_TYPE_ARTICLE_PICTURE      = "articlePicture"
	PATH_TYPE_USER_AVATAR          = "userAvatar"
	PATH_TYPE_ARTICLE_COVER        = "articleCover"
	PATH_TYPE_WEB_BACKGROUND_IMAGE = "webBackgroundImage"
	PATH_TYPE_WEB_AVATAR           = "webAvatar"
	PATH_TYPE_RANDOM_AVATAR        = "randomAvatar"
	PATH_TYPE_RANDOM_COVER         = "randomCover"
	PATH_TYPE_COMMENT_PICTURE      = "commentPicture"
	PATH_TYPE_INTERNET_MEME        = "internetMeme"
	PATH_TYPE_IM_GROUP_AVATAR      = "im/groupAvatar"
	PATH_TYPE_IM_GROUP_MESSAGE     = "im/groupMessage"
	PATH_TYPE_IM_FRIEND_MESSAGE    = "im/friendMessage"
	PATH_TYPE_FUNNY_URL            = "funnyUrl"
	PATH_TYPE_FUNNY_COVER          = "funnyCover"
	PATH_TYPE_FAVORITES_COVER      = "favoritesCover"
	PATH_TYPE_LOVE_COVER           = "love/bgCover"
	PATH_TYPE_LOVE_MAN             = "love/manCover"
	PATH_TYPE_LOVE_WOMAN           = "love/womanCover"

	/**
	 * 资源路径
	 */
	RESOURCE_PATH_TYPE_FRIEND     = "friendUrl"
	RESOURCE_PATH_TYPE_FUNNY      = "funny"
	RESOURCE_PATH_TYPE_FAVORITES  = "favorites"
	RESOURCE_PATH_TYPE_LOVE_PHOTO = "lovePhoto"

	/**
	 * 微言
	 */
	WEIYAN_TYPE_FRIEND = "friend"
	WEIYAN_TYPE_NEWS   = "news"
)
