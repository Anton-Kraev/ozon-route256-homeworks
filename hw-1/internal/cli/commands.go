package cli

type command struct {
	name    string
	desc    string
	example string
	args    []commandArg
}

type commandArg struct {
	name string
	desc string
}

const (
	help          = "help"
	exit          = "exit"
	receiveOrder  = "receive"
	returnOrder   = "return"
	deliverOrders = "deliver"
	clientOrders  = "olist"
	refundOrder   = "refund"
	refundsList   = "rlist"
)

// TODO: change commands description and args
var commandsList = []command{
	{help, "", "", []commandArg{}},
	{exit, "", "", []commandArg{}},
	{receiveOrder, "", "", []commandArg{}},
	{returnOrder, "", "", []commandArg{}},
	{deliverOrders, "", "", []commandArg{}},
	{clientOrders, "", "", []commandArg{}},
	{refundOrder, "", "", []commandArg{}},
	{refundsList, "", "", []commandArg{}},
}
