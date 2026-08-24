package locksmithwatch

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/jaimehisao/locksmith/internal/copy"
)

var (
	listPageTemplate = template.Must(template.New("list").Parse(`<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>LocksmithWatch SF</title></head>
<body>
{{if .Records}}
<ul>
{{range .Records}}<li><a href="/locksmiths/{{.ID}}">{{.Name}}</a></li>
{{end}}
</ul>
{{else}}
<p>{{.ReportsIndicated}}</p>
<p>{{.AppearRelated}}</p>
<p>{{.NoRecordFound}}</p>
{{end}}
</body>
</html>`))

	detailPageTemplate = template.Must(template.New("detail").Parse(`<!doctype html>
<html lang="en">
<head><meta charset="utf-8"><title>{{.Name}}</title></head>
<body>
<p>{{.Name}}</p>
<p>{{.Phone}}</p>
<p>{{.License}}</p>
<p><a href="{{.SourceURL}}">{{.SourceURL}}</a></p>
<p>{{.Confidence}}</p>
</body>
</html>`))
)

type listPageData struct {
	Records           []Record
	ReportsIndicated  string
	AppearRelated     string
	NoRecordFound     string
}

// EmptyStateNoRecordFound returns the allow-listed no-record phrase with DATE replaced.
func EmptyStateNoRecordFound(date time.Time) string {
	return strings.ReplaceAll(copy.NoRecordFoundOnDate, "DATE", date.Format("2006-01-02"))
}

// RenderListPage renders the published list page or empty state copy.
func RenderListPage(records []Record, now time.Time) (string, error) {
	data := listPageData{
		Records:          records,
		ReportsIndicated: copy.ReportsIndicated,
		AppearRelated:    copy.AppearRelated,
		NoRecordFound:    EmptyStateNoRecordFound(now),
	}
	var buf bytes.Buffer
	if err := listPageTemplate.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render list page: %w", err)
	}
	return buf.String(), nil
}

// RenderDetailPage renders a published record detail page.
func RenderDetailPage(record Record) (string, error) {
	var buf bytes.Buffer
	if err := detailPageTemplate.Execute(&buf, record); err != nil {
		return "", fmt.Errorf("render detail page: %w", err)
	}
	return buf.String(), nil
}
