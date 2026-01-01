package pipeline

import "_/home/codespace/Learn/pipeline"



type workerParams struct {
	stage int
	inChan <-chan Payload
	outChan chan<- Payload
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


func (p *Pipeline) Process