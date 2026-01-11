package cralwer

import (
	"context"
	"io"
	"net/url"

	"github.com/m0hh/Learn/crawler/pipeline"
)



type linkFetcher struct {
	urlGetter URLGetter
	netDetector PrivateNetwrokDetector
}


func NewLinkFetcher(urlGetter URLGetter, netDetector PrivateNetwrokDetector) *linkFetcher {
	return &linkFetcher{
		urlGetter: urlGetter,
		netDetector: netDetector,
	}
}


func(lf *linkFetcher) Process(ctx context.Context, payload pipeline.Payload) (pipeline.Payload, error) {
	cralwerPayload, ok := payload.(*cralwerPayload)
	if !ok {
		return nil, nil
	}

	if exclusionRegex.MatchString(cralwerPayload.Url){
		return nil,nil
	}

	if isPrivate, err := lf.isPrivate(cralwerPayload.Url); err != nil || isPrivate {

		return nil, nil
	}

	resp, err := lf.urlGetter.GetURL(cralwerPayload.Url)

	if err != nil {
		return nil, err
	}

	_,errr := io.Copy(&cralwerPayload.RawContent, resp.Body)
	_= resp.Body.Close()
	
	if errr != nil {
		return nil, errr
	}

	if resp.StatusCode != 200 || resp.Status > 299 {
		return nil, nil
	}

	return cralwerPayload, nil
}



func (lf *linkFetcher) isPrivate(rawUrl string) (bool,error) {

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return false, err
	}

	return lf.netDetector.IsPrivate( parsedUrl.Hostname())
}
