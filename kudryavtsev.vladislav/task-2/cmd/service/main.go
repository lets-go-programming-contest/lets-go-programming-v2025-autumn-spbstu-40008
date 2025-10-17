package main

import (
	"fmt"
)

func main() {
	var (
		n  int
		k  int
		ud = 15
		uh = 30
		tw int
		ts string
	)

	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}

	for i := 0; i < n; i++ {
		_, err = fmt.Scan(&k)
		if err != nil {
			return
		}
		for j := 0; j < k; j++ {
			_, err = fmt.Scan(&ts)
			if err != nil {
				return
			}

			_, err = fmt.Scan(&tw)
			if err != nil {
				return
			}

			switch ts {
			case "<=":
				if ud != -1 {
					if uh >= tw {
						uh = tw
					}
					if uh < ud {
						ud = -1
					}
				}
			case ">=":
				if ud != -1 {
					if ud <= tw {
						ud = tw
					}
					if ud > uh {
						ud = -1
					}
				}
			}
			fmt.Println(ud)
		}
		ud = 15
		uh = 30
	}
}
