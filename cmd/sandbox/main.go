package main
//
//import (
//	"fmt"
//	"reflect"
//)
//
//type metric struct {
//	mKey string
//	mGauge Gauge
//	mCounter Gauge
//}
//
//type metricValue interface {
//	Type() string
//	SetValue() metricValue
//}
//type gauge float64
//type counter int64
//
//func (g gauge) Type() string {
//	return reflect.TypeOf(g).Name()
//}
//
//func (c counter) Type() string {
//	return reflect.TypeOf(c).Name()
//}
//
//func (g gauge) SetValue(value metricValue) metricValue {
//	 return value
//}
//
//func (c counter) SetValue(value metricValue) metricValue {
//	return c + value
//}
//
//func main() {
//	variable := counter(63)
//
//	fmt.Println(variable.Type())
//	fmt.Println(variable)
//
//	variable.SetValue(12)
//
//	fmt.Println(variable)
//
//	variable.SetValue(12)
//
//	fmt.Println(variable)
//
//	metric := metric{
//		mKey: "test",
//		mValue: variable,
//	}
//
//	fmt.Println(metric)
//}
