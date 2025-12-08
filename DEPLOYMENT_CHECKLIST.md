# 🚀 Pre-Deployment Checklist

## Before Deploying to Production

### 1. Environment Variables
- [ ] Generate new `JWT_SECRET` (64+ characters)
  ```bash
  openssl rand -hex 32
  ```
- [ ] Set up Gmail App Password (not regular password)
  - Go to: https://myaccount.google.com/apppasswords
- [ ] Get HuggingFace token (free)
  - Go to: https://huggingface.co/settings/tokens
- [ ] Update `ALERT_EMAIL` to your HR email

### 2. Security
- [ ] Change default database credentials
- [ ] Update CORS origins in `internal/server/router.go` (line 32)
  ```go
  // Change from "*" to your frontend domain
  w.Header().Set("Access-Control-Allow-Origin", "https://your-frontend.com")
  ```
- [ ] Never commit `.env` file to Git

### 3. Database
- [ ] Run migrations after first deployment
  ```bash
  ./migrate -path migrations -database $DB_URL up
  ```

### 4. Testing
- [ ] Test Docker build locally
  ```bash
  docker compose up -d
  ```
- [ ] Test health endpoint
  ```bash
  curl http://localhost:3000/health
  ```
- [ ] Test registration
  ```bash
  curl -X POST http://localhost:3000/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"test123","name":"Test","company_name":"Test Co"}'
  ```

### 5. Render Deployment
- [ ] Push code to GitHub
- [ ] Connect Render to your repo
- [ ] Render will auto-detect `render.yaml`
- [ ] Add environment secrets in Render dashboard:
  - `HUGGING_FACE_TOKEN`
  - `SMTP_USER`
  - `SMTP_PASSWORD`
  - `ALERT_EMAIL`
  - `SLACK_WEBHOOK_URL` (optional)

### 6. Post-Deployment
- [ ] Check all service logs in Render
- [ ] Verify database migrations ran
- [ ] Test user registration
- [ ] Test survey creation
- [ ] Test email notifications

## Quick Deploy Commands

### Local Docker
```bash
# Start all services
docker compose up -d

# Check logs
docker compose logs -f

# Stop services
docker compose down
```

### Production URLs (after Render deploy)
- Frontend: `https://noburn-frontend.onrender.com`
- API: `https://noburn-api.onrender.com`
- Health: `https://noburn-api.onrender.com/health`

## Troubleshooting

**Services not starting?**
- Check Render logs
- Verify all env vars are set
- Ensure migrations ran

**Email not sending?**
- Use Gmail App Password (not regular password)
- Enable "Less secure app access" if needed

**Frontend can't reach API?**
- Update `VITE_API_URL` in Render frontend settings
- Check CORS configuration

## Support
- GitHub Issues: https://github.com/VinVorteX/NoBurn/issues
- Documentation: See README.md
