//go:build ignore

package main

import (
	"fmt"
	"os"
)

const overlay = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Twitch Helix Get Games advisory overlay",
    "version": "helix",
    "description": "Documents the game and category lookup and the returned identifiers, name and box art template."
  },
  "servers": [
    {
      "url": "https://api.twitch.tv"
    }
  ],
  "paths": {
    "/helix/games": {
      "get": {
        "operationId": "getTwitchGames",
        "summary": "Get information about specified categories or games",
        "description": "Supply at least one ID, exact name or IGDB ID. Send at most 100 values across the three query parameters.",
        "parameters": [
          {
            "name": "id",
            "in": "query",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              },
              "maxItems": 100
            },
            "description": "Twitch category or game ID.",
            "style": "form",
            "explode": true
          },
          {
            "name": "name",
            "in": "query",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              },
              "maxItems": 100
            },
            "description": "Exact Twitch category or game title.",
            "style": "form",
            "explode": true
          },
          {
            "name": "igdb_id",
            "in": "query",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              },
              "maxItems": 100
            },
            "description": "IGDB game ID.",
            "style": "form",
            "explode": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "id": {
                            "type": "string"
                          },
                          "name": {
                            "type": "string"
                          },
                          "box_art_url": {
                            "type": "string"
                          },
                          "igdb_id": {
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
      }
    }
  },
  "components": {
    "securitySchemes": {
      "twitchAppOAuth": {
        "type": "oauth2",
        "flows": {
          "clientCredentials": {
            "tokenUrl": "https://id.twitch.tv/oauth2/token",
            "scopes": {}
          }
        }
      },
      "twitchUserOAuth": {
        "type": "oauth2",
        "flows": {
          "authorizationCode": {
            "authorizationUrl": "https://id.twitch.tv/oauth2/authorize",
            "tokenUrl": "https://id.twitch.tv/oauth2/token",
            "scopes": {}
          }
        }
      },
      "twitchClientID": {
        "type": "apiKey",
        "in": "header",
        "name": "Client-Id",
        "description": "Registered client identifier, which must match the identifier associated with the OAuth token."
      }
    }
  },
  "security": [
    {
      "twitchAppOAuth": [],
      "twitchClientID": []
    },
    {
      "twitchUserOAuth": [],
      "twitchClientID": []
    }
  ],
  "x-apitools-overlay": {
    "provider_id": "twitch",
    "overlay_id": "twitch-helix-games-overlay",
    "official_openapi": false,
    "derived_from_docs": true,
    "source_refs": [
      "https://dev.twitch.tv/docs/api/reference/games#get-games",
      "https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/",
      "https://legal.twitch.com/en/legal/developer-agreement/"
    ],
    "source_note": "Documents the game and category lookup and the returned identifiers, name and box art template."
  }
}
`

func main() {
	const path = "catalog-openapi-cache/advisory-overlays/twitch-helix-games-overlay.json"
	if err := os.WriteFile(path, []byte(overlay), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
