# 🚀 Render Deployment Guide

## Quick Deploy (Recommended)

1. **Push to GitHub** (already done ✅)

2. **Go to Render Dashboard**
   - Visit: https://dashboard.render.com
   - Sign up/Login with GitHub

3. **Create New Blueprint**
   - Click "New +" → "Blueprint"
   - Connect your GitHub repo: `VinVorteX/NoBurn`
   - Render will auto-detect `render.yaml`

4. **Set Environment Variables**
   After services are created, add these secrets:
   
   **For API & Worker:**
   - `HUGGING_FACE_TOKEN`: Your HuggingFace token
   - `SMTP_USER`: Your Gmail address
   - `SMTP_PASSWORD`: Gmail App Password
   - `ALERT_EMAIL`: Email for alerts
   - `SLACK_WEBHOOK_URL`: (Optional) Slack webhook

5. **Run Database Migrations**
   ```bash
   # In Render Shell for API service
   ./migrate -path migrations -database $DB_URL up
   ```

6. **Done!** 🎉
   - API: `https://noburn-api.onrender.com`
   - Frontend: `https://noburn-frontend.onrender.com`

---

## Manual Deployment (Alternative)

### 1. PostgreSQL Database
- New → PostgreSQL
- Name: `noburn-postgres`
- Database: `noburn_db`
- User: `noburn_user`
- Plan: Free

### 2. Redis
- New → Redis
- Name: `noburn-redis`
- Plan: Free

### 3. API Service
- New → Web Service
- Connect GitHub repo
- Settings:
  - **Name**: `noburn-api`
  - **Environment**: Docker
  - **Dockerfile Path**: `./Dockerfile`
  - **Docker Command**: `./api`
  - **Plan**: Free

**Environment Variables:**
```
PORT=3000
ENV=production
DB_URL=<postgres-internal-url>
JWT_SECRET=<generate-random-64-char-string>
JWT_EXPIRES_IN=24h
REDIS_URL=<redis-internal-url>
HUGGING_FACE_TOKEN=<your-token>
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=<your-gmail>
SMTP_PASSWORD=<gmail-app-password>
ALERT_EMAIL=<alert-email>
SLACK_WEBHOOK_URL=<optional>
```

### 4. Worker Service
- New → Background Worker
- Same repo
- Settings:
  - **Name**: `noburn-worker`
  - **Environment**: Docker
  - **Dockerfile Path**: `./Dockerfile`
  - **Docker Command**: `./worker`
  - **Plan**: Free

**Environment Variables:** (Same as API)

### 5. Frontend
- New → Static Site
- Same repo
- Settings:
  - **Name**: `noburn-frontend`
  - **Build Command**: `cd frontend && npm install && npm run build`
  - **Publish Directory**: `frontend/build`

**Environment Variables:**
```
REACT_APP_API_URL=https://noburn-api.onrender.com
```

**Rewrite Rules:**
```
/*  /index.html  200
```

---

## Post-Deployment

### Run Migrations
```bash
# SSH into API service shell
./migrate -path migrations -database $DB_URL up
```

### Test Endpoints
```bash
# Health check
curl https://noburn-api.onrender.com/health

# Register
curl -X POST https://noburn-api.onrender.com/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"test123","name":"Admin","company_name":"Test Co"}'
```

---

## Important Notes

⚠️ **Free Tier Limitations:**
- Services sleep after 15 min inactivity
- 750 hours/month (shared across services)
- PostgreSQL: 1GB storage
- Redis: 25MB storage

💡 **Tips:**
- Use internal URLs for service-to-service communication
- Set `REDIS_URL` to internal Redis URL (faster)
- Set `DB_URL` to internal Postgres URL
- Frontend uses public API URL

🔒 **Security:**
- Never commit `.env` file
- Use Render's environment variables
- Regenerate JWT_SECRET for production
- Use Gmail App Passwords (not regular password)

---

## Troubleshooting

**Services not starting?**
- Check logs in Render dashboard
- Verify all environment variables are set
- Ensure Docker build succeeds

**Database connection failed?**
- Use internal database URL
- Check if migrations ran successfully

**Worker not processing jobs?**
- Verify REDIS_URL is correct
- Check worker logs for errors

**Frontend can't reach API?**
- Update `REACT_APP_API_URL` to API public URL
- Check CORS settings in API

---

## Cost Estimate

**Free Tier (Current Setup):**
- PostgreSQL: Free (1GB)
- Redis: Free (25MB)
- API: Free (750 hrs)
- Worker: Free (750 hrs)
- Frontend: Free (100GB bandwidth)

**Total: $0/month** 🎉

**Paid Tier (For Production):**
- PostgreSQL: $7/month (256MB RAM)
- Redis: $10/month (100MB)
- API: $7/month (always on)
- Worker: $7/month (always on)
- Frontend: Free

**Total: ~$31/month**
