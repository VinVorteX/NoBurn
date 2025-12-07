package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"
)

type NotificationService struct {
	slackWebhookURL string
	smtpHost        string
	smtpPort        int
	smtpUser        string
	smtpPassword    string
	client          *http.Client
}

func NewNotificationService(slackWebhook, smtpHost string, smtpPort int, smtpUser, smtpPassword string) *NotificationService {
	return &NotificationService{
		slackWebhookURL: slackWebhook,
		smtpHost:        smtpHost,
		smtpPort:        smtpPort,
		smtpUser:        smtpUser,
		smtpPassword:    smtpPassword,
		client:          &http.Client{Timeout: 5 * time.Second},
	}
}

// Slack notification
func (s *NotificationService) SendSlackAlert(userName, message string, riskScore float64) error {
	if s.slackWebhookURL == "" || s.slackWebhookURL == "https://hooks.slack.com/services/YOUR/WEBHOOK/URL" {
		// Silently skip if not configured
		return nil
	}

	color := "warning"
	if riskScore > 0.8 {
		color = "danger"
	}

	payload := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"color":      color,
				"title":      "🚨 High Churn Risk Alert",
				"text":       message,
				"fields": []map[string]interface{}{
					{
						"title": "Employee",
						"value": userName,
						"short": true,
					},
					{
						"title": "Risk Score",
						"value": fmt.Sprintf("%.0f%%", riskScore*100),
						"short": true,
					},
				},
				"footer": "NoBurn HR Analytics",
				"ts":     time.Now().Unix(),
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := s.client.Post(s.slackWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API error: %d", resp.StatusCode)
	}

	return nil
}

// Email notification using SMTP
func (s *NotificationService) SendEmailAlert(toEmail, userName, message string, riskScore float64) error {
	if s.smtpUser == "" || s.smtpPassword == "" {
		// Silently skip if not configured
		return nil
	}

	// Setup SMTP authentication
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	// Build email
	subject := fmt.Sprintf("⚠️ Churn Risk Alert: %s", userName)
	body := s.buildEmailHTML(userName, message, riskScore)
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.smtpUser, toEmail, subject, body))

	// Send email
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.smtpUser, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("SMTP error: %v", err)
	}

	return nil
}

func (s *NotificationService) buildEmailHTML(userName, message string, riskScore float64) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #ff6b6b; color: white; padding: 20px; border-radius: 5px; }
        .content { background: #f9f9f9; padding: 20px; margin-top: 20px; border-radius: 5px; }
        .risk-score { font-size: 24px; font-weight: bold; color: #ff6b6b; }
        .footer { margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🚨 High Churn Risk Alert</h2>
        </div>
        <div class="content">
            <p><strong>Employee:</strong> %s</p>
            <p><strong>Risk Score:</strong> <span class="risk-score">%.0f%%</span></p>
            <p><strong>Alert:</strong> %s</p>
            <hr>
            <p><strong>Recommended Actions:</strong></p>
            <ul>
                <li>Schedule immediate 1-on-1 meeting</li>
                <li>Review compensation and benefits</li>
                <li>Discuss career development opportunities</li>
                <li>Address any workplace concerns</li>
            </ul>
        </div>
        <div class="footer">
            <p>This is an automated alert from NoBurn HR Analytics</p>
        </div>
    </div>
</body>
</html>
`, userName, riskScore*100, message)
}

// Webhook for custom integrations
func (s *NotificationService) SendWebhook(webhookURL string, data map[string]interface{}) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL not provided")
	}

	jsonData, _ := json.Marshal(data)
	resp, err := s.client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook error: %d", resp.StatusCode)
	}

	return nil
}

// Send survey invitation email
func (s *NotificationService) SendSurveyInvitation(toEmail, subject, body string) error {
	if s.smtpUser == "" || s.smtpPassword == "" {
		return fmt.Errorf("SMTP not configured")
	}

	// Setup SMTP authentication
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	// Build email
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.smtpUser, toEmail, subject, body))

	// Send email
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	err := smtp.SendMail(addr, auth, s.smtpUser, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("SMTP error: %v", err)
	}

	return nil
}