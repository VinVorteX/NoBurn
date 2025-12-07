package models

import (
	"encoding/json"
	"testing"
)

func TestStringArrayScan(t *testing.T) {
	var arr StringArray
	
	jsonData := []byte(`["item1", "item2", "item3"]`)
	
	err := arr.Scan(jsonData)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if len(arr) != 3 {
		t.Errorf("Expected 3 items, got %d", len(arr))
	}

	if arr[0] != "item1" {
		t.Errorf("Expected item1, got %s", arr[0])
	}
}

func TestStringArrayValue(t *testing.T) {
	arr := StringArray{"item1", "item2", "item3"}
	
	value, err := arr.Value()
	if err != nil {
		t.Fatalf("Value failed: %v", err)
	}

	var result []string
	if err := json.Unmarshal(value.([]byte), &result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 items, got %d", len(result))
	}
}

func TestStringArrayScanNil(t *testing.T) {
	var arr StringArray
	
	err := arr.Scan(nil)
	if err != nil {
		t.Fatalf("Scan nil failed: %v", err)
	}

	if arr != nil {
		t.Error("Expected nil array")
	}
}