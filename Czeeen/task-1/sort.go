package mySort

func BubbleSort(ar []int) {
    for i := 0; i < len(ar); i++ {
        for j := len(ar) - 1; j > i; j-- {
            if ar[j-1] > ar[j] {
                ar[j-1], ar[j] = ar[j], ar[j-1]
            }
        }
    }
}
