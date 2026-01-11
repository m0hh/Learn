package cralwer

import (
	"html"
	"regexp"
	"strings"
	"sync"

	"github.com/microcosm-cc/bluemonday"
)


var (
	titleRegex         = regexp.MustCompile(`(?i)<title.*?>(.*?)</title>`)
	repeatedSpaceRegex = regexp.MustCompile(`\s+`)
)


type textExtractor struct {
	policyPool sync.Pool
}


func newTextExtractor() *textExtractor {
	return &textExtractor{
		policyPool: sync.Pool{
			New: func() interface{} {
				return bluemonday.StrictPolicy()
			},
		},
	}
}

func (te *textExtractor) Process(ctx context.Context, payload pipeline.Payload) (pipeline.Payload, error) {
	payload:= payload.(*cralwerPayload)
		
	policy := te.policyPool.Get().(*bluemonday.Policy)
	defer te.policyPool.Put(policy)

	if titleMatch := titleRegex.FindAllStringSubmatch(payload.RawContent.string()); len(titleMatch) == 2 {
		payload.Title = strings.TrimSpace(html.UnescapeString(repeatedSpaceRegex.ReplaceAllLiteralString(
			policy.Sanitize(titleMatch[1], " ")
		))

	payload.TextContent = strings.TrimSpace(html.UnescapeString(repeatedSpaceRegex.ReplaceAllLiteralString(
		policy.Sanitize(payload.RawContent.string(), " ")
	)))

	return payload, nil
}


