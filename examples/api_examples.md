# NoBurn API Examples

## Multi-Language AI Suggestions

### Get Retention Suggestions in English
```bash
curl -X GET "http://localhost:8080/api/retention-suggestions/5?language=en" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "5",
    "language": "en",
    "suggestions": [
      {
        "suggestion": "Schedule 1-on-1 meeting",
        "priority": "high",
        "category": "Low sentiment scores"
      },
      {
        "suggestion": "Discuss career goals",
        "priority": "high",
        "category": "Low sentiment scores"
      }
    ]
  }
}
```

### Get Retention Suggestions in Hindi
```bash
curl -X GET "http://localhost:8080/api/retention-suggestions/5?language=hi" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Retention Suggestions in Tamil
```bash
curl -X GET "http://localhost:8080/api/retention-suggestions/5?language=ta" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Supported Languages
- `en` - English
- `hi` - Hindi (हिंदी)
- `ta` - Tamil (தமிழ்)