package cralwer

import (
	"context"
	"net/http"
	"time"

	"github.com/blevesearch/bleve/v2/document"
	"github.com/m0hh/Learn/crawler/linkgraph/graph"
	"github.com/m0hh/Learn/crawler/pipeline"
	"github.com/m0hh/Learn/crawler/textindexer/index"
	"modernc.org/libc/uuid"
)


type URLGetter interface {
	GetURL(url string) (*http.Response, error)
}

type PrivateNetwrokDetector interface {
	IsPrivate(host string) (bool,error)
}

type Graph interface {
	UpsertLinks(link *graph.Link) error
	UpsertEdges(edge *graph.Edge) error
	RemoveStaleEdges(from uuid.UUID, updatedBefore time.Time) error
}



type Indexer interface {
	Index(doc *index.Document) error
}

type config struct {
	URLGetter              URLGetter
	PrivateNetworkDetector PrivateNetwrokDetector
	Graph                  Graph
	Indexer                Indexer
	FetchWorkers           int
}

type Cralwer struct {
	p *pipeline.Pipeline
}

func NewCrawler(cfg config) *Cralwer {
	return &Cralwer{
		p: assembleCralwerPipeline(cfg),
	}

}



func assembleCralwerPipeline(cfg config) *pipeline.Pipeline {
	return pipeline.NewPipeline(
		pipeline.NewFixedWorkerPool(
			NewLinkFetcher(cfg.URLGetter,cfg.PrivateNetworkDetector),		
			cfg.FetchWorkers,
		),
		pipeline.FIFO(newLinkExtractor(cfg.PrivateNetworkDetector)),
		pipeline.FIFO(newTextExtractor()),
		pipeline.NewBroadCast(
			newGraphUpdater(cfg.Graph),
			newTextIndexer(cfg.Indexer),
		),
	)
}


