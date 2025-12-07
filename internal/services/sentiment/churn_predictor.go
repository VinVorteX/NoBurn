package sentiment

import (
	"math"
)

type ChurnFeatures struct {
	AvgSentiment     float64
	ResponseRate     float64
	DaysInactive     int
	NegativeResponses int
	TotalResponses   int
	LastLoginDays    int
}

type ChurnPredictor struct{}

func NewChurnPredictor() *ChurnPredictor {
	return &ChurnPredictor{}
}

// Simple ML-like churn prediction using weighted features
func (cp *ChurnPredictor) PredictChurnRisk(features ChurnFeatures) float64 {
	// Weights based on HR research
	weights := map[string]float64{
		"sentiment":    0.50,  // Increased weight for sentiment
		"response":     0.15,
		"activity":     0.15,
		"engagement":   0.20,
	}

	// Normalize features to 0-1 scale
	sentimentScore := (1 - features.AvgSentiment) / 2 // Convert -1,1 to 0,1
	responseScore := 1 - features.ResponseRate
	activityScore := math.Min(float64(features.DaysInactive)/30, 1.0)
	
	engagementScore := 0.0
	if features.TotalResponses > 0 {
		engagementScore = float64(features.NegativeResponses) / float64(features.TotalResponses)
	}

	// Calculate weighted risk score
	riskScore := (sentimentScore * weights["sentiment"]) +
		(responseScore * weights["response"]) +
		(activityScore * weights["activity"]) +
		(engagementScore * weights["engagement"])

	// Apply sigmoid for smooth 0-1 output (increased sensitivity)
	return 1 / (1 + math.Exp(-8*(riskScore-0.4)))
}

func (cp *ChurnPredictor) GetRiskFactors(features ChurnFeatures) []string {
	factors := []string{}

	if features.AvgSentiment < -0.2 {
		factors = append(factors, "Low sentiment scores")
	}
	if features.ResponseRate < 0.5 {
		factors = append(factors, "Poor survey participation")
	}
	if features.DaysInactive > 7 {
		factors = append(factors, "Reduced activity")
	}
	if features.TotalResponses > 0 && float64(features.NegativeResponses)/float64(features.TotalResponses) > 0.6 {
		factors = append(factors, "Frequent negative feedback")
	}
	if features.LastLoginDays > 3 {
		factors = append(factors, "Infrequent system usage")
	}

	return factors
}

// Generate AI-powered retention suggestions based on risk factors
func (cp *ChurnPredictor) GenerateRetentionSuggestions(features ChurnFeatures, language string) []string {
	suggestions := map[string]map[string][]string{
		"en": {
			"sentiment": {"Schedule 1-on-1 meeting", "Discuss career goals", "Address work concerns"},
			"response":  {"Send personalized survey", "Improve communication", "Regular check-ins"},
			"activity":  {"Assign engaging projects", "Team collaboration", "Skill development"},
			"engagement": {"Recognition program", "Flexible work options", "Mentorship"},
		},
		"hi": {
			"sentiment": {"व्यक्तिगत बैठक करें", "करियर लक्ष्यों पर चर्चा", "कार्य संबंधी चिंताओं को हल करें"},
			"response":  {"व्यक्तिगत सर्वे भेजें", "संवाद में सुधार", "नियमित जांच"},
			"activity":  {"रोचक प्रोजेक्ट दें", "टीम सहयोग", "कौशल विकास"},
			"engagement": {"पहचान कार्यक्रम", "लचीले काम के विकल्प", "मार्गदर्शन"},
		},
		"ta": {
			"sentiment": {"தனிப்பட்ட சந்திப்பு", "தொழில் இலக்குகள் விவாதம்", "வேலை கவலைகள் தீர்க்க"},
			"response":  {"தனிப்பட்ட கணக்கெடுப்பு", "தொடர்பு மேம்படுத்த", "வழக்கமான சரிபார்ப்பு"},
			"activity":  {"சுவாரஸ்யமான திட்டங்கள்", "குழு ஒத்துழைப்பு", "திறன் மேம்பாடு"},
			"engagement": {"அங்கீகார திட்டம்", "நெகிழ்வான வேலை", "வழிகாட்டுதல்"},
		},
	}

	lang := language
	if _, exists := suggestions[lang]; !exists {
		lang = "en"
	}

	result := []string{}
	
	if features.AvgSentiment < -0.2 {
		result = append(result, suggestions[lang]["sentiment"]...)
	}
	if features.ResponseRate < 0.5 {
		result = append(result, suggestions[lang]["response"]...)
	}
	if features.DaysInactive > 7 {
		result = append(result, suggestions[lang]["activity"]...)
	}
	if features.TotalResponses > 0 && float64(features.NegativeResponses)/float64(features.TotalResponses) > 0.6 {
		result = append(result, suggestions[lang]["engagement"]...)
	}

	// Limit to top 4 suggestions
	if len(result) > 4 {
		result = result[:4]
	}

	return result
}