package infrastructure

import (
	"testing"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"google.golang.org/genai"
)

func TestStructToToolParamsUsesJSONFieldNames(t *testing.T) {
	params := StructToToolParams[entity.RecipeInputResp]()

	if _, ok := params["title"]; !ok {
		t.Fatal("StructToToolParams() did not use the JSON field name title")
	}
	if _, ok := params["Title"]; ok {
		t.Fatal("StructToToolParams() unexpectedly used the Go field name Title")
	}
	if params["ingredients"].Type != "array" || params["ingredients"].Items == nil {
		t.Fatalf("ingredients schema = %#v, want an array with item schema", params["ingredients"])
	}
}

func TestConvertToGenaiSchema(t *testing.T) {
	schema := convertToGenaiSchema(ToolParameter{
		"ingredients": {Type: "array", Items: &Property{Type: "string"}},
		"title":       {Type: "string"},
	})

	if schema["title"].Type != genai.TypeString {
		t.Fatalf("title type = %q, want %q", schema["title"].Type, genai.TypeString)
	}
	if schema["ingredients"].Type != genai.TypeArray || schema["ingredients"].Items.Type != genai.TypeString {
		t.Fatalf("ingredients schema = %#v, want string array", schema["ingredients"])
	}
}
