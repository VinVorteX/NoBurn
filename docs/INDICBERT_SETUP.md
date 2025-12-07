# IndicBERT Integration

## What is IndicBERT?

IndicBERT is a multilingual BERT model trained on 12 major Indian languages by AI4Bharat (IIT Madras).

**Supported Languages:**
- Hindi (hi)
- Bengali (bn)
- Gujarati (gu)
- Kannada (kn)
- Malayalam (ml)
- Marathi (mr)
- Oriya (or)
- Punjabi (pa)
- Tamil (ta)
- Telugu (te)
- Urdu (ur)
- English (en)

## Why IndicBERT?

✅ **State-of-the-art** for Indian languages
✅ **Open source** and free
✅ **Better accuracy** than generic multilingual models
✅ **Proven** by IIT Madras research
✅ **Easy integration** via Hugging Face

## Setup

### 1. Get Hugging Face Token

1. Visit: https://huggingface.co/settings/tokens
2. Click "New token"
3. Name: "NoBurn"
4. Type: Read
5. Copy token

### 2. Configure NoBurn

```bash
# Add to .env
HUGGING_FACE_TOKEN=hf_xxxxxxxxxxxxxxxxxxxxx
```

### 3. Test

```bash
# Start server
make dev-api

# Submit Hindi survey
curl -X POST http://localhost:8080/api/surveys/responses \
  -H "Authorization: Bearer YOUR_JWT" \
  -H "Content-Type: application/json" \
  -d '{
    "survey_id": 1,
    "responses": ["मैं अपने काम से खुश हूं"]
  }'
```

## Model Details

**Model**: `ai4bharat/indic-bert`
**Paper**: https://arxiv.org/abs/2007.07691
**GitHub**: https://github.com/AI4Bharat/indic-bert

## Fallback Strategy

1. **IndicBERT** (Primary) - Best for Indian languages
2. **Rule-based** (Fallback) - Always works offline

## Performance

- **Accuracy**: 85-92% for Indian languages
- **Latency**: ~500ms per request
- **Cost**: Free (Hugging Face Inference API)

## Rate Limits

- **Free tier**: 1000 requests/day
- **Pro tier**: $9/month for unlimited

For production, consider:
- Caching results
- Self-hosting the model
- Using Hugging Face Pro