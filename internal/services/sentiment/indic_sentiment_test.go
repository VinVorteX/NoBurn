package sentiment

import "testing"

func TestAnalyzeSentiment(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		language string
		wantPos  bool
	}{
		{
			name:     "Positive English",
			text:     "I am very happy and satisfied with my work",
			language: "en",
			wantPos:  true,
		},
		{
			name:     "Negative English",
			text:     "I am frustrated and sad about the situation",
			language: "en",
			wantPos:  false,
		},
		{
			name:     "Positive Hindi",
			text:     "मैं बहुत खुश हूं",
			language: "hi",
			wantPos:  true,
		},
		{
			name:     "Negative Hindi",
			text:     "मैं बहुत दुखी हूं",
			language: "hi",
			wantPos:  false,
		},
		{
			name:     "Neutral text",
			text:     "The meeting is scheduled for tomorrow",
			language: "en",
			wantPos:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := AnalyzeSentiment(tt.text, tt.language)

			if score < -1 || score > 1 {
				t.Errorf("Sentiment score %f out of range [-1,1]", score)
			}

			if tt.wantPos && score <= 0 {
				t.Errorf("Expected positive sentiment, got %f", score)
			}

			if !tt.wantPos && tt.text != "The meeting is scheduled for tomorrow" && score >= 0 {
				t.Errorf("Expected negative sentiment, got %f", score)
			}
		})
	}
}