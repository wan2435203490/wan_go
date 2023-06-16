package blog

import (
	"gorm.io/gorm"
	"wan_go/pkg/vo/blog"
)

//todo tinyint or bit

// 用户信息表
type User struct {
	gorm.Model
	ID           int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;primaryKey;NOT NULL"`
	UserName     string `gorm:"column:user_name;type:VARCHAR(32);"`
	Password     string `gorm:"column:password;type:VARCHAR(128);"`
	PhoneNumber  string `gorm:"column:phone_number;type:VARCHAR(16);"`
	Email        string `gorm:"column:email;type:VARCHAR(32);"`
	UserStatus   bool   `gorm:"column:user_status;type:TINYINT(1);NOT NULL"`
	Gender       int8   `gorm:"column:gender;type:TINYINT(2);"`
	OpenId       string `gorm:"column:open_id;type:VARCHAR(128);"`
	Avatar       string `gorm:"column:avatar;type:VARCHAR(256);"`
	Admire       string `gorm:"column:admire;type:VARCHAR(32);"`
	Introduction string `gorm:"column:introduction;type:VARCHAR(4096);"`
	UserType     int8   `gorm:"column:user_type;type:TINYINT(2);NOT NULL"`

	//CreateTime time.Time `gorm:"column:create_time;type:DATETIME;"`
	//UpdateTime time.Time `gorm:"column:update_time;type:DATETIME;"`
	UpdateBy string `gorm:"column:update_by;type:VARCHAR(32);"`
	//Deleted    bool      `gorm:"column://Deleted;type:TINYINT(1);NOT NULL"`
}

// 文章表
type Article struct {
	gorm.Model
	ID              int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	UserId          int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	SortId          int32  `gorm:"column:sort_id;type:INT;NOT NULL"`
	LabelId         int32  `gorm:"column:label_id;type:INT;NOT NULL"`
	ArticleCover    string `gorm:"column:article_cover;type:VARCHAR(256);"`
	ArticleTitle    string `gorm:"column:article_title;type:VARCHAR(32);NOT NULL"`
	ArticleContent  string `gorm:"column:article_content;type:TEXT;NOT NULL"`
	ViewCount       int32  `gorm:"column:view_count;type:INT;NOT NULL"`
	LikeCount       int32  `gorm:"column:like_count;type:INT;NOT NULL"`
	ViewStatus      bool   `gorm:"column:view_status;type:TINYINT(1);NOT NULL"`
	Password        string `gorm:"column:password;type:VARCHAR(128);"`
	RecommendStatus bool   `gorm:"column:recommend_status;type:TINYINT(1);NOT NULL"`
	CommentStatus   bool   `gorm:"column:comment_status;type:TINYINT(1);NOT NULL"`
	//CreateTime      time.Time `gorm:"column:create_time;type:DATETIME;"`
	//UpdateTime      time.Time `gorm:"column:update_time;type:DATETIME;"`
	UpdateBy string `gorm:"column:update_by;type:VARCHAR(32);"`
	//Deleted         int8      `gorm:"column://Deleted;type:TINYINT(1);NOT NULL"`
}

// 文章评论表
type Comment struct {
	gorm.Model
	ID              int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	Source          int32  `gorm:"column:source;type:INT;NOT NULL"`
	Type            string `gorm:"column:type;type:VARCHAR(32);NOT NULL"`
	ParentCommentId int32  `gorm:"column:parent_comment_id;type:INT;NOT NULL"`
	UserId          int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	FloorCommentId  int32  `gorm:"column:floor_comment_id;type:INT;"`
	ParentUserId    int32  `gorm:"column:parent_user_id;type:INT;"`
	LikeCount       int32  `gorm:"column:like_count;type:INT;NOT NULL"`
	CommentContent  string `gorm:"column:comment_content;type:VARCHAR(1024);NOT NULL"`
	CommentInfo     string `gorm:"column:comment_info;type:VARCHAR(256);"`
	//CreateTime      time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 分类
type Sort struct {
	ID              int32    `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	SortName        string   `gorm:"column:sort_name;type:VARCHAR(32);NOT NULL"`
	SortDescription string   `gorm:"column:sort_description;type:VARCHAR(256);NOT NULL"`
	SortType        int8     `gorm:"column:sort_type;type:TINYINT(2);NOT NULL"`
	Priority        int32    `gorm:"column:priority;type:INT;"`
	CountOfSort     int32    `gorm:"-:all"`
	Labels          []*Label `gorm:"-:all"`
}

// 标签
type Label struct {
	ID               int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	SortId           int32  `gorm:"column:sort_id;type:INT;NOT NULL"`
	LabelName        string `gorm:"column:label_name;type:VARCHAR(32);NOT NULL"`
	LabelDescription string `gorm:"column:label_description;type:VARCHAR(256);NOT NULL"`
	CountOfLabel     int32  `gorm:"-:all"`
}

// 树洞
type TreeHole struct {
	gorm.Model
	ID      int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	Avatar  string `gorm:"column:avatar;type:VARCHAR(256);"`
	Message string `gorm:"column:message;type:VARCHAR(64);NOT NULL"`
	//CreateTime time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 微言表
type WeiYan struct {
	gorm.Model
	ID        int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	UserId    int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	LikeCount int32  `gorm:"column:like_count;type:INT;NOT NULL"`
	Content   string `gorm:"column:content;type:VARCHAR(1024);NOT NULL"`
	Type      string `gorm:"column:type;type:VARCHAR(32);NOT NULL"`
	Source    int32  `gorm:"column:source;type:INT;"`
	IsPublic  bool   `gorm:"column:is_public;type:TINYINT(1);NOT NULL"`
	//CreateTime time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 网站信息表
type WebInfo struct {
	ID              int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	WebName         string `gorm:"column:web_name;type:VARCHAR(16);NOT NULL"`
	WebTitle        string `gorm:"column:web_title;type:VARCHAR(512);NOT NULL"`
	Notices         string `gorm:"column:notices;type:VARCHAR(512);"`
	Footer          string `gorm:"column:footer;type:VARCHAR(256);NOT NULL"`
	BackgroundImage string `gorm:"column:background_image;type:VARCHAR(256);"`
	Avatar          string `gorm:"column:avatar;type:VARCHAR(256);NOT NULL"`
	RandomAvatar    string `gorm:"column:random_avatar;type:TEXT;"`
	RandomName      string `gorm:"column:random_name;type:VARCHAR(4096);"`
	RandomCover     string `gorm:"column:random_cover;type:TEXT;"`
	WaifuJson       string `gorm:"column:waifu_json;type:TEXT;"`
	Status          bool   `gorm:"column:status;type:TINYINT(1);NOT NULL"`
}

// 资源路径
type ResourcePath struct {
	gorm.Model
	ID           int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	Title        string `gorm:"column:title;type:VARCHAR(64);NOT NULL"`
	Classify     string `gorm:"column:classify;type:VARCHAR(32);"`
	Cover        string `gorm:"column:cover;type:VARCHAR(256);"`
	Url          string `gorm:"column:url;type:VARCHAR(256);"`
	Introduction string `gorm:"column:introduction;type:VARCHAR(1024);"`
	Type         string `gorm:"column:type;type:VARCHAR(32);NOT NULL"`
	Status       bool   `gorm:"column:status;type:TINYINT(1);NOT NULL"`
	Remark       string `gorm:"column:remark;type:TEXT;"`
	//CreateTime   time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 资源信息
type Resource struct {
	gorm.Model
	ID       int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	UserId   int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	Type     string `gorm:"column:type;type:VARCHAR(32);NOT NULL"`
	Path     string `gorm:"column:path;type:VARCHAR(256);NOT NULL"`
	Size     int32  `gorm:"column:size;type:INT;"`
	MimeType string `gorm:"column:mime_type;type:VARCHAR(256);"`
	Status   bool   `gorm:"column:status;type:TINYINT(1);NOT NULL"`
	//CreateTime time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 家庭信息
type Family struct {
	gorm.Model
	ID             int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	UserId         int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	BgCover        string `gorm:"column:bg_cover;type:VARCHAR(256);NOT NULL"`
	ManCover       string `gorm:"column:man_cover;type:VARCHAR(256);NOT NULL"`
	WomanCover     string `gorm:"column:woman_cover;type:VARCHAR(256);NOT NULL"`
	ManName        string `gorm:"column:man_name;type:VARCHAR(32);NOT NULL"`
	WomanName      string `gorm:"column:woman_name;type:VARCHAR(32);NOT NULL"`
	Timing         string `gorm:"column:timing;type:VARCHAR(32);NOT NULL"`
	CountdownTitle string `gorm:"column:countdown_title;type:VARCHAR(32);"`
	CountdownTime  string `gorm:"column:countdown_time;type:VARCHAR(32);"`
	Status         bool   `gorm:"column:status;type:TINYINT(1);NOT NULL"`
	FamilyInfo     string `gorm:"column:family_info;type:VARCHAR(1024);"`
	LikeCount      int32  `gorm:"column:like_count;type:INT;NOT NULL"`
	//CreateTime     time.Time `gorm:"column:create_time;type:DATETIME;"`
	//UpdateTime     time.Time `gorm:"column:update_time;type:DATETIME;"`
}

// 好友
type ImChatUserFriend struct {
	gorm.Model
	ID           int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	UserId       int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	FriendId     int32  `gorm:"column:friend_id;type:INT;NOT NULL"`
	FriendStatus int8   `gorm:"column:friend_status;type:TINYINT(2);NOT NULL"`
	Remark       string `gorm:"column:remark;type:VARCHAR(32);"`
	//CreateTime   time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 聊天群
type ImChatGroup struct {
	gorm.Model
	ID           int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	GroupName    string `gorm:"column:group_name;type:VARCHAR(32);NOT NULL"`
	MasterUserId int32  `gorm:"column:master_user_id;type:INT;NOT NULL"`
	Avatar       string `gorm:"column:avatar;type:VARCHAR(256);"`
	Introduction string `gorm:"column:introduction;type:VARCHAR(128);"`
	Notice       string `gorm:"column:notice;type:VARCHAR(1024);"`
	InType       bool   `gorm:"column:in_type;type:TINYINT(1);NOT NULL"`
	GroupType    int8   `gorm:"column:group_type;type:TINYINT(2);NOT NULL"`
	//CreateTime   time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 聊天群成员
type ImChatGroupUser struct {
	gorm.Model
	ID           int32  `gorm:"column:id;type:INT;AUTO_INCREMENT;NOT NULL"`
	GroupId      int32  `gorm:"column:group_id;type:INT;NOT NULL"`
	UserId       int32  `gorm:"column:user_id;type:INT;NOT NULL"`
	VerifyUserId int32  `gorm:"column:verify_user_id;type:INT;"`
	Remark       string `gorm:"column:remark;type:VARCHAR(1024);"`
	AdminFlag    bool   `gorm:"column:admin_flag;type:TINYINT(1);NOT NULL"`
	UserStatus   int8   `gorm:"column:user_status;type:TINYINT(2);NOT NULL"`
	//CreateTime   time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 单聊记录
type ImChatUserMessage struct {
	gorm.Model
	ID            int64  `gorm:"column:id;type:BIGINT;AUTO_INCREMENT;NOT NULL"`
	FromId        int32  `gorm:"column:from_id;type:INT;NOT NULL"`
	ToId          int32  `gorm:"column:to_id;type:INT;NOT NULL"`
	Content       string `gorm:"column:content;type:VARCHAR(1024);NOT NULL"`
	MessageStatus bool   `gorm:"column:message_status;type:TINYINT(1);NOT NULL"`
	//CreateTime    time.Time `gorm:"column:create_time;type:DATETIME;"`
}

// 群聊记录
type ImChatUserGroupMessage struct {
	gorm.Model
	ID      int64  `gorm:"column:id;type:BIGINT;AUTO_INCREMENT;NOT NULL"`
	GroupId int32  `gorm:"column:group_id;type:INT;NOT NULL"`
	FromId  int32  `gorm:"column:from_id;type:INT;NOT NULL"`
	ToId    int32  `gorm:"column:to_id;type:INT;"`
	Content string `gorm:"column:content;type:VARCHAR(1024);NOT NULL"`
	//CreateTime time.Time `gorm:"column:create_time;type:DATETIME;"`
}

func (to *WebInfo) Copy(from *WebInfo) {
	to.ID = from.ID
	to.WebName = from.WebName
	to.WebTitle = from.WebTitle
	to.Notices = from.Notices
	to.Footer = from.Footer
	to.BackgroundImage = from.BackgroundImage
	to.Avatar = from.Avatar
	to.RandomAvatar = from.RandomAvatar
	to.RandomName = from.RandomName
	to.RandomCover = from.RandomCover
	to.WaifuJson = from.WaifuJson
	to.Status = from.Status
}

func (to *ResourcePath) Copy(from *blog.ResourcePathVO) {
	to.ID = from.ID
	to.Title = from.Title
	to.Classify = from.Classify
	to.Cover = from.Cover
	to.Url = from.Url
	to.Introduction = from.Introduction
	to.Type = from.Type
	to.Status = from.Status
	to.Remark = from.Remark
	to.CreatedAt = from.CreatedAt
}
