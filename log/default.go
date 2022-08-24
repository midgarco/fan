package log

import "fmt"

type DefaultLogger struct{}

func (*DefaultLogger) Debug(msg string)                         { fmt.Println(msg) }
func (*DefaultLogger) Info(msg string)                          { fmt.Println(msg) }
func (*DefaultLogger) Warn(msg string)                          { fmt.Println(msg) }
func (*DefaultLogger) Error(msg string)                         { fmt.Println(msg) }
func (*DefaultLogger) Fatal(msg string)                         { fmt.Println(msg) }
func (*DefaultLogger) Debugf(msg string, params ...interface{}) { fmt.Printf(msg, params...) }
func (*DefaultLogger) Infof(msg string, params ...interface{})  { fmt.Printf(msg, params...) }
func (*DefaultLogger) Warnf(msg string, params ...interface{})  { fmt.Printf(msg, params...) }
func (*DefaultLogger) Errorf(msg string, params ...interface{}) { fmt.Printf(msg, params...) }
func (*DefaultLogger) Fatalf(msg string, params ...interface{}) { fmt.Printf(msg, params...) }
