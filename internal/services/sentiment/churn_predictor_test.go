package sentiment

import "testing"

func TestPredictChurnRisk(t *testing.T) {
	predictor := NewChurnPredictor()

	tests := []struct {
		name     string
		features ChurnFeatures
		wantHigh bool
	}{
		{
			name: "High risk employee",
			features: ChurnFeatures{
				AvgSentiment:      -0.5,
				ResponseRate:      0.3,
				DaysInactive:      15,
				NegativeResponses: 8,
				TotalResponses:    10,
			},
			wantHigh: true,
		},
		{
			name: "Low risk employee",
			features: ChurnFeatures{
				AvgSentiment:      0.7,
				ResponseRate:      0.9,
				DaysInactive:      1,
				NegativeResponses: 1,
				TotalResponses:    10,
			},
			wantHigh: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := predictor.PredictChurnRisk(tt.features)

			if risk < 0 || risk > 1 {
				t.Errorf("Risk score %f out of range [0,1]", risk)
			}

			if tt.wantHigh && risk < 0.5 {
				t.Errorf("Expected high risk, got %f", risk)
			}

			if !tt.wantHigh && risk > 0.5 {
				t.Errorf("Expected low risk, got %f", risk)
			}
		})
	}
}

func TestGetRiskFactors(t *testing.T) {
	predictor := NewChurnPredictor()

	features := ChurnFeatures{
		AvgSentiment:      -0.5,
		ResponseRate:      0.3,
		DaysInactive:      10,
		NegativeResponses: 7,
		TotalResponses:    10,
	}

	factors := predictor.GetRiskFactors(features)

	if len(factors) == 0 {
		t.Error("Expected risk factors to be identified")
	}

	expectedFactors := map[string]bool{
		"Low sentiment scores":       true,
		"Poor survey participation":  true,
		"Reduced activity":           true,
		"Frequent negative feedback": true,
	}

	for _, factor := range factors {
		if !expectedFactors[factor] {
			t.Errorf("Unexpected risk factor: %s", factor)
		}
	}
}

func TestGenerateRetentionSuggestions(t *testing.T) {
	predictor := NewChurnPredictor()

	features := ChurnFeatures{
		AvgSentiment: -0.5,
		ResponseRate: 0.3,
	}

	tests := []struct {
		language string
	}{
		{"en"},
		{"hi"},
		{"ta"},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			suggestions := predictor.GenerateRetentionSuggestions(features, tt.language)

			if len(suggestions) == 0 {
				t.Error("Expected suggestions to be generated")
			}

			if len(suggestions) > 4 {
				t.Errorf("Expected max 4 suggestions, got %d", len(suggestions))
			}
		})
	}
}