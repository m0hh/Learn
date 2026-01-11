package cralwer

import (
	"context"
	"net/url"
	"regexp"

	"github.com/m0hh/Learn/crawler/pipeline"
)

var (
	exclusionRegex = regexp.MustCompile(`(?i)\.(?:jpg|jpeg|png|gif|ico|css|js)$`)
	baseHrefRegex  = regexp.MustCompile(`(?i)<base.*?href\s*?=\s*?"(.*?)\s*?"`)
	findLinksRegex = regexp.MustCompile(`(?i)<a.*?href\s*?=\s*?"\s*?(.*?)\s*?".*?>`)
	NoFollowRegex  = regexp.MustCompile(`(?i)rel\s*?=\s*?"?nofollow"?`)
)

type linkExtractor struct {
	netDetector PrivateNetwrokDetector
}

func newLinkExtractor(netDetector PrivateNetwrokDetector) *linkExtractor {
	return &linkExtractor{
		netDetector: netDetector,
	}
}

func (le *linkExtractor) Process(ctx context.Context, payload pipeline.Payload) (pipeline.Payload, error) {
	cralwerPayload, ok := payload.(*cralwerPayload)
	if !ok {
		return nil, nil
	}

	relTo, err := url.Parse(cralwerPayload.Url)
	if err != nil {
		return nil, err
	}

	content := cralwerPayload.RawContent.String()

	if baseMatch := baseHrefRegex.FindStringSubmatch(content); len(baseMatch) == 2 {
		if baseUrl := resolveUrl(relTo, enureHasTrailingSlash(baseMatch[1])); baseUrl != nil {
			relTo = baseUrl
		}

	}

	seenMap := make(map[string]struct{})

	for _, match := range findLinksRegex.FindAllStringSubmatch(content, -1) {
		link := resolveUrl(relTo, match[1])
		if !le.retainLink(relTo.Hostname(), link) {
			continue
		}

		link.Fragment = ""
		linkStr := link.String()

		if _, seen := seenMap[linkStr]; seen {
			continue
		}

		if exclusionRegex.MatchString(linkStr) {
			continue
		}

		seenMap[linkStr] = struct{}{}

		if NoFollowRegex.MatchString(match[0]) {
			cralwerPayload.NoFollowLinks = append(cralwerPayload.NoFollowLinks, linkStr)
		} else {
			cralwerPayload.Links = append(cralwerPayload.Links, linkStr)
		}

	}
	return cralwerPayload, nil
}

func (le *linkExtractor) retainLink(srcHost string, link *url.URL) bool {
	if link == nil {
		return false
	}

	if link.Scheme != "http" && link.Scheme != "https" {
		return false
	}

	if link.Hostname() == srcHost {
		return true
	}

	if isPrivaate, err := le.netDetector.IsPrivate(link.Hostname()); err != nil || isPrivaate {
		return false
	}

	return true
}

func enureHasTrailingSlash(s string) string {
	if s[len(s)-1] != '/' {
		s = s + "/"
	}
	return s
}

func resolveUrl(relTO *url.URL, target string) *url.URL {
	lenTarget := len(target)

	if lenTarget == 0 {
		return nil
	}

	if lenTarget >= 1 && target[0] == '/' {
		if lenTarget >= 2 && target[1] == '/' {
			target = relTO.Scheme + ":" + target
		}
	}

	if targetUrl, err := url.Parse(target); err == nil {
		return relTO.ResolveReference(targetUrl)
	}

	return nil
}
