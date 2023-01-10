package ziva

type Options struct {

	// goroutine num
	Num int

	// Queue all task
	Queue TodoQueue

	CreateQueue CreateQueue

	StartFunc CallbackFunc

	SucceedFunc CallbackFunc

	RetryFunc CallbackFunc

	FailedFunc CallbackFunc

	CompleteFunc CallbackFunc

	// ProxyIP http or https proxy ip
	ProxyIP string

	// ProxyLib proxy ips
	ProxyLib *ProxyLib

	SheepTime int

	TimeOut int

	IsDebug bool
}
