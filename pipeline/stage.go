package pipeline

import (
	"context"
	"sync"
)



type fifo struct {
	processor Processor
}

func FIFO(proc Processor) StageRunner {
	return &fifo{
		processor: proc,
	}
}



func( r fifo) Run(ctx context.Context, params StageParams){
	for {
		select {
		case <-ctx.Done():
			return
		case payloadIn,ok := <- params.StageInput():
			if !ok{
				return
			}
			payloadOut, err := r.processor.Process(ctx, payloadIn)

			if err != nil {
				wrappedError := xerrors.Errorf("pipeline stage %d: %w",params.StageIndex(), err)
				maybeEmitError(wrappedError, params.Error())
				return
			}

			if payloadOut == nil {
				payloadIn.MarkAsProcessed()
				continue
			}

			select {
			case <-ctx.Done():
				return
			case params.StageOutput() <- payloadOut:
			}
			
		}
	}
}



type fixedWorkerPool struct {
	fifos []StageRunner
}


func NewFixedWorkerPool (proc Processor, numWorkers int) StageRunner {

	if numWorkers <=0 {
		panic("Number of workers cannot be 0 or less")
	}

	fifos := make([]StageRunner,numWorkers)

	for i:=0;i<numWorkers;i++ {
		fifos[i] = FIFO(proc)
	}

	return &fixedWorkerPool{
		fifos: fifos,
	}
}


func (p *fixedWorkerPool) Run( ctx context.Context, params StageParams) {
	var wg sync.WaitGroup

	for i :=0; i<len(p.fifos) ;i++ {
		wg.Add(1)
		go func(fifoIndex int) {
			p.fifos[fifoIndex].Run(ctx, params)
			wg.Done()
		}(i)
	}

	wg.Wait()
}


type DynamicWorkerPool struct{
	processor Processor
	tokenPool chan struct{}
}


func NewDynamicWorkerPool(proc Processor, maxWorkers int) StageRunner {
	
	if maxWorkers <=0 {
		panic("Number of workers cannot be 0 or less")
	}

	tokenPool := make(chan struct{}, maxWorkers)

	for i:=0; i< maxWorkers; i++{
		tokenPool <- struct{}{}	
	}

	return &DynamicWorkerPool{
		processor: proc,
		tokenPool: tokenPool,
	}
}

func (p *DynamicWorkerPool) Run(ctx context.Context, params StageParams) {
	stop:
		for {
			select{
			case <- ctx.Done():
				break stop

			case payLoadIn, ok := <- params.StageInput():
				if !ok{
					break stop
				}
				var token struct{}
				select {
				case <-ctx.Done():
					break stop
				case token = <- p.tokenPool:	
				}

				go func(payloadIn Payload, token struct{}) {
					defer func(){
						p.tokenPool <- token
					}()

					payloadOut,err := p.processor.Process(ctx, payloadIn)

					if err != nil {
						wrappedError := xerrors.Errorf("pipeline stage %d: %w",params.StageIndex(), err)
						maybeEmitError(wrappedError, params.Error())
						return
					}
					
					if payloadOut == nil {
						payloadIn.MarkAsProcessed()
						return
					}

					select {
						case <-ctx.Done():
						case params.StageOutput() <- payloadOut:
					}

			}(payLoadIn, token)
}
		}

		for i:= 0;i<cap(p.tokenPool);i++ {
			<- p.tokenPool
		}

}


type broadCast struct {
	runners []StageRunner
}


func NewBroadCast(processors ...Processor) StageRunner {
	runners := make([]StageRunner, len(processors))
	for i:= 0 ;i< len(processors); i++ {
		if len(processors) <=0 {
			panic("Number of workers cannot be 0 or less")
		}
		runners[i] = FIFO(processors[i])

	}
	return & broadCast{
		runners: runners,
	}

}

func (b *broadCast) Run(ctx context.Context, params StageParams) {
	var wg sync.WaitGroup
	inCh := make([]chan Payload, len(b.runners))

	for i:=0; i< len(b.runners); i++ {
		go func(runnerIndex int){
			wg.Add(1)
			inCh[runnerIndex] = make(chan Payload)
			stageParams := &workerParams{
				stage:  params.StageIndex(),
				inChan:  inCh[runnerIndex],
				outChan: params.StageOutput(),
				errorChan:     params.Error(),
			}
			b.runners[runnerIndex].Run(ctx, stageParams)
			wg.Done()
		}(i)

	}

stop:
	for {
		select {
		case <- ctx.Done():
			break stop
		case payloadIn, ok := <- params.StageInput():
			if !ok{
				break stop
			}

			for i:=0; i< len(b.runners); i++ {
				payloadClone := payloadIn.Clone()
				select {
				case <- ctx.Done():
					break stop
				case inCh[i] <- payloadClone:
				}
			}
			payloadIn.MarkAsProcessed()
		}
	}

	for i:=0; i< len(b.runners); i++ {
		close(inCh[i])
	}

	wg.Wait()
		
}