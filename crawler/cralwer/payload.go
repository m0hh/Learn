package cralwer

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/m0hh/Learn/crawler/pipeline"
	"github.com/pborman/uuid"
)

var (
	_ pipeline.Payload = *(cralwerPayload)(nil)
	payloadPool = sync.Pool{
		New: func() interface{} {
			return new(cralwerPayload)
		},
	}
)

type cralwerPayload struct {
	LinkId uuid.UUID
	Url   string
	RetrievedAt time.Time
	RawContent bytes.Buffer
	NoFollowLinks []string
	Links []string
	Title string
	TextContent string
}


func(p *cralwerPayload) Clone() pipeline.Payload {
	newP := payloadPool.Get().(*cralwerPayload)
	newP.LinkId = p.LinkId
	newP.Url = p.Url
	newP.RetrievedAt = p.RetrievedAt
	newP.NoFollowLinks = append([]string(nil), p.NoFollowLinks...)
	newP.Links = append([]string(nil), p.Links...)
	newP.Title = p.Title
	newP.TextContent = p.TextContent
	_,err := io.Copy(&newP.RawContent,&p.RawContent)
	if err != nil {
		panic("[BUG] failed to clone payload RawContent: " + err.Error())
	}
	return newP
}

func (p *cralwerPayload) MarkAsProcessed( ) {
	p.Url = p.Url[:0]
	p.NoFollowLinks = p.NoFollowLinks[:0]
	p.Links = p.Links[:0]
	p.Title = p.Title[:0]
	p.TextContent = p.TextContent[:0]
	p.RawContent.Reset()
	payloadPool.Put(p)
}
