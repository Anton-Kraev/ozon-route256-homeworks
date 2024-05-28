package cli

type commandArg struct {
	name string
	desc string
}

type command struct {
	name string
	desc string
	args []commandArg
}

const (
	help          = "help"
	receiveOrder  = "receive_order"
	returnOrder   = "return_order"
	deliverOrders = "deliver_orders"
	clientOrders  = "client_orders"
	refundOrder   = "refund_order"
	refundsList   = "refunds_list"
)

// TODO: change commands description and args
var commandsList = []command{
	{help, "", []commandArg{}},
	{receiveOrder, "", []commandArg{}},
	{returnOrder, "", []commandArg{}},
	{deliverOrders, "", []commandArg{}},
	{clientOrders, "", []commandArg{}},
	{refundOrder, "", []commandArg{}},
	{refundsList, "", []commandArg{}},
}
