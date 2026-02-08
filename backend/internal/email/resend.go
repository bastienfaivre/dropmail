package email

import (
	"fmt"
	"strings"
	"time"

	"backend/internal/model"

	"github.com/resend/resend-go/v3"
)

// Client wraps the Resend SDK to send summary emails.
type Client struct {
	client *resend.Client
	from   string
	to     string
}

// NewClient creates a new Resend email client.
func NewClient(apiKey, from, to string) *Client {
	return &Client{
		client: resend.NewClient(apiKey),
		from:   from,
		to:     to,
	}
}

// SendDailySummary sends an email summarizing recent submissions.
func (c *Client) SendDailySummary(submissions []model.Submission) error {
	count := len(submissions)
	subject := fmt.Sprintf("Dropmail Daily Summary — %d submission(s)", count)

	html := buildSummaryHTML(submissions)

	params := &resend.SendEmailRequest{
		From:    c.from,
		To:      []string{c.to},
		Subject: subject,
		Html:    html,
	}

	_, err := c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send summary email: %w", err)
	}

	return nil
}

func buildSummaryHTML(submissions []model.Submission) string {
	count := len(submissions)
	var b strings.Builder

	b.WriteString("<h2>Daily Submission Summary</h2>")
	b.WriteString(fmt.Sprintf("<p><strong>%d</strong> new submission(s) in the last 24 hours.</p>", count))

	if count == 0 {
		return b.String()
	}

	b.WriteString("<table border=\"1\" cellpadding=\"8\" cellspacing=\"0\" style=\"border-collapse:collapse\">")
	b.WriteString("<tr><th>Email</th><th>Source</th><th>Date</th></tr>")
	for _, s := range submissions {
		b.WriteString(fmt.Sprintf(
			"<tr><td>%s</td><td>%s</td><td>%s</td></tr>",
			s.Email,
			s.Source,
			s.CreatedAt.Format(time.RFC822),
		))
	}
	b.WriteString("</table>")

	return b.String()
}
