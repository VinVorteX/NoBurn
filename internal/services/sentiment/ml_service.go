package sentiment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type MLRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

type MLResponse struct {
	Sentiment float64 `json:"sentiment"`
	Confidence float64 `json:"confidence"`
}

type MLService struct {
	apiURL string
	client *http.Client
}

func NewMLService(apiURL string) *MLService {
	return &MLService{
		apiURL: apiURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *MLService) AnalyzeSentiment(text, language string) (float64, error) {
	// Primary: Use IndicBERT for Indian languages
	score, err := s.AnalyzeWithIndicBERT(text, language)
	if err == nil && score != 0 {
		return score, nil
	}
	
	// Fallback: Rule-based analysis
	return AnalyzeSentiment(text, language), nil
}

// IndicBERT for Indian language sentiment (Primary model)
func (s *MLService) AnalyzeWithIndicBERT(text, language string) (float64, error) {
	// Use ai4bharat/indic-bert for all Indian languages
	model := "ai4bharat/indic-bert"
	
	// For English, use specialized model
	if language == "en" {
		model = "cardiffnlp/twitter-roberta-base-sentiment-latest"
	}


	reqBody := map[string]interface{}{
		"inputs": text,
	}

	jsonData, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", 
		fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model),
		bytes.NewBuffer(jsonData))
	
	req.Header.Set("Authorization", "Bearer "+s.apiURL) // apiURL stores HF token
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("HF API error: %v", err)
		return AnalyzeSentiment(text, language), nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("HF Response [%d]: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return AnalyzeSentiment(text, language), nil
	}

	var result [][]map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Parse error: %v", err)
		return AnalyzeSentiment(text, language), nil
	}

	// Convert labels to sentiment score
	if len(result) > 0 && len(result[0]) > 0 {
		var maxScore float64
		var maxLabel string
		
		// Find highest confidence prediction
		for _, pred := range result[0] {
			label, _ := pred["label"].(string)
			score, _ := pred["score"].(float64)
			
			if score > maxScore {
				maxScore = score
				maxLabel = label
			}
		}
		
		// Convert to sentiment score
		switch maxLabel {
		case "LABEL_2", "POSITIVE", "positive":
			return maxScore, nil
		case "LABEL_0", "NEGATIVE", "negative":
			return -maxScore, nil
		case "LABEL_1", "NEUTRAL", "neutral":
			return 0, nil
		}
	}

	// If no valid result, use fallback
	return AnalyzeSentiment(text, language), nil
}