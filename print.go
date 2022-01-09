package main

import "fmt"

func print(material []Material, idArray *map[int]entity) {
	for _, v := range material {
		fmt.Println((*idArray)[v.TypeID].Name.En, "______", v.Quantity)
		fmt.Println("volume:", (*idArray)[v.TypeID].Volume*float32(v.Quantity))
	}
}