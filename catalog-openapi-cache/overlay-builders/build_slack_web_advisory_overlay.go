//go:build ignore

package main

import (
	"encoding/json"
	"os"
)

func main() {
	const postMessageDocs = "https://docs.slack.dev/reference/methods/chat.postMessage/"
	const historyDocs = "https://docs.slack.dev/reference/methods/conversations.history/"
	const oauthDocs = "https://docs.slack.dev/authentication/installing-with-oauth/"
	text := map[string]any{"type": "string"}
	boolean := map[string]any{"type": "boolean"}
	content := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title": "Slack Web API Advisory Overlay", "version": "1.0.0",
			"description": "Reviewed JSON message posting and public-channel history from Slack's current documentation.",
		},
		"servers": []map[string]any{{"url": "https://slack.com/api"}},
		"x-apitools-overlay": map[string]any{
			"provider_id": "slack", "overlay_id": "slack-web-message-history-overlay",
			"official_openapi": false, "derived_from_docs": true,
			"source_refs": []string{postMessageDocs, historyDocs, oauthDocs},
			"source_note": "The archived Swagger uses classic chat scopes and an aggregate history scope list. This overlay represents a small current subset.",
		},
		"components": map[string]any{"securitySchemes": map[string]any{
			"slackOAuth2": map[string]any{
				"type": "oauth2", "flows": map[string]any{
					"authorizationCode": map[string]any{
						"authorizationUrl": "https://slack.com/oauth/v2/authorize",
						"tokenUrl":         "https://slack.com/api/oauth.v2.access",
						"scopes": map[string]string{
							"chat:write":       "Send messages to participating conversations.",
							"channels:history": "Read messages in eligible public channels.",
						},
					},
				},
			},
		}},
		"paths": map[string]any{
			"/chat.postMessage": map[string]any{"post": map[string]any{
				"operationId":  "postSlackMessage",
				"summary":      "Post a text message to a participating Slack conversation",
				"externalDocs": map[string]string{"url": postMessageDocs},
				"security":     []map[string][]string{{"slackOAuth2": {"chat:write"}}},
				"requestBody": map[string]any{
					"required": true, "content": map[string]any{
						"application/json": map[string]any{"schema": map[string]any{
							"type": "object", "required": []string{"channel", "text"},
							"properties": map[string]any{"channel": text, "text": text, "thread_ts": text},
						}},
					},
				},
				"responses": response(map[string]any{"ok": boolean, "channel": text, "ts": text, "error": text}),
			}},
			"/conversations.history": map[string]any{"get": map[string]any{
				"operationId":  "getSlackPublicChannelHistory",
				"summary":      "Read messages in an eligible public Slack channel",
				"externalDocs": map[string]string{"url": historyDocs},
				"security":     []map[string][]string{{"slackOAuth2": {"channels:history"}}},
				"parameters": []map[string]any{
					{"name": "channel", "in": "query", "required": true, "schema": text},
					{"name": "cursor", "in": "query", "schema": text},
					{"name": "limit", "in": "query", "schema": map[string]any{"type": "integer", "minimum": 1}},
				},
				"responses": response(map[string]any{
					"ok": boolean, "has_more": boolean, "error": text,
					"messages":          map[string]any{"type": "array", "items": map[string]any{"type": "object", "properties": map[string]any{"ts": text, "text": text}}},
					"response_metadata": map[string]any{"type": "object", "properties": map[string]any{"next_cursor": text}},
				}),
			}},
		},
	}
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("catalog-openapi-cache/advisory-overlays/slack-web-message-history-overlay.json", append(data, '\n'), 0o644); err != nil {
		panic(err)
	}
}

func response(properties map[string]any) map[string]any {
	return map[string]any{"200": map[string]any{
		"description": "Slack Web API result. Inspect ok and error for method failures.",
		"content": map[string]any{"application/json": map[string]any{"schema": map[string]any{
			"type": "object", "required": []string{"ok"}, "properties": properties,
		}}},
	}}
}
