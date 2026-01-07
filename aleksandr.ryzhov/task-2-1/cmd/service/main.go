package main

import "fmt"

/*
Задача 1 Борьба за место под солнцем.

Постановка задачи:
В офисе компании Y у каждого из N отделов есть свой кондиционер,
на котором может быть выставлена температура в диапазоне от 15 до 30 градусов.
В каждом отделе работает K сотрудников. Каждый сотрудник, приходя в офис,
устанавливает желаемое значение температурной границы (не больше или не меньше T).
Напишите программу, которая после каждого прибывшего сотрудника
будет выводить оптимальную температуру для всего отдела.
Если такой температуры нет выведете -1.
Первым подается на вход число N - количество отделов (от 1 до 1000).
Вторым числом подается число K - количество сотрудников (от 1 до 1000).
После идут данные температуры: <= / >= *числовое значение*

Выходные данные:
В случае успешного подбора программа должны вывести числовое
значение, соответствующее оптимальной температуре.
В случае неудачи требуется вывести -1

Входные данные 	Выходные данные
2
1
>= 30			30
6
>= 18			18
<= 23			18
>= 20			20
<= 27			20
<= 21			20
>= 28			-1
*/

func main() {
	var depCount int
	var emplCount int
	var minTemp int
	var maxTemp int
	var reqTemp int
	var reqTempInfo string

	_, err := fmt.Scan(&depCount)
	if err != nil || depCount < 1 {
		fmt.Println("Incorrect number of departments", err)
		return
	}

	for range depCount {
		_, err := fmt.Scan(&emplCount)
		if err != nil || emplCount < 1 {
			fmt.Println("Incorrect number of employees", err)
			return
		}

		maxTemp, minTemp = 30, 15

		for range emplCount {
			_, err := fmt.Scan(&reqTempInfo, &reqTemp)
			if err != nil {
				fmt.Println("Incorrect temperature information", err)
				return
			}

			switch reqTempInfo {
			case "<=":
				maxTemp = min(maxTemp, reqTemp)
			case ">=":
				minTemp = max(minTemp, reqTemp)
			default:
				fmt.Println("Incorrect temperature information")
				return
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
				continue
			}

			fmt.Println(minTemp)
		}
	}
}
