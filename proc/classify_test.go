package proc

import (
	"testing"
)

func TestClassifier(t *testing.T) {
	categories := []Category{
		{
			Name:      "Technology",
			Keywords:  []string{"computer", "software", "hardware"},
			Phrases:   []string{"artificial intelligence", "machine learning"},
			Contexts:  map[string][]string{"data": {"analysis", "processing"}},
			Excluders: []string{"biology"},
		},
		{
			Name:     "Science",
			Keywords: []string{"research", "experiment", "laboratory"},
			Phrases:  []string{"scientific method", "hypothesis testing"},
		},
	}

	classifier := &Classifier{}
	classifier.Init(categories)

	tests := []struct {
		name               string
		input              string
		expectedCategory   string
		expectedConfidence bool // true if confidence should be > 0
		expectedMatches    bool // true if matches should not be empty
	}{
		{
			name:               "Basic keyword match",
			input:              "The computer is running new software",
			expectedCategory:   "Technology",
			expectedConfidence: true,
			expectedMatches:    true,
		},
		{
			name:               "Phrase match",
			input:              "Artificial intelligence is transforming industries",
			expectedCategory:   "Technology",
			expectedConfidence: true,
			expectedMatches:    true,
		},
		{
			name:               "Context match",
			input:              "The data analysis shows interesting patterns",
			expectedCategory:   "Technology",
			expectedConfidence: true,
			expectedMatches:    true,
		},
		{
			name:               "Excluder test",
			input:              "The biology of computer systems",
			expectedCategory:   "Unknown",
			expectedConfidence: false,
			expectedMatches:    false,
		},
		{
			name:               "Stop words handling",
			input:              "The and but or research in laboratory",
			expectedCategory:   "Science",
			expectedConfidence: true,
			expectedMatches:    true,
		},
		{
			name:               "Unknown category",
			input:              "The weather is nice today",
			expectedCategory:   "Unknown",
			expectedConfidence: false,
			expectedMatches:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.Classify(tt.input)

			if result.Category != tt.expectedCategory {
				t.Errorf("Classify() category = %v, want %v", result.Category, tt.expectedCategory)
			}

			if tt.expectedConfidence && result.Confidence <= 0 {
				t.Errorf("Classify() confidence = %v, want > 0", result.Confidence)
			}

			if tt.expectedMatches && len(result.Matches) == 0 {
				t.Errorf("Classify() matches is empty, want non-empty")
			}
		})
	}
}
