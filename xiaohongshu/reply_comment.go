package xiaohongshu

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

// ReplyCommentAction 表示回复评论动作
type ReplyCommentAction struct {
	page *rod.Page
}

// NewReplyCommentAction 创建回复评论动作
func NewReplyCommentAction(page *rod.Page) *ReplyCommentAction {
	return &ReplyCommentAction{page: page}
}

// ReplyToComment 回复指定评论
func (r *ReplyCommentAction) ReplyToComment(ctx context.Context, feedID, xsecToken, commentID, content string) error {
	page := r.page.Context(ctx).Timeout(60 * time.Second)

	// 构建详情页 URL
	url := makeFeedDetailURL(feedID, xsecToken)
	logrus.Infof("打开笔记详情页准备回复评论: %s", url)

	// 导航到详情页
	page.MustNavigate(url)
	page.MustWaitDOMStable()
	time.Sleep(1 * time.Second)

	// 定位目标评论并点击回复按钮
	// 评论区域的选择器: div.comment-item[data-id="commentID"]
	commentSelector := fmt.Sprintf(`div.comment-item[data-v-comment-id="%s"]`, commentID)

	// 尝试找到目标评论
	comment, err := page.Element(commentSelector)
	if err != nil {
		logrus.Warnf("未找到评论元素，尝试备用选择器: %v", err)
		// 备用方案：遍历评论找到匹配的
		comment, err = r.findCommentByID(page, commentID)
		if err != nil {
			return fmt.Errorf("未找到评论 %s: %w", commentID, err)
		}
	}

	// 点击评论的回复按钮
	replyBtn, err := comment.Element("span.reply-btn")
	if err != nil {
		// 尝试其他可能的回复按钮选择器
		replyBtn, err = comment.Element("span.reply")
		if err != nil {
			return fmt.Errorf("未找到回复按钮: %w", err)
		}
	}

	replyBtn.MustClick()
	time.Sleep(500 * time.Millisecond)

	// 输入回复内容
	inputElem := page.MustElement("div.input-box div.content-edit p.content-input")
	inputElem.MustInput(content)
	time.Sleep(500 * time.Millisecond)

	// 点击发送按钮
	submitButton := page.MustElement("div.bottom button.submit")
	submitButton.MustClick()
	time.Sleep(1 * time.Second)

	logrus.Infof("成功回复评论 %s", commentID)
	return nil
}

// findCommentByID 通过遍历评论列表找到目标评论
func (r *ReplyCommentAction) findCommentByID(page *rod.Page, commentID string) (*rod.Element, error) {
	// 获取所有评论项
	comments, err := page.Elements("div.comment-item, div.parent-comment")
	if err != nil {
		return nil, fmt.Errorf("获取评论列表失败: %w", err)
	}

	for _, comment := range comments {
		// 检查评论的 data 属性或内部元素
		id, err := comment.Attribute("data-v-comment-id")
		if err == nil && id != nil && *id == commentID {
			return comment, nil
		}

		// 检查其他可能的 ID 属性
		id, err = comment.Attribute("data-comment-id")
		if err == nil && id != nil && *id == commentID {
			return comment, nil
		}
	}

	return nil, fmt.Errorf("未找到 ID 为 %s 的评论", commentID)
}
