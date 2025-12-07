package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type AIService struct {
	hfToken string
	client  *http.Client
}

func NewAIService(hfToken string) *AIService {
	return &AIService{
		hfToken: hfToken,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

type RetentionSuggestion struct {
	Suggestion string `json:"suggestion"`
	Priority   string `json:"priority"`
	Category   string `json:"category"`
}

// Generate AI-powered retention suggestions using LLM
func (s *AIService) GenerateRetentionSuggestions(riskFactors []string, language string) ([]RetentionSuggestion, error) {
	prompt := s.buildPrompt(riskFactors, language)
	
	// Try Hugging Face API first
	suggestions, err := s.callHuggingFace(prompt, language)
	if err == nil {
		return suggestions, nil
	}
	
	// Fallback to rule-based suggestions
	return s.getRuleBasedSuggestions(riskFactors, language), nil
}

func (s *AIService) buildPrompt(riskFactors []string, language string) string {
	langMap := map[string]string{
		"en": "English",
		"hi": "Hindi",
		"ta": "Tamil",
	}
	
	langName := langMap[language]
	if langName == "" {
		langName = "English"
	}
	
	factorsStr := strings.Join(riskFactors, ", ")
	
	prompts := map[string]string{
		"en": fmt.Sprintf("As an HR expert, suggest 4 specific retention strategies for an employee showing these risk factors: %s. Be concise and actionable.", factorsStr),
		"hi": fmt.Sprintf("एक HR विशेषज्ञ के रूप में, इन जोखिम कारकों वाले कर्मचारी के लिए 4 विशिष्ट प्रतिधारण रणनीतियाँ सुझाएं: %s। संक्षिप्त और कार्रवाई योग्य रहें।", factorsStr),
		"ta": fmt.Sprintf("HR நிபுணராக, இந்த ஆபத்து காரணிகளைக் கொண்ட ஊழியருக்கு 4 குறிப்பிட்ட தக்கவைப்பு உத்திகளை பரிந்துரைக்கவும்: %s. சுருக்கமாகவும் செயல்படக்கூடியதாகவும் இருக்கவும்.", factorsStr),
	}
	
	return prompts[language]
}

func (s *AIService) callHuggingFace(prompt, language string) ([]RetentionSuggestion, error) {
	model := "mistralai/Mistral-7B-Instruct-v0.2"
	
	reqBody := map[string]interface{}{
		"inputs": prompt,
		"parameters": map[string]interface{}{
			"max_new_tokens": 200,
			"temperature":    0.7,
		},
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model),
		bytes.NewBuffer(jsonData))
	
	req.Header.Set("Authorization", "Bearer "+s.hfToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HF API error: %d", resp.StatusCode)
	}
	
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("no response from HF")
	}
	
	// Parse generated text into suggestions
	generatedText := result[0]["generated_text"].(string)
	return s.parseAISuggestions(generatedText), nil
}

func (s *AIService) parseAISuggestions(text string) []RetentionSuggestion {
	lines := strings.Split(text, "\n")
	suggestions := []RetentionSuggestion{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || len(line) < 10 {
			continue
		}
		
		// Remove numbering like "1.", "2.", etc.
		line = strings.TrimLeft(line, "0123456789.-) ")
		
		if len(suggestions) < 4 {
			suggestions = append(suggestions, RetentionSuggestion{
				Suggestion: line,
				Priority:   "medium",
				Category:   "ai_generated",
			})
		}
	}
	
	return suggestions
}

func (s *AIService) getRuleBasedSuggestions(riskFactors []string, language string) []RetentionSuggestion {
	templates := map[string]map[string][]string{
		"en": {
			"Low sentiment scores":        {"Schedule 1-on-1 meeting", "Discuss career goals", "Address work concerns"},
			"Poor survey participation":   {"Send personalized survey", "Improve communication", "Regular check-ins"},
			"Reduced activity":            {"Assign engaging projects", "Team collaboration", "Skill development"},
			"Frequent negative feedback":  {"Recognition program", "Flexible work options", "Mentorship"},
			"Infrequent system usage":     {"Re-engagement program", "Training sessions", "Team activities"},
		},
		"hi": {
			"Low sentiment scores":        {"व्यक्तिगत बैठक करें", "करियर लक्ष्यों पर चर्चा", "कार्य संबंधी चिंताओं को हल करें"},
			"Poor survey participation":   {"व्यक्तिगत सर्वे भेजें", "संवाद में सुधार", "नियमित जांच"},
			"Reduced activity":            {"रोचक प्रोजेक्ट दें", "टीम सहयोग", "कौशल विकास"},
			"Frequent negative feedback":  {"पहचान कार्यक्रम", "लचीले काम के विकल्प", "मार्गदर्शन"},
			"Infrequent system usage":     {"पुनः जुड़ाव कार्यक्रम", "प्रशिक्षण सत्र", "टीम गतिविधियाँ"},
		},
		"ta": {
			"Low sentiment scores":        {"தனிப்பட்ட சந்திப்பு", "தொழில் இலக்குகள் விவாதம்", "வேலை கவலைகள் தீர்க்க"},
			"Poor survey participation":   {"தனிப்பட்ட கணக்கெடுப்பு", "தொடர்பு மேம்படுத்த", "வழக்கமான சரிபார்ப்பு"},
			"Reduced activity":            {"சுவாரஸ்யமான திட்டங்கள்", "குழு ஒத்துழைப்பு", "திறன் மேம்பாடு"},
			"Frequent negative feedback":  {"அங்கீகார திட்டம்", "நெகிழ்வான வேலை", "வழிகாட்டுதல்"},
			"Infrequent system usage":     {"மீண்டும் ஈடுபாடு திட்டம்", "பயிற்சி அமர்வுகள்", "குழு நடவடிக்கைகள்"},
		},
	}
	
	lang := language
	if _, exists := templates[lang]; !exists {
		lang = "en"
	}
	
	suggestions := []RetentionSuggestion{}
	seen := make(map[string]bool)
	
	for _, factor := range riskFactors {
		if suggList, exists := templates[lang][factor]; exists {
			for _, sugg := range suggList {
				if !seen[sugg] && len(suggestions) < 4 {
					suggestions = append(suggestions, RetentionSuggestion{
						Suggestion: sugg,
						Priority:   s.getPriority(factor),
						Category:   factor,
					})
					seen[sugg] = true
				}
			}
		}
	}
	
	return suggestions
}

func (s *AIService) getPriority(factor string) string {
	highPriority := []string{"Low sentiment scores", "Frequent negative feedback"}
	for _, hp := range highPriority {
		if factor == hp {
			return "high"
		}
	}
	return "medium"
}

// Translate text using AI
func (s *AIService) TranslateText(text, targetLang string) (string, error) {
	modelMap := map[string]string{
		"hi": "Helsinki-NLP/opus-mt-en-hi",
		"ta": "Helsinki-NLP/opus-mt-en-ta",
	}
	
	model, exists := modelMap[targetLang]
	if !exists {
		return text, nil // Return original if no translation needed
	}
	
	reqBody := map[string]interface{}{
		"inputs": text,
	}
	
	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model),
		bytes.NewBuffer(jsonData))
	
	req.Header.Set("Authorization", "Bearer "+s.hfToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return text, err
	}
	defer resp.Body.Close()
	
	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return text, err
	}
	
	if len(result) > 0 {
		if translated, ok := result[0]["translation_text"].(string); ok {
			return translated, nil
		}
	}
	
	return text, nil
}