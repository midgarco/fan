package worker

type Details struct {
	Id      int64
	Payload interface{}
}

type Result struct {
	Id      int64
	Payload interface{}
	Error   error
}
