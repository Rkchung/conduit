package common

// for errors, look at package io: EOF
const (
	ProviderClientListenerAddr = ":8000"
	ClientListenerAddr         = ":8001"
	ProviderListenerAddr       = ":8002"
	LoggerAddr                 = ":8003"
	RegionalMasterListenerAddr = ":8004"
    LocalMasterListenerAddr    = ":8005"
)

type Nothing struct{}

type JobRequest struct {
	// Todo: Job Request parameters
}

func NewNothing() *Nothing {
	return new(Nothing)
}

type RequestProviderReply struct {
	Addr string
}

type RequestProviderError struct {
	Msg string
}

func (e *RequestProviderError) Error() string {
	return e.Msg
}

type RequestRegionalMasterReply struct {
	Addr string
}

type GetJobRequestFromLogReply struct {
	Req JobRequest
}

type ProviderRegisterArgs struct {
	pID string
}

type ClientJoinLeaveArgs struct {
    pID string
}

type PingArgs struct {
	Addr string
}

type ClientRegisterArgs struct {
	pID string
}

type Executable struct {
	Interpreter string
	Content     []byte
}

type ExecutionReply struct {
	Output []byte
}
