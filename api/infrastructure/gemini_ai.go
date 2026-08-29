package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"google.golang.org/genai"
)

type Tool struct {
	Name        string
	Description string
	Parameters  ToolParameter
}

type ToolParameter map[string]Property

type Property struct {
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Items       *Property `json:"items,omitempty"`
}

type GeminiAI struct {
	Client *genai.Client
	Model  string
}

const GeminiModel = "gemini-3.5-flash-lite"

func NewGeminiAI(ctx context.Context) (*GeminiAI, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("GeminiAI client creation failed: %w", err)
	}
	return &GeminiAI{Client: client, Model: GeminiModel}, nil
}

func convertToGenaiSchema(params ToolParameter) map[string]*genai.Schema {
	schema := make(map[string]*genai.Schema, len(params))
	for name, property := range params {
		schema[name] = propertySchema(property)
	}
	return schema
}

func propertySchema(property Property) *genai.Schema {
	schema := &genai.Schema{Type: genaiType(property.Type), Description: property.Description}
	if property.Items != nil {
		schema.Items = propertySchema(*property.Items)
	}
	return schema
}

func genaiType(kind string) genai.Type {
	switch kind {
	case "string":
		return genai.TypeString
	case "integer":
		return genai.TypeInteger
	case "number", "float64", "float32":
		return genai.TypeNumber
	case "boolean":
		return genai.TypeBoolean
	case "array", "slice":
		return genai.TypeArray
	case "object", "struct":
		return genai.TypeObject
	default:
		return genai.TypeString
	}
}

func StructToToolParams[T any]() ToolParameter {
	t := reflect.TypeOf((*T)(nil)).Elem()
	params := ToolParameter{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName == "-" {
			continue
		}
		if jsonName == "" {
			jsonName = field.Name
		}
		property := Property{Type: field.Type.Kind().String(), Description: field.Name}
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			property.Type = "array"
			property.Items = &Property{Type: field.Type.Elem().Kind().String()}
		}
		params[jsonName] = property
	}
	return params
}

func MakeToolFromStruct[T any](name, desc string) Tool {
	return Tool{Name: name, Description: desc, Parameters: StructToToolParams[T]()}
}

func GetGeminiJSONResp[T any](g *GeminiAI, ctx context.Context, prompt string, tools []Tool) (T, error) {
	genaiTools := make([]*genai.Tool, 0, len(tools))
	for _, tool := range tools {
		genaiTools = append(genaiTools, &genai.Tool{FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters: &genai.Schema{
				Type:       genai.TypeObject,
				Properties: convertToGenaiSchema(tool.Parameters),
			},
		}}})
	}

	resp, err := g.Client.Models.GenerateContent(ctx, g.Model, genai.Text(prompt), &genai.GenerateContentConfig{Tools: genaiTools})
	if err != nil {
		return *new(T), fmt.Errorf("error generating content: %w", err)
	}
	functionCalls := resp.FunctionCalls()
	if len(functionCalls) == 0 {
		return *new(T), fmt.Errorf("no function call returned from AI")
	}

	jsonBytes, err := json.Marshal(functionCalls[0].Args)
	if err != nil {
		return *new(T), fmt.Errorf("failed to marshal function call args: %w", err)
	}
	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return *new(T), fmt.Errorf("failed to unmarshal json to struct: %w", err)
	}
	return result, nil
}

func GetGeminiStringResp(g *GeminiAI, ctx context.Context, prompt string) (string, error) {
	resp, err := g.Client.Models.GenerateContent(ctx, g.Model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}
	text := strings.TrimSpace(resp.Text())
	if text == "" {
		return "", fmt.Errorf("no text returned from AI")
	}
	return text, nil
}
