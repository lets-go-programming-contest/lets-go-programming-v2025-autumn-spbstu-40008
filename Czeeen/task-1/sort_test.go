	package mySort

import "testing"

func TestBubbleSort(t *testing.T) {
    ar := []int{5, 35, 12, 24, 21, 98, 6, 10}
    expected := []int{5, 6, 10, 12, 21, 24, 35, 98}

    BubbleSort(ar)

    for i, v := range ar {
        if v != expected[i] {
            t.Errorf("Expected %v, got %v", expected, ar)
            break
        }
    }
}
