package cls

import (
	"time"

	clssdk "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

type status int

const (
	Failed status = iota
	Success
)

type Result struct {
	Status status

	RequestID    string
	ErrorCode    string
	ErrorMessage string
	Time         time.Time
}

type Handle func(results ...Result)

type callback struct {
	fn Handle
}

func (c *callback) Success(result *clssdk.Result) {
	if c.fn == nil {
		return
	}

	var results []Result
	attemptList := result.GetReservedAttempts()
	for _, attempt := range attemptList {
		var res Result
		if attempt.Success {
			res.Status = Success
		}
		res.ErrorCode = attempt.ErrorCode
		res.ErrorMessage = attempt.ErrorMessage
		res.RequestID = attempt.RequestId
		res.Time = time.Unix(attempt.TimeStampMs/1000, 0)
		results = append(results, res)
	}
	c.fn(results...)
}

func (c *callback) Fail(result *clssdk.Result) {
	if c.fn == nil {
		return
	}

	var res Result
	if result.IsSuccessful() {
		res.Status = Success
	}
	res.ErrorCode = result.GetErrorCode()
	res.ErrorMessage = result.GetErrorMessage()
	res.RequestID = result.GetRequestId()
	res.Time = time.Unix(result.GetTimeStampMs()/1000, 0)
	c.fn(res)
}
