package main

import (
  "container/heap"
  "fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] } 
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
  *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
  old := *h
  n := len(old)
  x := old[n-1]
  *h = old[0 : n-1]
  return x
}



func main(){
  rating := &IntHeap{}
  heap.Init(rating)

  var (
    nDishes int
    nRating int
    temp int
  )

  _, err := fmt.Scan(&nDishes)
  if err != nil {
    fmt.Println("Invalid number")
    return
  }

  for range nDishes{
      _, err = fmt.Scan(&temp)
      if err != nil{
        fmt.Println("Invalid number")
        return
      }

      heap.Push(rating, temp)
  }
  
  _, err = fmt.Scan(&nRating)
  if err != nil {
    fmt.Println("Invalid number")
    return
  }

  for i := 1; i < nRating; i++ {
    heap.Pop(rating)
  }

  fmt.Println(heap.Pop(rating))
}
