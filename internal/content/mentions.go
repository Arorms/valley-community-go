package content

import (
	"regexp"
	"strings"
)

var (
	mentionFencedCodePattern = regexp.MustCompile("(?s)```.*?```")
	mentionInlineCodePattern = regexp.MustCompile("`[^`]*`")
	mentionPattern           = regexp.MustCompile(`(?:^|[^\w@])@([A-Za-z0-9_.-]{1,64})`)
)

// Mentions 解析 Markdown 中的 @username，跳过代码块/行内代码，返回去重后的小写用户名。
func Mentions(source string) []string {
	text := mentionFencedCodePattern.ReplaceAllString(source, " ")
	text = mentionInlineCodePattern.ReplaceAllString(text, " ")

	matches := mentionPattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]struct{}, len(matches))
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		name := strings.ToLower(match[1])
		if _, ok := seen[name]; ok || name == "" {
			continue
		}
		seen[name] = struct{}{}
		result = append(result, name)
	}
	return result
}
