package main

import (
 "container/heap"
 "fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
 val, ok := x.(int)
 if !ok {
  return
 }

 *h = append(*h, val)
}

func (h *IntHeap) Pop() any {
 old := *h
 n := len(old)
 x := old[n-1]
 *h = old[:n-1]

 return x
}

func main() {
 var (
  numElements int
  value       int
  kthIndex    int
  minHeap     = &IntHeap{}
  sortedHeap  []int
 )

 _, err := fmt.Scan(&numElements)
 if err != nil {
  return
 }

 for i := 0; i < numElements; i++ {
  _, err = fmt.Scan(&value)
  if err != nil {
   return
  }

  heap.Push(minHeap, value)
 }

 _, err = fmt.Scan(&kthIndex)
 if err != nil {
  return
 }

 for minHeap.Len() > 0 {
  item := heap.Pop(minHeap)

  val, ok := item.(int)
  
  if !ok {
   return
  }

  sortedHeap = append(sortedHeap, val)
 }

 fmt.Println(sortedHeap[len(sortedHeap)-kthIndex])
}