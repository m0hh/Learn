package cralwer

import (
	"context"

	"github.com/m0hh/Learn/crawler/pipeline"
	"github.com/m0hh/Learn/crawler/textindexer/index"
)


type textIndexer struct {
	indexer Indexer
}

func newTextIndexer(indexer Indexer) *textIndexer{
	return &textIndexer{
		indexer: indexer,
	}
}


func (ti *textIndexer) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*cralwerPayload)

	doc := &index.Document{
		URL:         payload.Url,
		Title:       payload.Title,
		Content:     payload.TextContent,
		LinkID:       payload.LinkId, 
	}
	
	if err := ti.indexer.Index(doc); err != nil {
		return nil, err
	}
	return p, nil
}