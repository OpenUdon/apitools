//go:build ignore

package main

import (
	"fmt"
	"os"
)

const overlay = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Google Ads accessible customers advisory overlay",
    "version": "v26",
    "description": "Documents the published REST binding for the single CustomerService listing RPC. The protobuf document is the primary source for that RPC."
  },
  "servers": [
    {
      "url": "https://googleads.googleapis.com"
    }
  ],
  "paths": {
    "/v26/customers:listAccessibleCustomers": {
      "get": {
        "operationId": "listAccessibleCustomers",
        "summary": "List accessible customers",
        "description": "Return resource names for customers directly accessible to the authenticated principal.",
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "resourceNames": {
                      "type": "array",
                      "items": {
                        "type": "string"
                      }
                    }
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
    "securitySchemes": {
      "googleAdsOAuth": {
        "type": "oauth2",
        "flows": {
          "authorizationCode": {
            "authorizationUrl": "https://accounts.google.com/o/oauth2/v2/auth",
            "tokenUrl": "https://oauth2.googleapis.com/token",
            "scopes": {
              "https://www.googleapis.com/auth/adwords": "Access Google Ads"
            }
          }
        }
      },
      "googleAdsDeveloperToken": {
        "type": "apiKey",
        "in": "header",
        "name": "developer-token"
      }
    }
  },
  "security": [
    {
      "googleAdsOAuth": [
        "https://www.googleapis.com/auth/adwords"
      ],
      "googleAdsDeveloperToken": []
    }
  ],
  "x-apitools-overlay": {
    "provider_id": "google-ads",
    "overlay_id": "google-ads-v26-accessible-customers-overlay",
    "official_openapi": false,
    "derived_from_docs": true,
    "source_refs": [
      "https://developers.google.com/google-ads/api/rest/auth",
      "https://developers.google.com/google-ads/api/docs/concepts/call-structure",
      "https://developers.google.com/google-ads/api/reference/rpc/v26/CustomerService/ListAccessibleCustomers",
      "https://developers.google.com/identity/protocols/oauth2/web-server"
    ],
    "source_note": "Documents the published REST binding for the single CustomerService listing RPC. The protobuf document is the primary source for that RPC."
  }
}
`

func main() {
	const path = "catalog-openapi-cache/advisory-overlays/google-ads-v26-accessible-customers-overlay.json"
	if err := os.WriteFile(path, []byte(overlay), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
