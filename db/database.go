package db

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"

	"cfs/proc"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Init() error {
	var err error
	d.db, err = sql.Open("sqlite3", "cfs.db")
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS classifications (
			item TEXT PRIMARY KEY,
			category TEXT NOT NULL,
			confidence REAL NOT NULL,
			matches TEXT
		);
		CREATE TABLE IF NOT EXISTS categories (
			name TEXT PRIMARY KEY,
			keywords TEXT NOT NULL,
			phrases TEXT NOT NULL,
			contexts TEXT NOT NULL,
			excluders TEXT NOT NULL
		);
	`)
	return err
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *Database) Seed() {
	var categories = []proc.Category{
		{
			Name: "Technology",
			Keywords: []string{
				"computer", "software", "program", "code", "algorithm",
				"database", "network", "server", "application", "system",
			},
			Phrases: []string{
				"artificial intelligence",
				"machine learning",
				"deep learning",
				"neural network",
				"cloud computing",
			},
			Contexts: map[string][]string{
				"development": {"software", "web", "app", "mobile"},
				"data":        {"processing", "analysis", "storage"},
				"security":    {"cyber", "network", "encryption"},
			},
			Excluders: []string{"recipe", "cook", "bake", "ingredient"},
		},
		{
			Name: "Food and Cooking",
			Keywords: []string{
				"cook", "recipe", "food", "meal", "ingredient",
				"kitchen", "dish", "taste", "flavor", "cuisine",
			},
			Phrases: []string{
				"healthy eating",
				"meal prep",
				"cooking instructions",
				"recipe guide",
				"food preparation",
			},
			Contexts: map[string][]string{
				"preparation": {"cook", "bake", "grill", "roast"},
				"ingredients": {"fresh", "organic", "raw", "dried"},
				"taste":       {"delicious", "savory", "sweet", "spicy"},
			},
			Excluders: []string{"computer", "program", "code", "algorithm"},
		},
	}

	for _, category := range categories {
		if err := d.AddCategory(category); err != nil {
			panic(err)
		}
	}
}

func (d *Database) AddClassification(item string, result proc.ClassificationResult) error {
	matchesJSON, err := json.Marshal(result.Matches)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(
		"INSERT OR REPLACE INTO classifications (item, category, confidence, matches) VALUES (?, ?, ?, ?)",
		item, result.Category, result.Confidence, string(matchesJSON),
	)
	return err
}

func (d *Database) GetClassifications() ([]proc.ClassificationResult, error) {
	rows, err := d.db.Query("SELECT item, category, confidence, matches FROM classifications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []proc.ClassificationResult
	for rows.Next() {
		var result proc.ClassificationResult
		var matchesJSON string
		err := rows.Scan(&result.Item, &result.Category, &result.Confidence, &matchesJSON)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(matchesJSON), &result.Matches); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, rows.Err()
}

func (d *Database) GetClassification(item string) (proc.ClassificationResult, error) {
	var result proc.ClassificationResult
	var matchesJSON string
	err := d.db.QueryRow(
		"SELECT item, category, confidence, matches FROM classifications WHERE item = ?",
		item,
	).Scan(&result.Item, &result.Category, &result.Confidence, &matchesJSON)
	if err != nil {
		return proc.ClassificationResult{}, err
	}
	if err := json.Unmarshal([]byte(matchesJSON), &result.Matches); err != nil {
		return proc.ClassificationResult{}, err
	}
	return result, nil
}

func (d *Database) AddCategory(category proc.Category) error {
	keywords, err := json.Marshal(category.Keywords)
	if err != nil {
		return err
	}
	phrases, err := json.Marshal(category.Phrases)
	if err != nil {
		return err
	}
	contexts, err := json.Marshal(category.Contexts)
	if err != nil {
		return err
	}
	excluders, err := json.Marshal(category.Excluders)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(
		"INSERT OR REPLACE INTO categories (name, keywords, phrases, contexts, excluders) VALUES (?, ?, ?, ?, ?)",
		category.Name, string(keywords), string(phrases), string(contexts), string(excluders),
	)
	return err
}

func (d *Database) GetCategories() ([]proc.Category, error) {
	rows, err := d.db.Query("SELECT name, keywords, phrases, contexts, excluders FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []proc.Category
	for rows.Next() {
		var cat proc.Category
		var keywordsJSON, phrasesJSON, contextsJSON, excludersJSON string
		err := rows.Scan(&cat.Name, &keywordsJSON, &phrasesJSON, &contextsJSON, &excludersJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(keywordsJSON), &cat.Keywords); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(phrasesJSON), &cat.Phrases); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(contextsJSON), &cat.Contexts); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(excludersJSON), &cat.Excluders); err != nil {
			return nil, err
		}

		categories = append(categories, cat)
	}
	return categories, rows.Err()
}

func (d *Database) GetCategory(name string) (proc.Category, error) {
	var cat proc.Category
	var keywordsJSON, phrasesJSON, contextsJSON, excludersJSON string
	err := d.db.QueryRow(
		"SELECT name, keywords, phrases, contexts, excluders FROM categories WHERE name = ?",
		name,
	).Scan(&cat.Name, &keywordsJSON, &phrasesJSON, &contextsJSON, &excludersJSON)
	if err != nil {
		return proc.Category{}, err
	}

	if err := json.Unmarshal([]byte(keywordsJSON), &cat.Keywords); err != nil {
		return proc.Category{}, err
	}
	if err := json.Unmarshal([]byte(phrasesJSON), &cat.Phrases); err != nil {
		return proc.Category{}, err
	}
	if err := json.Unmarshal([]byte(contextsJSON), &cat.Contexts); err != nil {
		return proc.Category{}, err
	}
	if err := json.Unmarshal([]byte(excludersJSON), &cat.Excluders); err != nil {
		return proc.Category{}, err
	}

	return cat, nil
}

func (d *Database) Cleanup() error {
	_, err := d.db.Exec(`
		DELETE FROM classifications WHERE item LIKE 'test%';
		DELETE FROM categories WHERE name LIKE 'Test%';
	`)
	return err
}
