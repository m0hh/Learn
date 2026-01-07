package cralwer

import "regexp"


var (
	exclusionRegex = regexp.MustCompile(`(?i)\.(?:jpg|jpeg|png|gif|ico|css|js)$`)
	baseHrefRegex = regexp.MustCompile(`(?i)<base.*?href\s*?=\s*?"(.*?)\s*?"`)
	findLinksRegex = regexp.MustCompile(`(?i)<a.*?href\s*?=\s*?"\s*?(.*?)\s*?".*?>`)
	NoFollowRegex = regexp.MustCompile(`(?i)rel\s*?=\s*?"?nofollow"?`)
)

 type linkExtractor struct {
	newDetector PrivateNetwrokDetector
}

func newLinkExtractor(netDetector PrivateNetwrokDetector) *linkExtractor {
	return &linkExtractor{
		newDetector: netDetector,
	}
}


func (le *linkExtractor) Proc