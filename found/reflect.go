package main

import (
	"fmt"
	"reflect"
)

// struct with reflection

type Person struct {
	Name string
	Age  int
}

func main() {

	// refelct

	x := 20

	rv := reflect.ValueOf(x)
	rt := rv.Type()

	fmt.Println("value", rv)
	fmt.Println("type", rt)
	fmt.Println("kind", rt.Kind())
	fmt.Println("kind", rt.Kind() == reflect.Int)
	fmt.Println("kind", rt.Kind() == reflect.String)
	fmt.Println("type", rv.IsZero())

	var x2 int = 20
	rve2 := reflect.ValueOf(&x2).Elem()
	rv2 := reflect.ValueOf(&x2)
	rt2 := rv2.Type()

	fmt.Println("original value", x2)
	rve2.SetInt(1)
	fmt.Println("modified value", x2)

	fmt.Println("modified value", rt2)

	// struct with reflection

	p := Person{
		Name: "mohmmed",
		Age:  21,
	}

	reflectv := reflect.ValueOf(p)

	for i := 0; i < reflectv.NumField(); i++ {
		fmt.Println("field", i, ":", reflectv.Field(i))
	}

	v12 := reflect.ValueOf(&p).Elem()
	nameField := v12.FieldByName("Name")

	if nameField.CanSet() {
		nameField.SetString("alawy")
	}

	fmt.Println("newName", p.Name)

	// method with reflect

	



}
