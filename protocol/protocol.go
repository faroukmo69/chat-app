package protocol

type JoinArgs struct {
	Username string
}

type JoinReply struct {
	Success bool
	Message string
}

type MessageArgs struct {
	Username string
	Text     string
}

type StandardReply struct {
	Success bool
	Message string
}

type PollArgs struct {
	Username string
}

type PollReply struct {
	Messages []string
}

type LeaveArgs struct {
	Username string
}
