package csv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"opportunities/internal/schemas"
)

const (
	openingCSVChunkSize = 100
)

var expectedHeader = []string{"role", "company", "location", "remote", "link", "salary"}

type ParsedOpening struct {
	LineNumber int
	Opening    schemas.Openings
}

type RowError struct {
	LineNumber int
	Message    string
}

func ValidateHeader(content []byte) error {
	reader := csv.NewReader(strings.NewReader(string(content)))
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("invalid csv header: %w", err)
	}

	if len(header) != len(expectedHeader) {
		return fmt.Errorf("invalid csv header. expected %v", expectedHeader)
	}

	for i := range expectedHeader {
		if strings.TrimSpace(strings.ToLower(header[i])) != expectedHeader[i] {
			return fmt.Errorf("invalid csv header. expected %v", expectedHeader)
		}
	}

	return nil
}

func ParseAndValidate(content []byte) ([]ParsedOpening, []RowError, error) {
	reader := csv.NewReader(strings.NewReader(string(content)))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("invalid csv format: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("csv file is empty")
	}

	if err := validateHeaderRow(rows[0]); err != nil {
		return nil, nil, err
	}

	rawRows := rows[1:]
	if len(rawRows) == 0 {
		return []ParsedOpening{}, nil, nil
	}

	parsed := make([]ParsedOpening, 0, len(rawRows))
	rowErrors := make([]RowError, 0)

	for start := 0; start < len(rawRows); start += openingCSVChunkSize {
		end := start + openingCSVChunkSize
		if end > len(rawRows) {
			end = len(rawRows)
		}

		chunkRows := rawRows[start:end]
		chunkResults := make([]chunkParseResult, len(chunkRows))
		wg := sync.WaitGroup{}

		for i, row := range chunkRows {
			index := i
			lineNumber := start + i + 2
			data := row
			wg.Add(1)

			go func() {
				defer wg.Done()
				chunkResults[index] = parseRow(lineNumber, data)
			}()
		}

		wg.Wait()

		for _, result := range chunkResults {
			if result.Err != nil {
				rowErrors = append(rowErrors, RowError{
					LineNumber: result.LineNumber,
					Message:    result.Err.Error(),
				})
				continue
			}

			parsed = append(parsed, ParsedOpening{
				LineNumber: result.LineNumber,
				Opening:    result.Opening,
			})
		}
	}

	return parsed, rowErrors, nil
}

type chunkParseResult struct {
	LineNumber int
	Opening    schemas.Openings
	Err        error
}

func validateHeaderRow(header []string) error {
	if len(header) != len(expectedHeader) {
		return fmt.Errorf("invalid csv header. expected %v", expectedHeader)
	}

	for i := range expectedHeader {
		if strings.TrimSpace(strings.ToLower(header[i])) != expectedHeader[i] {
			return fmt.Errorf("invalid csv header. expected %v", expectedHeader)
		}
	}

	return nil
}

func parseRow(lineNumber int, row []string) chunkParseResult {
	if len(row) != len(expectedHeader) {
		return chunkParseResult{
			LineNumber: lineNumber,
			Err:        fmt.Errorf("invalid column count, expected %d, got %d", len(expectedHeader), len(row)),
		}
	}

	role := strings.TrimSpace(row[0])
	company := strings.TrimSpace(row[1])
	location := strings.TrimSpace(row[2])
	remoteRaw := strings.TrimSpace(row[3])
	link := strings.TrimSpace(row[4])
	salaryRaw := strings.TrimSpace(row[5])

	if role == "" {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("role is required")}
	}

	if company == "" {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("company is required")}
	}

	if location == "" {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("location is required")}
	}

	if link == "" {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("link is required")}
	}

	remote, err := strconv.ParseBool(remoteRaw)
	if err != nil {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("remote must be a boolean")}
	}

	salary, err := strconv.ParseInt(salaryRaw, 10, 64)
	if err != nil {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("salary must be an integer")}
	}

	if salary <= 0 {
		return chunkParseResult{LineNumber: lineNumber, Err: fmt.Errorf("salary must be greater than zero")}
	}

	return chunkParseResult{
		LineNumber: lineNumber,
		Opening: schemas.Openings{
			Role:     role,
			Company:  company,
			Location: location,
			Remote:   remote,
			Link:     link,
			Salary:   salary,
		},
	}
}
