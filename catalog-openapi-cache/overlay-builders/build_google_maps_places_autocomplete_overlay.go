//go:build ignore

package main

import (
	"fmt"
	"os"
)

const overlay = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Places API (New) Autocomplete advisory overlay",
    "version": "v1",
    "description": "Documents the default text prediction request, select optional inputs, prediction identifiers and display text. The full Google Discovery document is the primary source."
  },
  "servers": [
    {
      "url": "https://places.googleapis.com"
    }
  ],
  "paths": {
    "/v1/places:autocomplete": {
      "post": {
        "operationId": "autocompletePlaces",
        "summary": "Return place and query predictions for the given input",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "input": {
                    "type": "string",
                    "description": "Text used to produce predictions."
                  },
                  "languageCode": {
                    "type": "string"
                  },
                  "regionCode": {
                    "type": "string"
                  },
                  "includeQueryPredictions": {
                    "type": "boolean"
                  },
                  "sessionToken": {
                    "type": "string"
                  }
                },
                "required": [
                  "input"
                ]
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "suggestions": {
                      "type": "array",
                      "items": {
                        "type": "object",
                        "properties": {
                          "placePrediction": {
                            "type": "object",
                            "properties": {
                              "place": {
                                "type": "string"
                              },
                              "placeId": {
                                "type": "string"
                              },
                              "text": {
                                "type": "object",
                                "properties": {
                                  "text": {
                                    "type": "string"
                                  }
                                }
                              }
                            }
                          },
                          "queryPrediction": {
                            "type": "object",
                            "properties": {
                              "text": {
                                "type": "object",
                                "properties": {
                                  "text": {
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
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "placesKey": {
        "type": "apiKey",
        "in": "header",
        "name": "X-Goog-Api-Key"
      },
      "placesOAuth": {
        "type": "oauth2",
        "flows": {
          "authorizationCode": {
            "authorizationUrl": "https://accounts.google.com/o/oauth2/v2/auth",
            "tokenUrl": "https://oauth2.googleapis.com/token",
            "scopes": {
              "https://www.googleapis.com/auth/maps-platform.places.autocomplete": "Autocomplete",
              "https://www.googleapis.com/auth/maps-platform.places": "Places API",
              "https://www.googleapis.com/auth/cloud-platform": "Cloud Platform"
            }
          }
        }
      }
    }
  },
  "security": [
    {
      "placesKey": []
    },
    {
      "placesOAuth": [
        "https://www.googleapis.com/auth/maps-platform.places.autocomplete"
      ]
    },
    {
      "placesOAuth": [
        "https://www.googleapis.com/auth/maps-platform.places"
      ]
    },
    {
      "placesOAuth": [
        "https://www.googleapis.com/auth/cloud-platform"
      ]
    }
  ],
  "x-apitools-overlay": {
    "provider_id": "google-maps-platform",
    "overlay_id": "google-maps-places-autocomplete-overlay",
    "official_openapi": false,
    "derived_from_docs": true,
    "source_refs": [
      "https://developers.google.com/maps/documentation/places/web-service/place-autocomplete",
      "https://developers.google.com/maps/documentation/places/web-service/reference/rest/v1/places/autocomplete",
      "https://developers.google.com/maps/documentation/places/web-service/get-api-key?setupProd=configure",
      "https://developers.google.com/identity/protocols/oauth2/web-server"
    ],
    "source_note": "Documents the default text prediction request, select optional inputs, prediction identifiers and display text. The full Google Discovery document is the primary source."
  }
}
`

func main() {
	const path = "catalog-openapi-cache/advisory-overlays/google-maps-places-autocomplete-overlay.json"
	if err := os.WriteFile(path, []byte(overlay), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
