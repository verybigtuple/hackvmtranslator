package main

import "strings"

const (
	bootstrap = "!!bootstrap"
	mainf     = "!mainf"
)

type trResult struct {
	Name    string
	Builder *strings.Builder
}

type resPriotityQueue []*trResult

// sort.Interface

func (pq resPriotityQueue) Len() int {
	return len(pq)
}

func (pq resPriotityQueue) Less(i, j int) bool {
	return pq[i].Name < pq[j].Name
}

func (pq resPriotityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// heap.Interface

func (pq *resPriotityQueue) Push(x interface{}) {
	item := x.(*trResult)
	*pq = append(*pq, item)
}

func (pq *resPriotityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}
