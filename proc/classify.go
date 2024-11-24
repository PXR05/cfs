package proc

import (
	"fmt"
	"strings"
	"unicode"
)

type Classifier struct {
	categories []Category
	stopWords  map[string]bool
}

func (sc *Classifier) Init(categories []Category) {
	sc.categories = categories
	sc.stopWords = makeStopWords()
}

func makeStopWords() map[string]bool {
	words := []string{"the", "is", "at", "which", "on", "a", "an", "and", "or", "but", "in", "with", "to", "for"}
	stopWords := make(map[string]bool)
	for _, word := range words {
		stopWords[word] = true
	}
	return stopWords
}

func (sc *Classifier) tokenize(text string) []string {
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	filtered := make([]string, 0)
	for _, word := range words {
		if !sc.stopWords[word] {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func (sc *Classifier) Classify(sentence string) ClassificationResult {
	words := sc.tokenize(sentence)
	sentenceLower := strings.ToLower(sentence)

	results := make([]ClassificationResult, 0)

	for _, category := range sc.categories {
		score := 0.0
		matches := make([]string, 0)

		// excluded
		excluded := false
		for _, excluder := range category.Excluders {
			if strings.Contains(sentenceLower, strings.ToLower(excluder)) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		// keywords
		for _, keyword := range category.Keywords {
			if strings.Contains(sentenceLower, strings.ToLower(keyword)) {
				score += 1.0
				matches = append(matches, keyword)
			}
		}

		// phrases
		for _, phrase := range category.Phrases {
			if strings.Contains(sentenceLower, strings.ToLower(phrase)) {
				score += 2.0
				matches = append(matches, phrase)
			}
		}

		// contextual
		for context, relatedWords := range category.Contexts {
			if strings.Contains(sentenceLower, context) {
				for _, related := range relatedWords {
					if strings.Contains(sentenceLower, related) {
						score += 1.5
						matches = append(matches, fmt.Sprintf("%s-%s", context, related))
					}
				}
			}
		}

		// normalize
		confidence := score / float64(len(words))
		if confidence > 0 {
			results = append(results, ClassificationResult{
				Category:   category.Name,
				Confidence: confidence,
				Matches:    matches,
			})
		}
	}

	if len(results) > 0 {
		bestResult := results[0]
		for _, result := range results {
			if result.Confidence > bestResult.Confidence {
				bestResult = result
			}
		}
		return bestResult
	}

	return ClassificationResult{
		Category:   "Unknown",
		Confidence: 0.0,
		Matches:    nil,
	}
}
