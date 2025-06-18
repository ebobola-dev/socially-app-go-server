package nullable

import (
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	Present bool
	Value   *T
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Present = true
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("Nullable.Unmarshal: %w", err)
	}
	n.Value = &v
	return nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Present {
		return []byte("null"), nil
	}
	if n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}

func StringEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
