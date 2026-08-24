package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/jaimehisao/locksmith/internal/copy"
)

func noRecordFoundOn(date time.Time) string {
	return strings.Replace(copy.NoRecordFoundOnDate, "DATE", date.Format("2006-01-02"), 1)
}

func confidenceText(confidence int) string {
	return fmt.Sprintf("%d %s", confidence, copy.ReportsIndicated)
}
