package request

type PostCreateRequest struct {
	Title           string `json:"title"`
	MarkdownContent string `json:"markdown_content"`
}
