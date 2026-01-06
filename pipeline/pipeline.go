package pipeline

import (
	"context"
	"mime/multipart"
	"sync"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/xerrors"
)

type workerParams struct {
	stage     int
	inChan    <-chan Payload
	outChan   chan<- Payload
	errorChan chan<- error
}

func (wp *workerParams) StageIndex() int {
	return wp.stage
}

func (wp *workerParams) StageInput() <-chan Payload {
	return wp.inChan
}

func (wp *workerParams) StageOutput() chan<- Payload {
	return wp.outChan
}

func (wp *workerParams) Error() chan<- error {
	return wp.errorChan
}

type Pipeline struct {
	Stages []StageRunner
}

func NewPipeline(stages ...StageRunner) *Pipeline {
	return &Pipeline{
		Stages: stages,
	}
}

func (p *Pipeline) Process(ctx context.Context, source Source, sink Sink) error {
	var wg sync.WaitGroup
	pCtx,ctxCancelFin := contextWithCancel(ctx)
	stageChan := make([]chan Payload, len(p.Stages)+1)
	errorChan := make(chan error, len(p.Stages)+2)
	for i := 0; i < len(p.Stages); i++ {
		stageChan[i] := make(chan Payload)
	}

	for i := 0; i< len(p.Stages); i++ {
		wg.Add(1)
		go func(stageIndex int) {
			defer wg.Done()
			p.Stages[stageIndex].Run(pCtx, &workerParams{
				stage:    stageIndex,
				inChan:    stageChan[stageIndex],
				outChan :   stageChan[stageIndex+1],
				errorChan: errorChan,
			})

			close(stageChan[stageIndex+1])
			wg.Done()
		}(i)

		wg.Add(2)


		wg.Add(2)

		go func() {
			defer wg.Done()
			sourceWorker(pCtx,source, stageChan[i], errorChan )
			close(stageChan[0])
		}()

		go func(){
			defer wg.Done()
			sinkWorker(pCtx, sink, stageChan[len(p.Stages)], errorChan)
		}()	


		go func() {
			wg.Wait()
			ctxCancelFin()
			close(errorChan)
		}()

		var err error

		for pErr := range errorChan {
			err = multierror.Append(err, pErr)
			ctxCancelFin()
		}

		return err
	}


}



func sourceWorker(ctx context.Context, source Source, outChan chan<- Payload, errorChan chan<- error, ){

	for source.Next(ctx) {
		payload := source.Payload()
		select {
		case outChan <- payload:
		case <-ctx.Done():
			return
		}
	}

	if err := source.Error(); err != nil{
		wrappedError := xerrors.Errorf("pipeline error: %w", err)
		maybeEmitError(wrappedError, errorChan)
	}
}


func sinkWorker(ctx context.Context, sink Sink, inChan <- chan Payload, errChan chan <- error){
	for {
		select {
		case payload, ok :=  <- inChan:
			if !ok {
				return
			}
			if err := sink.Consume(ctx, payload); err != nil {
				wrappedError := xerrors.Errorf("pipeline error: %w", err)
				maybeEmitError(wrappedError, errChan)
				return	
			}
			payload.MarkAsProcessed()
		case <- ctx.Done():
			return	
		}
	}
}

func maybeEmitError(err error, errChan chan <- error){
	select {
		case errChan <- err:
		default:
	}
}