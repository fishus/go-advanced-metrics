package agent

import (
	"sync"
)

var wgAgent sync.WaitGroup

func Shutdown() {
	wgAgent.Wait()
}
