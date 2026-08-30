package entity

import (
	"encoding/json"
	"testing"
)

func TestRecipeInputRespUnmarshalAcceptsStringBoolean(t *testing.T) {
	var got RecipeInputResp
	if err := json.Unmarshal([]byte(`{"title":"炒め物","is_tempered":"false"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.IsTempered {
		t.Fatal("expected string false to be accepted")
	}
}

func TestRecipeInputRespUnmarshalAcceptsBoolean(t *testing.T) {
	var got RecipeInputResp
	if err := json.Unmarshal([]byte(`{"title":"炒め物","is_tempered":true}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !got.IsTempered {
		t.Fatal("expected boolean true to be accepted")
	}
}
