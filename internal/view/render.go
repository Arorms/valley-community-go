package view

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model"
	"gitlab.arorms.cn/arorms/applications/valley-community-go/internal/model/viewmodel"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"safeHTML": func(value string) template.HTML {
			return template.HTML(value)
		},
		"fmtTime": func(value time.Time) string {
			return value.Local().Format("2006-01-02 15:04")
		},
		"relativeTime": relativeTime,
		"initials":     viewmodel.Initials,
		"add":          func(a, b int) int { return a + b },
		"sub":          func(a, b int) int { return a - b },
		"notificationText": func(kind model.NotificationType) string {
			switch kind {
			case model.NotificationTypeComment:
				return "评论了你的帖子"
			case model.NotificationTypeReply:
				return "回复了你的评论"
			case model.NotificationTypeMention:
				return "在评论中提到了你"
			case model.NotificationTypePostLike:
				return "赞了你的帖子"
			case model.NotificationTypeCommentLike:
				return "赞了你的评论"
			default:
				return "有新动态"
			}
		},
		"dict": func(values ...any) map[string]any {
			result := make(map[string]any, len(values)/2)
			for i := 0; i+1 < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					continue
				}
				result[key] = values[i+1]
			}
			return result
		},
	}

	templates, err := template.New("root").Funcs(funcs).ParseFS(
		TemplatesFS(),
		"partials/*.html",
		"pages/*.html",
	)
	if err != nil {
		return nil, err
	}
	return &Renderer{templates: templates}, nil
}

func (r *Renderer) Render(w io.Writer, page string, data any) error {
	return r.templates.ExecuteTemplate(w, page, data)
}

func relativeTime(value time.Time) string {
	delta := time.Since(value)
	switch {
	case delta < time.Minute:
		return "刚刚"
	case delta < time.Hour:
		return fmt.Sprintf("%d 分钟前", int(delta.Minutes()))
	case delta < 24*time.Hour:
		return fmt.Sprintf("%d 小时前", int(delta.Hours()))
	case delta < 7*24*time.Hour:
		return fmt.Sprintf("%d 天前", int(delta.Hours()/24))
	default:
		return value.Local().Format("2006-01-02")
	}
}
