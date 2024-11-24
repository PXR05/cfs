# CFS - Text Classification

A lightweight text classification service written in Go that categorizes text input using keyword matching, phrase detection, contextual analysis and exclusion rules.

## Features

- Rule-based text classification
- Multiple classification methods:
  - Keyword matching
  - Exact phrase detection
  - Contextual analysis
  - Exclusion rules
- SQLite persistence
- REST API
- Confidence scoring

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/PXR05/cfs.git
   cd cfs
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Build the application:

   ```sh
   go build -o cfs
   ```

## Usage

### Starting the Server

```sh
./cfs
```

The server starts on port 8080 by default.

### Classification Rules

Categories are defined with the following structure:

```json
{
  "Name": "Category",
  "Keywords": ["word1", "word2"],
  "Phrases": ["exact phrase1", "exact phrase2"],
  "Contexts": {
    "context1": ["related1", "related2"]
  },
  "Excluders": ["exclude1", "exclude2"]
}
```

- Keywords: Single words that indicate the category
- Phrases: Exact phrases to match
- Contexts: Related words that increase confidence when found together
- Excluders: Words that disqualify a text from a category

## API Reference

### Categories

#### Create Categories

- `POST /cfs/c`
- Request body: Array of category objects
- Response: 201 Created

#### Get Categories

- `GET /cfs/c`
- Response: List of all categories

#### Get Category

- `GET /cfs/c?category={name}`
- Response: Single category object

### Classifications

#### Create Classifications

- `POST /cfs/i`
- Request body: `{"Items": ["text1", "text2"]}`
- Response: Classification results

#### Get Classifications

- `GET /cfs/i`
- Response: All stored classifications

#### Get Classification

- `GET /cfs/i?item={text}`
- Response: Classification result for specific text

## Example

Creating a category:

```sh
curl -X POST http://localhost:8080/cfs/c -d '[{
    "Name": "Technology",
    "Keywords": ["computer", "software"],
    "Phrases": ["artificial intelligence"],
    "Contexts": {
        "data": ["analysis"]
    },
    "Excluders": ["biology"]
}]'
```

Classifying text:

```sh
curl -X POST http://localhost:8080/cfs/i -d '{
    "Items": ["The computer is running new software"]
}'
```

## Testing

Run the test suite:

```sh
go test ./...
```

## License

MIT License
