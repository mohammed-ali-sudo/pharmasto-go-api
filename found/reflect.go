package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	Xinit := 50
	vinit := reflect.ValueOf(&Xinit).Elem() // addressable/settable

	printInfo(vinit)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		<-ticker.C
		if vinit.CanSet() && vinit.Kind() == reflect.Int {
			vinit.SetInt(vinit.Int() + 7)
		}
		printInfo(vinit)
	}
}

func printInfo(v reflect.Value) {
	fmt.Printf("value=%v type=%s\n", v.Interface(), v.Type())
}
