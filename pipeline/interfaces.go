package pipeline

import "context"

type Payload interface{
	Clone() Payload	
	MarkAsProcessed()
}

type Processor interface {
	Process(context.Context, Payload) (Payload, error)
}

type ProcessFunc func(context.Context,Payload)(Payload,error)

func (f ProcessFunc) Process(ctx context.Context, p Payload) (Payload, error) {
	return f(ctx,p)
}


type StageParams interface{
	StageIndex() int

	StageInput()  <-chan Payload 

	StageOutput() chan<- Payload

	Error() chan<- error
}


type StageRunner interface {
	Run(context.Context, StageParams)
}

type Source interface{
	Next(context.Context) bool
	Payload() Payload	
	Error() error
}

type Sink interface{
	Consume(context.Context, Payload) error
}