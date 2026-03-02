package storage

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"war-report/internal/config"
	"war-report/internal/templates"
)

func ExtractPreviousReport() string {
	filename := "index.html"
	content, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	htmlStr := string(content)
	cardStart := "<!-- MR_REPORT_CARD_START -->"
	cardEnd := "<!-- MR_REPORT_CARD_END -->"

	startIdx := strings.Index(htmlStr, cardStart)
	if startIdx == -1 {
		return ""
	}

	endIdx := strings.Index(htmlStr[startIdx:], cardEnd)
	if endIdx == -1 {
		return ""
	}
	endIdx += startIdx

	cardBlock := htmlStr[startIdx : endIdx+len(cardEnd)]

	contentStartTag := `<div class="content">`
	contentEndTag := `</div>`

	contentStart := strings.Index(cardBlock, contentStartTag)
	if contentStart == -1 {
		return ""
	}

	contentEnd := strings.LastIndex(cardBlock, contentEndTag)
	if contentEnd == -1 || contentEnd < contentStart {
		return ""
	}

	previousContent := strings.TrimSpace(cardBlock[contentStart+len(contentStartTag) : contentEnd])

	timestampRegex := regexp.MustCompile(`<div class="timestamp">REPORTE GENERADO: ([^<]+)</div>`)
	matches := timestampRegex.FindStringSubmatch(cardBlock)
	timestamp := "Fecha desconocida"
	if len(matches) > 1 {
		timestamp = matches[1]
	}

	return fmt.Sprintf("[Timestamp: %s]\n%s", timestamp, previousContent)
}

func UpdateHTML(newReportContent string) error {
	filename := "index.html"
	nowDisplay := time.Now().UTC().Format("02-Jan-2006 15:04 UTC")
	currentYear := fmt.Sprintf("%d", time.Now().Year())

	reportsStart := "<!-- MR_REPORTS_START -->"
	reportsEnd := "<!-- MR_REPORTS_END -->"
	cardStart := "<!-- MR_REPORT_CARD_START -->"
	cardEnd := "<!-- MR_REPORT_CARD_END -->"

	existingReportsHTML := ""
	content, err := os.ReadFile(filename)
	if err == nil {
		htmlStr := string(content)
		s := strings.Index(htmlStr, reportsStart)
		e := strings.Index(htmlStr, reportsEnd)
		if s != -1 && e != -1 && e > s {
			existingReportsHTML = htmlStr[s+len(reportsStart) : e]
		}
	}

	newBlock := fmt.Sprintf(`
            %s
            <div class="report-card">
                <div class="timestamp">REPORTE GENERADO: %s</div>
                <div class="content">
                    %s
                </div>
            </div>
            %s
    `, cardStart, nowDisplay, newReportContent, cardEnd)

	cardRe := regexp.MustCompile(fmt.Sprintf(`(?s)%s.*?%s`, regexp.QuoteMeta(cardStart), regexp.QuoteMeta(cardEnd)))
	oldCards := cardRe.FindAllString(existingReportsHTML, -1)

	var allCards []string
	allCards = append(allCards, newBlock)
	for i := 0; i < len(oldCards) && len(allCards) < config.MaxReports; i++ {
		allCards = append(allCards, strings.TrimSpace(oldCards[i]))
	}

	finalReportsSection := strings.Join(allCards, "\n\n")

	baseHTMLFormatted := strings.Replace(templates.BaseHTML, "{current_year}", currentYear, 1)

	placeholder := "            <!-- MR_REPORTS_START -->\n            <!-- MR_REPORTS_END -->"
	replacement := fmt.Sprintf("            <!-- MR_REPORTS_START -->\n%s\n            <!-- MR_REPORTS_END -->", finalReportsSection)

	finalHTML := strings.Replace(baseHTMLFormatted, placeholder, replacement, 1)

	return os.WriteFile(filename, []byte(finalHTML), 0644)
}
