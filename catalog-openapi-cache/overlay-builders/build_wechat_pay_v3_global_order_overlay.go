//go:build ignore

package main

import (
	"fmt"
	"os"
)

const overlay = `{
  "openapi": "3.0.3",
  "info": {
    "title": "WeChat Pay API v3 Global order query advisory overlay",
    "version": "v3",
    "description": "Documents the common-mode query by the WeChat Pay order number. Institutional mode and the merchant order number route are outside this subset."
  },
  "servers": [
    {
      "url": "https://apihk.mch.weixin.qq.com"
    }
  ],
  "paths": {
    "/v3/global/transactions/id/{id}": {
      "get": {
        "operationId": "queryWeChatPayOrder",
        "summary": "Query WeChat Pay order status by order ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "schema": {
              "type": "string",
              "minLength": 1,
              "maxLength": 32
            },
            "required": true,
            "description": "WeChat Pay order number."
          },
          {
            "name": "mchid",
            "in": "query",
            "schema": {
              "type": "string",
              "minLength": 1,
              "maxLength": 32
            },
            "required": true,
            "description": "Common-mode merchant ID assigned by WeChat Pay."
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
                    "mchid": {
                      "type": "string",
                      "minLength": 1,
                      "maxLength": 32
                    },
                    "id": {
                      "type": "string",
                      "minLength": 1,
                      "maxLength": 32
                    },
                    "out_trade_no": {
                      "type": "string",
                      "minLength": 1,
                      "maxLength": 32
                    },
                    "trade_state": {
                      "type": "string",
                      "enum": [
                        "SUCCESS",
                        "REFUND",
                        "NOTPAY",
                        "CLOSED",
                        "REVOKED",
                        "USERPAYING",
                        "PAYERROR"
                      ]
                    },
                    "amount": {
                      "type": "object",
                      "properties": {
                        "total": {
                          "type": "integer",
                          "description": "Total order amount in currency subunits."
                        },
                        "currency": {
                          "type": "string",
                          "minLength": 3,
                          "maxLength": 3
                        }
                      },
                      "required": [
                        "total",
                        "currency"
                      ]
                    }
                  },
                  "required": [
                    "mchid",
                    "id",
                    "out_trade_no",
                    "trade_state",
                    "amount"
                  ]
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
      "wechatPaySignature": {
        "type": "http",
        "scheme": "wechatpay2-sha256-rsa2048",
        "description": "Symbolic signed Authorization header for a WeChat Pay merchant request."
      }
    }
  },
  "security": [
    {
      "wechatPaySignature": []
    }
  ],
  "x-apitools-overlay": {
    "provider_id": "wechat-pay",
    "overlay_id": "wechat-pay-v3-global-order-overlay",
    "official_openapi": false,
    "derived_from_docs": true,
    "source_refs": [
      "https://pay.wechatpay.cn/doc/global/v3/en/4012356497",
      "https://pay.wechatpay.cn/doc/global/v3/en/4012357147",
      "https://pay.wechatpay.cn/doc/global/v3/en/4012357105"
    ],
    "source_note": "Documents the common-mode query by the WeChat Pay order number. Institutional mode and the merchant order number route are outside this subset."
  }
}
`

func main() {
	const path = "catalog-openapi-cache/advisory-overlays/wechat-pay-v3-global-order-overlay.json"
	if err := os.WriteFile(path, []byte(overlay), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
