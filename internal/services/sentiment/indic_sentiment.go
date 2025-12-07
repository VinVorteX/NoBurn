package sentiment

import (
	"strings"
)

// Mock sentiment analysis for Indian languages
// In production, integrate with actual ML models or APIs
func AnalyzeSentiment(text string, language string) float64 {
	text = strings.ToLower(text)
	
	// Simple keyword-based sentiment (replace with actual ML model)
	positiveWords := map[string][]string{
		"en": {"good", "great", "excellent", "happy", "satisfied", "love", "wonderful", "amazing", "fantastic", "awesome"},
		"hi": {"अच्छा", "बहुत अच्छा", "खुश", "संतुष्ट", "प्रसन्न"},
	}
	
	negativeWords := map[string][]string{
		"en": {"bad", "terrible", "awful", "sad", "frustrated", "hate"},
		"hi": {"बुरा", "भयानक", "दुखी", "परेशान", "गुस्सा"},
	}
	
	score := 0.0
	wordCount := 0
	
	words := strings.Fields(text)
	for _, word := range words {
		wordCount++
		
		// Check positive words
		for _, pos := range positiveWords[language] {
			if strings.Contains(word, pos) {
				score += 1.0
				break
			}
		}
		
		// Check negative words
		for _, neg := range negativeWords[language] {
			if strings.Contains(word, neg) {
				score -= 1.0
				break
			}
		}
	}
	
	if wordCount == 0 {
		return 0.0
	}
	
	// Normalize to -1 to 1 range
	normalized := score / float64(wordCount)
	if normalized > 1.0 {
		normalized = 1.0
	} else if normalized < -1.0 {
		normalized = -1.0
	}
	
	return normalized
}