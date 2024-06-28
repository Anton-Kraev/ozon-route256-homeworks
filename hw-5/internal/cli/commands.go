package cli

import "slices"

type (
	commandName string

	commandInfo struct {
		name  commandName
		desc  string
		usage string
	}
)

const (
	help       commandName = "help"
	exit       commandName = "exit"
	numWorkers commandName = "nworkers"

	receiveOrder  commandName = "receive"
	returnOrder   commandName = "return"
	deliverOrders commandName = "deliver"
	clientOrders  commandName = "olist"
	refundOrder   commandName = "refund"
	refundsList   commandName = "rlist"

	addWrap commandName = "addwrap"
)

var (
	controlCommandsNames = []commandName{help, exit, numWorkers}
	controlCommandsInfo  = []commandInfo{
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
			"nworkers --num=4",
		},
	}

	orderCommandsNames = []commandName{receiveOrder, returnOrder, deliverOrders, clientOrders, refundOrder, refundsList}
	orderCommandsInfo  = []commandInfo{
		{
			receiveOrder,
			"receive order from a courier",
			"receive --clientID=123 --orderID=456 --weight=100 --cost=100 --storedUntil=dd.mm.yyyy-hh:mm:ss [--wrap=box]",
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

	wrapCommandsNames = []commandName{addWrap}
	wrapCommandsInfo  = []commandInfo{
		{
			addWrap,
			"add a new available order wrap",
			"addwrap --name=box --weight=1000 --cost=10",
		},
	}

	allCommandsInfo = slices.Concat(controlCommandsInfo, orderCommandsInfo, wrapCommandsInfo)
)
