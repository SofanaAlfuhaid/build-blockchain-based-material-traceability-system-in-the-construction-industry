package main

import(
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
var _mainLogger = shim.NewLogger("TraceabilitySmartContract")

//SmartContract represents the main entart contract

type SmartContract struct {
	traceability *Trace
}

// Init initalizes the chaincode 
func (sc *SmartContract) Init(stub shim.ChaincodeStubInterface)  pb.Response {
    _mainLogger.Infof("Inside the init method")
	sc.traceability = new(Trace)
	return shim.Success(nil) 
}

// Invoke is the entry point for any transaction

func  (sc *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	var response pb.Response
	action, args := stub.GetFunctionAndParameters()
	switch action {
		case "purchaseOrder":
			response = sc.traceability.PurchaseOrder(stub)
		case "supplierRecOrderSts":
			response = sc.traceability.supplierRecOrderSts(stub)
		case "createOrderBySupplier":
			response = sc.traceability.createOrderBySupplier(stub)	
		case "logisticApproval":
			response = sc.traceability.logisticApproval(stub)	
		case "inventoryManagerReceipt":
			response = sc.traceability.InventoryManagerReceipt(stub)
		case "inventoryApproval":
			response = sc.traceability.inventoryApproval(stub)	
		case "formenConsumption":
			response = sc.traceability.FormenConsumption(stub)	
		case "stockRelease":
	        	response = sc.traceability.stockRelease(stub)
		case "consumptionApproval":
			response = sc.traceability.consumptionApproval(stub)																								
		case "displayOrderStatus":
			response = sc.traceability.displayOrderStatus(stub)
		case "consumptionApprovalForPouring":
			response = sc.traceability.consumptionApprovalForPouring(stub)
		case "materialQuery":
			response = sc.traceability.materialQuery(stub, args)						
		default:
			response = shim.Error("Invalid function name provided") 
	}
	return response	
}


func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
       _mainLogger.Criticalf("Error starting the chaincode: %v", err)
	}else {
		_mainLogger.Info("|| STARTING TRACEABILITY CHAINCODE ||")
	}
}
