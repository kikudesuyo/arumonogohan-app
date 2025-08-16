package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Tool は Gemini Tool 情報をラップ
type Tool struct {
	Name        string
	Description string
	Parameters  ToolParameter
}

// ToolParameter は JSON Schema 用
type ToolParameter map[string]Property

type Property struct {
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Items       *Property `json:"items,omitempty"`
}

// GeminiAI は Gemini 用クライアント
type GeminiAI struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
}

var GeminiModel = "gemini-1.5-flash"

// NewGeminiAI は Gemini クライアントを生成
func NewGeminiAI(ctx context.Context) (*GeminiAI, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("GeminiAI client creation failed: %v", err)
	}
	model := client.GenerativeModel(GeminiModel)
	return &GeminiAI{Client: client, Model: model}, nil
}

func convertToGenaiSchema(params ToolParameter) map[string]*genai.Schema {
	schema := make(map[string]*genai.Schema)
	for k, v := range params {
		var genaiType genai.Type
		switch v.Type {
		case "string":
			genaiType = genai.TypeString
		case "array":
			genaiType = genai.TypeArray
		case "object":
			genaiType = genai.TypeObject
		}
		s := &genai.Schema{Type: genaiType, Description: v.Description}

		if v.Items != nil {
			var itemGenaiType genai.Type
			switch v.Items.Type {
			case "string":
				itemGenaiType = genai.TypeString
			case "integer":
				itemGenaiType = genai.TypeInteger
			case "number":
				itemGenaiType = genai.TypeNumber
			case "boolean":
				itemGenaiType = genai.TypeBoolean
			}
			s.Items = &genai.Schema{Type: itemGenaiType}
		}
		schema[k] = s
	}
	return schema
}

// StructToToolParams は任意の構造体 T から ToolParameter を生成
func StructToToolParams[T any]() ToolParameter {
	t := reflect.TypeOf((*T)(nil)).Elem()
	params := ToolParameter{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typ := f.Type.Kind().String()
		prop := Property{
			Type:        typ,
			Description: f.Tag.Get("json"), // json tag を説明として使う例
		}
		if f.Type.Kind() == reflect.Slice {
			prop.Type = "array"
			prop.Items = &Property{Type: f.Type.Elem().Kind().String()}
		}
		params[f.Name] = prop
	}
	return params
}

// MakeToolFromStruct は構造体 T から Tool を生成
func MakeToolFromStruct[T any](name, desc string) Tool {
	return Tool{
		Name:        name,
		Description: desc,
		Parameters:  StructToToolParams[T](),
	}
}

// GetJSONResp は Gemini の FunctionCall の Args を任意の構造体にマッピング
func GetJSONResp[T any](g *GeminiAI, ctx context.Context, prompt string, tools []Tool) (T, error) {
	var result T
	modelCopy := *g.Model

	var genaiTools []*genai.Tool
	for _, t := range tools {
		genaiTools = append(genaiTools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        t.Name,
					Description: t.Description,
					Parameters: &genai.Schema{
						Type:       genai.TypeObject,
						Properties: convertToGenaiSchema(t.Parameters),
					},
				},
			},
		})
	}
	modelCopy.Tools = genaiTools

	resp, err := modelCopy.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return result, fmt.Errorf("error generating content: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return result, fmt.Errorf("no content returned from AI")
	}

	fc, ok := resp.Candidates[0].Content.Parts[0].(genai.FunctionCall)
	if !ok {
		return result, fmt.Errorf("unexpected content type: %T", resp.Candidates[0].Content.Parts[0])
	}

	jsonBytes, err := json.Marshal(fc.Args)
	if err != nil {
		return result, fmt.Errorf("failed to marshal function call args: %v", err)
	}

	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal json to struct: %v", err)
	}

	return result, nil
}
