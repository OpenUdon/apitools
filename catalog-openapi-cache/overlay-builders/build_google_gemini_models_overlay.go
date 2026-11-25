//go:build ignore

package main

import (
	"fmt"
	"os"
)

const overlay = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Google Gemini model listing advisory overlay",
    "version": "v1beta",
    "description": "Metadata for listing available Gemini model IDs with pagination."
  },
  "servers": [{"url": "https://generativelanguage.googleapis.com"}],
  "paths": {
    "/v1beta/models": {
      "get": {
        "operationId": "listGeminiModels",
        "summary": "List available Gemini models",
        "parameters": [
          {"name": "pageSize", "in": "query", "description": "Maximum models requested per page. The server returns at most 1,000 even if a larger size is requested. The default is 50.", "schema": {"type": "integer", "format": "int32", "default": 50}},
          {"name": "pageToken", "in": "query", "schema": {"type": "string"}}
        ],
        "responses": {
          "200": {
            "description": "Available models and continuation token",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "models": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "required": ["name", "version"],
                        "properties": {
                          "name": {"type": "string"},
                          "displayName": {"type": "string"},
                          "version": {"type": "string"},
                          "supportedGenerationMethods": {"type": "array", "items": {"type": "string"}}
                        }
                      }
                    },
                    "nextPageToken": {"type": "string"}
                  }
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {"geminiApiKey": {"type": "apiKey", "in": "header", "name": "x-goog-api-key"}}
  },
  "security": [{"geminiApiKey": []}],
  "x-apitools-overlay": {
    "provider_id": "google-gemini",
    "overlay_id": "google-gemini-models-overlay",
    "official_openapi": false,
    "derived_from_docs": true,
    "source_refs": [
      "https://aistudio.google.com/api-reference/models/list-models",
      "https://aistudio.google.com/api-reference/models",
      "https://aistudio.google.com/docs/api-key"
    ],
    "source_note": "Advisory metadata for the documented v1beta model listing operation, with selected response fields. The Google Discovery source remains the primary catalog reference."
  }
}
`

func main() {
	const path = "catalog-openapi-cache/advisory-overlays/google-gemini-models-overlay.json"
	if err := os.WriteFile(path, []byte(overlay), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
