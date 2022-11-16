package ziva

type Options struct {
	Num int

	Queue TodoQueue

	CreateQueue CreateQueue

	StartFunc CallbackFunc

	SucceedFunc CallbackFunc

	RetryFunc CallbackFunc

	FailedFunc CallbackFunc

	CompleteFunc CallbackFunc

	ProxyIP string

	SheepTime int

	TimeOut int

	IsDebug bool
}
