package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var pipeline Out
	connector := func(in In, done In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				select {
				case <-done:
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				}
			}
		}()
		return out
	}

	pipeline = in
	for i := range stages {
		pipeline = stages[i](pipeline)
		pipeline = connector(pipeline, done)
	}

	return pipeline
}
