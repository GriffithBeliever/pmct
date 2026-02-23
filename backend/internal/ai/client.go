package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const defaultModel = "claude-sonnet-4-6"

// Client wraps the Anthropic SDK for media tracker AI features.
type Client struct {
	client *anthropic.Client
	model  string
}

// NewClient creates a new Anthropic AI client.
func NewClient(apiKey string) *Client {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Client{client: c, model: defaultModel}
}

// Complete sends a prompt and returns the full text response.
func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(c.model)),
		MaxTokens: anthropic.F(int64(2048)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		return "", fmt.Errorf("anthropic complete: %w", err)
	}
	if len(msg.Content) == 0 {
		return "", fmt.Errorf("empty response from anthropic")
	}

	var result string
	for _, block := range msg.Content {
		if block.Type == "text" {
			result += block.Text
		}
	}
	return result, nil
}

// StreamComplete streams a response, sending each token to the out channel.
func (c *Client) StreamComplete(ctx context.Context, prompt string, out chan<- string) error {
	stream := c.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(c.model)),
		MaxTokens: anthropic.F(int64(4096)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})

	for stream.Next() {
		event := stream.Current()
		if delta, ok := event.AsUnion().(anthropic.ContentBlockDeltaEvent); ok {
			if text, ok := delta.Delta.AsUnion().(anthropic.TextDelta); ok {
				select {
				case out <- text.Text:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}
	return nil
}
