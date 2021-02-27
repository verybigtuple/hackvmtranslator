package main

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	queue := resPriotityQueue{
		&trResult{"YFile.vm", nil},
		&trResult{"ZFile.vm", nil},
	}
	heap.Init(&queue)
	heap.Push(&queue, &trResult{bootstrap, nil})
	heap.Push(&queue, &trResult{mainf, nil})
	heap.Push(&queue, &trResult{"XFile.vm", nil})

	want := [...]string{bootstrap, mainf, "XFile.vm", "YFile.vm", "ZFile.vm"}

	for i := 0; i < len(want); i++ {
		actual := heap.Pop(&queue).(*trResult)
		if actual.Name != want[i] {
			t.Errorf("%d item: %s; want: %s", i, actual.Name, want[i])
			return
		}
	}
}
