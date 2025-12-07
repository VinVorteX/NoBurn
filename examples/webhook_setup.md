# Webhook Integration Setup

## Slack Integration

### 1. Create Slack App
1. Go to https://api.slack.com/apps
2. Click "Create New App" → "From scratch"
3. Name: "NoBurn HR Analytics"
4. Select your workspace

### 2. Enable Incoming Webhooks
1. Go to "Incoming Webhooks"
2. Toggle "Activate Incoming Webhooks" to ON
3. Click "Add New Webhook to Workspace"
4. Select channel (e.g., #hr-alerts)
5. Copy webhook URL

### 3. Configure NoBurn
```bash
# Add to .env
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXX
```

### 4. Test Slack Alert
```bash
curl -X POST http://localhost:8080/webhooks/slack \
  -H "Content-Type: application/json" \
  -d '{
    "type": "url_verification",
    "challenge": "test123"
  }'
```

### 5. Slack Event Subscriptions (Optional)
For receiving messages from Slack:
1. Go to "Event Subscriptions"
2. Enable Events
3. Request URL: `https://your-domain.com/webhooks/slack`
4. Subscribe to: `message.channels`

---

## Email Integration (SendGrid)

### 1. Create SendGrid Account
1. Sign up at https://sendgrid.com
2. Verify your sender email
3. Create API Key

### 2. Configure NoBurn
```bash
# Add to .env
SENDGRID_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
ALERT_EMAIL=hr@yourcompany.com
```

### 3. Test Email Alert
The system automatically sends emails when:
- Employee churn risk > 70%
- Survey sentiment is very negative
- Employee becomes inactive

### 4. Inbound Email Parsing (Optional)
For receiving survey responses via email:
1. Go to SendGrid → Settings → Inbound Parse
2. Add hostname: `surveys.yourcompany.com`
3. Destination URL: `https://your-domain.com/webhooks/email`

---

## Custom Webhook Integration

### Survey Submission Webhook
```bash
curl -X POST http://localhost:8080/webhooks/survey \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 5,
    "responses": [
      "I am satisfied with my work",
      "Team collaboration is good"
    ],
    "source": "external_form"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "accepted",
    "message": "Survey response queued for processing"
  }
}
```

---

## Notification Examples

### Slack Alert Format
```
🚨 High Churn Risk Alert

Employee: John Doe
Risk Score: 85%

Employee has high churn risk: 0.85

NoBurn HR Analytics
```

### Email Alert Format
- Subject: ⚠️ Churn Risk Alert: John Doe
- HTML formatted with risk score and recommended actions
- Includes actionable steps for HR

---

## Webhook Security

### Verify Slack Requests
```go
// Verify Slack signature
slackSignature := r.Header.Get("X-Slack-Signature")
slackTimestamp := r.Header.Get("X-Slack-Request-Timestamp")
// Implement signature verification
```

### Verify SendGrid Webhooks
```go
// Verify SendGrid signature
sendgridSignature := r.Header.Get("X-Twilio-Email-Event-Webhook-Signature")
// Implement signature verification
```

---

## Monitoring

Check webhook delivery status:
```bash
# View worker logs
docker-compose logs -f worker

# Check notification queue
redis-cli LLEN asynq:default:pending
```