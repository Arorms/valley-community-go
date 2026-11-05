package migrations

import "embed"

// FS 嵌入所有迁移 SQL，供 cmd/migrate 使用。
//
//go:embed *.sql
var FS embed.FS
