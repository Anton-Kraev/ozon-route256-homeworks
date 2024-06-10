package cli

type command struct {
	name  string
	desc  string
	usage string
}

const (
	help          = "help"
	exit          = "exit"
	numWorkers    = "numWorkers"
	receiveOrder  = "receive"
	returnOrder   = "return"
	deliverOrders = "deliver"
	clientOrders  = "olist"
	refundOrder   = "refund"
	refundsList   = "rlist"
)

var commandsList = []command{
	{
		help,
		"print a list of available commands with their description",
		"help",
	},
	{
		exit,
		"shutdown the application",
		"exit",
	},
	{
		numWorkers,
		"change the number of workers in worker pool",
		"numWorkers --num=4",
	},
	{
		receiveOrder,
		"receive order from a courier",
		"receive --clientID=123 --orderID=456 --storedUntil=dd.mm.yyyy-hh:mm:ss",
	},
	{
		returnOrder,
		"return order to the courier",
		"return --orderID=123",
	},
	{
		deliverOrders,
		"deliver orders to the client",
		"deliver --orders=1,2,3,4,5",
	},
	{
		clientOrders,
		"print list of client orders",
		"olist --clientID=123 [--lastN=10] [--inStorage]",
	},
	{
		refundOrder,
		"accept a refund from client",
		"refund --clientID=123 --orderID=456",
	},
	{
		refundsList,
		"print list of refunds",
		"rlist [--pageN=3] [--perPage=10]",
	},
}
