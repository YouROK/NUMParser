package tasker

import (
	"NUMParser/utils"
	"log"
	"sync"
)

type Tasker struct {
	tasks   []func()
	threads int
	wa      sync.WaitGroup
	mu      sync.Mutex
}

func New(threads int) *Tasker {
	t := new(Tasker)
	t.threads = threads
	return t
}

func (t *Tasker) Add(wrk func()) {
	t.mu.Lock()
	t.tasks = append(t.tasks, wrk)
	t.wa.Add(1)
	t.mu.Unlock()
}

func (t *Tasker) Run() {
	if t.threads > 0 {
		utils.PForLim(t.tasks, t.threads, func(i int, fn func()) {
			log.Println("Task", i+1, "/", len(t.tasks))
			fn()
			t.wa.Done()
		})
	} else {
		for i := range t.tasks {
			t.tasks[i]()
			t.wa.Done()
		}
	}
}

func (t *Tasker) Wait() {
	if len(t.tasks) > 0 {
		t.wa.Wait()
	}
}
