package main

import (
	"fmt"
	"reflect"
)

// Metricer interface
type Metricer interface {
	Type() string
	Get() any
}

func printType(metric Metricer) {
	fmt.Println(metric.Type())
}

func printValue(metric Metricer) {
	fmt.Println(metric.Get())
}

// Gauge
type Gauge float64
type Counter int64

func (g Gauge) Type() string {
	return reflect.TypeOf(g).Name()
}

func (c Counter) Type() string {
	 return reflect.TypeOf(c).Name()
}

func main() {
	//var test = Gauge("33.4") //counter := 1

	//fmt.Println(test)

	//array := map[string]Metricer{}
	//
	//array["1"] = Gauge(1)
	//array["2"] = Counter(counter)
	//array["3"] = Counter(5)

	//for k, v := range array{
	//	fmt.Println(k)
	//	printType(v)
	//	printValue(v)
	//	fmt.Println("---")
	//}
}
