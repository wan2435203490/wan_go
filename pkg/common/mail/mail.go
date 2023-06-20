package mail

import (
	json2 "encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"gorm.io/gorm/utils"
	"net/smtp"
	"sync/atomic"
	"wan_go/pkg/common/config"
	"wan_go/pkg/common/constant/blog_const"
	"wan_go/pkg/common/log"
	blogVO "wan_go/pkg/vo/blog"
)

var (
	/**
	 * 1. 来源人名
	 * 2. 来源内容
	 */
	OriginalText = "<hr style=\"border: 1px dashed #ef859d2e;margin: 20px 0\">\n" +
		"            <div>\n" +
		"                <div style=\"font-size: 18px;font-weight: bold;color: #C5343E\">\n" +
		"                    %s\n" +
		"                </div>\n" +
		"                <div style=\"margin-top: 6px;font-size: 16px;color: #000000\">\n" +
		"                    <p>\n" +
		"                        %s\n" +
		"                    </p>\n" +
		"                </div>\n" +
		"            </div>"

	/**
	 * 1. 网站名称
	 * 2. 邮件类型
	 * 3. 发件人
	 * 4. 发件内容
	 * 5. originalText
	 * 6. 网站名称
	 */
	MailText = "<div style=\"font-family: serif;line-height: 22px;padding: 30px\">\n" +
		"    <div style=\"display: flex;justify-content: center;width: 100%%;max-width: 900px;background-image: url('" + blog_const.DOWNLOAD_URL + "webBackgroundImage/Sara11667042705239112');background-size: cover;border-radius: 10px\"></div>\n" +
		"    <div style=\"margin-top: 20px;display: flex;flex-direction: column;align-items: center\">\n" +
		"        <div style=\"margin: 10px auto 20px;text-align: center\">\n" +
		"            <div style=\"line-height: 32px;font-size: 26px;font-weight: bold;color: #000000\">\n" +
		"                嘿！你收到一条来自 %s 的一条新消息。\n" +
		"            </div>\n" +
		"            <div style=\"font-size: 16px;font-weight: bold;color: rgba(0, 0, 0, 0.19);margin-top: 21px\">\n" +
		"                %s\n" +
		"            </div>\n" +
		"        </div>\n" +
		"        <div style=\"min-width: 250px;max-width: 800px;min-height: 128px;background: #F7F7F7;border-radius: 10px;padding: 32px\">\n" +
		"            <div>\n" +
		"                <div style=\"font-size: 18px;font-weight: bold;color: #C5343E\">\n" +
		"                    %s\n" +
		"                </div>\n" +
		"                <div style=\"margin-top: 6px;font-size: 16px;color: #000000\">\n" +
		"                    <p>\n" +
		"                        %s\n" +
		"                    </p>\n" +
		"                </div>\n" +
		"            </div>\n" +
		"            %s\n" +
		"            <a style=\"width: 150px;height: 38px;background: #ef859d38;border-radius: 32px;display: flex;align-items: center;justify-content: center;text-decoration: none;margin: 40px auto 0\"\n" +
		"               href=\"https://2fun.top\" target=\"_blank\">\n" +
		"                <span style=\"color: #DB214B\">有朋自远方来</span>\n" +
		"            </a>\n" +
		"        </div>\n" +
		"        <div style=\"margin-top: 20px;font-size: 12px;color: #00000045\">\n" +
		"            此邮件由 %s 自动发出，直接回复无效（一天最多发送 " + utils.ToString(blog_const.COMMENT_IM_MAIL_COUNT) + " 条通知邮件），退订请联系站长。\n" +
		"        </div>\n" +
		"    </div>\n" +
		"</div>"

	/**
	 * 发件人
	 */
	ReplyMail   = "你之前的评论收到来自 %s 的回复"
	CommentMail = "你的文章 %s 收到来自 %s 的评论"
	MessageMail = "你收到来自 %s 的留言"
	LoveMail    = "你收到来自 %s 的祝福"
	ImMail      = "你收到来自 %s 的消息"
)

func SendMail(to []string, subject, text string) {
	toBytes, err := json2.Marshal(to)
	log.Info("SendMail", "发送邮件===================")
	log.Info("SendMail", toBytes)
	log.Info("SendMail", subject, text)
	// 实例化返回一个结构体
	e := email.NewEmail()

	// From：谁发来的
	e.From = config.Config.Mail.Username

	// To：发给谁的
	e.To = to

	// 抄送,这个地方抄送的意思是：这个邮件在发送后还可以抄送给谁
	//e.Bcc = []string{"1111@qq.com"}
	//e.Cc = []string{"2222@qq.com"}

	// 主题，标题
	e.Subject = subject

	// 普通文本内容，支持html
	e.Text = []byte(text)
	//e.HTML = []byte("<h1>html 邮件测试 </h1>")

	auth := smtp.PlainAuth("", config.Config.Mail.Username, config.Config.Mail.Password, config.Config.Mail.Host)
	err = e.Send(config.Config.Mail.Host+":"+config.Config.Mail.Port, auth)
	if err == nil {
		log.Info("SendMail", "发送成功==================")
	} else {
		log.Info("SendMail", "发送失败==================%s", err.Error())
	}
}

func SendSimpleCommentMail(vo *blogVO.CommentVO, mails []string, fromName, toName, webName string,
	sendCount *atomic.Int32, sendSuccess func()) {
	SendCommentMail(vo, mails, "", fromName, toName, webName, "", sendCount, sendSuccess)
}

func SendCommentMail(vo *blogVO.CommentVO, mails []string, articleTitle, fromName, toName, webName, mailContent string,
	sendCount *atomic.Int32, sendSuccess func()) {

	sourceName := ""
	if vo.Type == blog_const.COMMENT_TYPE_ARTICLE.Code {
		sourceName = articleTitle
	}

	commentMail := getCommentMail(
		vo.Type,
		sourceName,
		fromName,
		vo.CommentContent,
		toName,
		webName,
		mailContent,
	)

	if sendCount.Load() < blog_const.COMMENT_IM_MAIL_COUNT {

		SendMail(mails, "您有一封来自"+webName+"的回执！", commentMail)

		if sendCount == nil {
			sendSuccess()
		} else {
			//直接set 1
			sendCount.Swap(1)
		}
	}

}

/**
 * source：0留言 其他是文章标题
 * fromName：评论人
 * toName：被评论人
 */
func getCommentMail(commentType, source, fromName, fromContent, toName, webName, mailContent string) string {
	var mailType string

	if len(toName) > 0 {
		mailType = fmt.Sprintf(ReplyMail, fromName)
	} else {
		switch commentType {
		case blog_const.COMMENT_TYPE_MESSAGE.Code:
			mailType = fmt.Sprintf(MessageMail, fromName)
		case blog_const.COMMENT_TYPE_ARTICLE.Code:
			mailType = fmt.Sprintf(CommentMail, source, fromName)
		case blog_const.COMMENT_TYPE_LOVE.Code:
			mailType = fmt.Sprintf(LoveMail, fromName)
		default:
			break
		}
	}

	return fmt.Sprintf(MailText,
		webName,
		mailType,
		fromName,
		fromContent,
		mailContent,
		webName)
}
