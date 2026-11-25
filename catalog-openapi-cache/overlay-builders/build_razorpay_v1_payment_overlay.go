//go:build ignore

package main

import (
	"fmt"
	"os"
)

const overlay = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Razorpay Payments v1 read advisory overlay",
    "version": "v1",
    "description": "Documents a merchant payment lookup and selected payment status and amount fields."
  },
  "servers": [
    {
      "url": "https://api.razorpay.com"
    }
  ],
  "paths": {
    "/v1/payments/{id}": {
      "get": {
        "operationId": "fetchRazorpayPayment",
        "summary": "Fetch a payment using its ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "schema": {
              "type": "string"
            },
            "required": true,
            "description": "Identifier of the payment."
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
                    "id": {
                      "type": "string"
                    },
                    "amount": {
                      "type": "integer",
                      "description": "Payment amount in currency subunits."
                    },
                    "currency": {
                      "type": "string"
                    },
                    "status": {
                      "type": "string"
                    },
                    "captured": {
                      "type": "boolean"
                    },
                    "created_at": {
                      "type": "integer",
                      "format": "int64",
                      "description": "Unix timestamp."
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
      "razorpayBasic": {
        "type": "http",
        "scheme": "basic",
        "description": "Merchant key ID and key secret as Basic credentials."
      }
    }
  },
  "security": [
    {
      "razorpayBasic": []
    }
  ],
  "x-apitools-overlay": {
    "provider_id": "razorpay",
    "overlay_id": "razorpay-v1-payment-overlay",
    "official_openapi": false,
    "derived_from_docs": true,
    "source_refs": [
      "https://razorpay.com/docs/api/payments/fetch-with-id",
      "https://razorpay.com/docs/api/authentication",
      "https://razorpay.com/terms/"
    ],
    "source_note": "Documents a merchant payment lookup and selected payment status and amount fields."
  }
}
`

func main() {
	const path = "catalog-openapi-cache/advisory-overlays/razorpay-v1-payment-overlay.json"
	if err := os.WriteFile(path, []byte(overlay), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
