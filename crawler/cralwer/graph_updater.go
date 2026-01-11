package cralwer

import (
	"context"
	"time"

	"github.com/m0hh/Learn/crawler/linkgraph/graph"
	"github.com/m0hh/Learn/crawler/pipeline"
)

type graphUpdater struct {
	updater Graph
}

func newGraphUpdater(updater Graph) *graphUpdater {
	return &graphUpdater{
		updater: updater,
	}
}

func (u *graphUpdater) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {

	payload := p.(*cralwerPayload)

	src := &graph.Link{
		ID:          payload.LinkId,
		URL:         payload.Url,
		RetrievedAt: time.Now(),
	}

	if err := u.updater.UpsertLinks(src); err != nil {
		return nil, err
	}

	for _, dstLinks := range payload.NoFollowLinks {
		dst := &graph.Link{
			URL: dstLinks,
		}
		if err := u.updater.UpsertLinks(dst); err != nil {
			return nil, err
		}

	}

	removeOldEdgesTime := time.Now()
	for _, dstLinlk := range payload.Links {
		dst := &graph.Link{
			URL: dstLinlk,
		}
		if err := u.updater.UpsertLinks(dst); err != nil {
			return nil, err
		}
		if err := u.updater.UpsertEdges(&graph.Edge{From: src.ID, To: dst.ID}); err != nil {
			return nil, err
		}
	}

	if err := u.updater.RemoveStaleEdges(src.ID, removeOldEdgesTime); err != nil {
		return nil, err
	}
	return p, nil
}
