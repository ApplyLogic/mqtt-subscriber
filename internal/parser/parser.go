package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

type Parser struct {
	// Add any configuration for the parser here
}

type RawData struct {
	Data []byte
}

type ParsedData struct {
	Timestamp  time.Time              `json:"timestamp"`
	DeviceID   string                 `json:"device_id"`
	Values     map[string]interface{} `json:"values"`
	Raw        []byte                 `json:"raw"`
	RawMessage string                 `json:"raw_message"`
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(data []byte) ([]byte, error) {
	// This is a simplified example parser that extracts a device ID and some values
	// You'll need to replace this with your actual parsing logic

	// Example: assuming first 6 bytes are device ID and the rest is some format
	// that we'll convert to JSON
	if len(data) < 8 {
		return nil, fmt.Errorf("data too short to parse")
	}

	deviceID := fmt.Sprintf("%x", data[:6])

	// Example: assume data after device ID contains key-value pairs
	// This is just a placeholder - replace with actual parsing logic
	values := map[string]interface{}{
		"reading1": int(data[6]),
		"reading2": int(data[7]),
		// Add more readings as needed
	}

	parsedData := ParsedData{
		Timestamp:  time.Now(),
		DeviceID:   deviceID,
		Values:     values,
		Raw:        data,
		RawMessage: string(data),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(parsedData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling to JSON: %v", err)
	}

	return jsonData, nil
}
