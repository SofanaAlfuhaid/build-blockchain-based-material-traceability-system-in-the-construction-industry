package main

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

var  _TraceingLogger = shim.NewLogger("Traceing")

//Events

const _CreateEvent = "PO_ORDER"
const _CreateOrderEvent = "CO_ORDER"
const _SupplyRecOrderSts = "SR_ORDER"
const _LogiscticApproval = "LA_ORDER"
const _InventoryManager = "IM_ORDER"
const _InventoryApproval = "IA_ORDER"
const _Foremen = "M_FOREMEN"
const _StockRelease = "SR_MANAGER"
const _ConsumtionOrder = "CO_ORDER"
const _DisplayOrder = "DO_STATUS"
const _ApprovalForPouring = "AP_STATUS"
const _UpdateEvent = "TRU_OWNERSHIP_CHANGE"

// Trace manages all Traceing related transactions
type Trace struct {
}


type PurchaseOrder struct {
	ObjType        string `json:"obj"`
	PONumber       string `json:"po"`
	ItemNumber     string `json:"itemno"`
	ItemName       string `json:"itemname"`
	Description    string `json:"desc"`
	Quantity       float64   `json:"quan"`
	UintPrice      string `json:"uprice"`
	Amount         string `json:"amt"`
	// State          string `json:"state"`
	ShipTo         string `json:"addr"` // buyer location
	DeliveryDue    string `json:"delvry"`  //at what time buyer want the delivery
	BuyerID        string `json:"buyerid"`
	SupplierID     string `json:"suppid"`	
	Creator        string `json:"crt"`
	UpdateBy	   string `json:"uby"`
	CreateTs       string `json:"cts"`
	UpdateTs       string `json:"uts"` 
	PoStatus       string `json:"posts"`        //create 	, Inprogress, reject, Instock
	Ownership      string `json:"ownr"`
	Standard       []float64  `json:"standard"`
	
	//Create Delivery Order , 
	//on delivery Order ownership change to carrierid, 
	//PoStauts change to Inprogress
	CarrierId      string `json:"cid"`    // shipment company id  who is responsible for driveing the truck
	// Location    same as shipto
	ShipmentId     string `json:"shid"`  //delivery order number
	Truckno        string `json:"trno"`
	RegulatorId    string `json:"regid"`    
	DoStatus       string `json:"dosts"`   //pending ,expecting confirmation from regulator, shipped, dispute, arrived,
	GTIN           string `json:"gtin"`  
	
	
	// Logistics Approval

	// DoStatus change to shiped or dispute


	//Inventory Manager
	InvMngId       string `json:"invmngid"`
	ExpDate        string `json:"expdate"`
	StockLocation  string `json:"gis"` // zone1
	GoodReceipt    string `json:"grept"`
	GRStatus       string `json:"grsts"`    // pending , received, backorder


	//Inventory Aproval (Regulator)

	//DoStaus change to arrived 
	// PoStatus change to inStock
	//GRStatus change to received
	//ownership change to InventoryManagerId
	//GRStatus change to backorder in case of dispute
	Innerdia       float64 `json:"innerdia"`
	Outerdia       float64 `json:"outerdia"`
	Wallwidth      float64 `json:"wallwidth"`
	StadBatchWeght float64 `json:"stanbathweght"`	

	
	//Foremen
	ForemenUpdate  []ForemenType  `json:"foremenupdate"`
}

type ForemenType struct 	{
	PONumber        string `json:"po"`
	ForemenId       string `json:"foreid"`
	Purpuse         string `json:"purps"`
	ForDesc         string `json:"fdesc"`  //same as description
	CCOrder         string `json:"ccorder"`  // created
	CONumber        string `json:"conum"`
	PQuantity       float64   `json:"pquanty"`
	UpdateTs        string `json:"futs"`

	
	//Inventory Manager Stock release for Regulator confirmation

	//CCorder would be change to regulation pending
	BatchId         string `json:"batchid"`
	

	//Regulator consumption Approval

	BatchWeight     string `json:"bweght"`
	// on confirmation order CCOrder change to ready to use
	// Ownership change to formenId
	

	// Foremen Display Orders Status

	// CCOrder status change to expecting confirmation from regulator for pouring


	// Regulator Consumption Approval

	Density         string `json:"density"`
	//CCOrder change to ready to be poured
	//when dispute foremen should redo the work
}


var dltDomainNames = map[string]string{
	"cst.track.com":        "CS", //Construction
	"sup.track.com":        "SP", //Supplier
	"reg.track.com":        "RG", //Regulator
}

func validEnumEntry(input string, enumMap map[string]string) bool {
	if _, isEntryExists := enumMap[input]; !isEntryExists {
		return false
	}
	return true
}

func isValidDomainName(domainName string) (bool, string) {
	if !validEnumEntry(domainName, dltDomainNames) {
		return false, "WARNING: Needs to be a valid Domain name"
	}
	return true, ""
}


// Create the purchase order
func (tr *Trace) PurchaseOrder(stub shim.ChaincodeStubInterface) peer.Response {
	_TraceingLogger.Infof("Purchase Order")
	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
       return shim.Error("Invalid number of arguments provided")
	}	
	var orderToSave PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &orderToSave)
	if err != nil {
		return shim.Error("Invalid json provided as input")
	}
	// enCert, err := cid.GetX509Certificate(stub)
	// id, errrr :=cid.GetID(stub) 
	// fmt.Println("enCert",enCert.Subject.CommonName, "ssss", err)
	// fmt.Println("id....", id,"....",errrr )
	authorize, creator :=tr.getInvokerIdentity(stub)
	if authorize == false {
		return shim.Error("Unauthorized Access")

	}
	if recordBytes, _ :=stub.GetState(orderToSave.PONumber); len(recordBytes) > 0 {
        return shim.Error("Order already exist please provide unique purchase order number")
	}

	orderToSave.Creator = creator
	orderToSave.UpdateBy = creator
	orderToSave.ObjType  = "PurchaseOrder"
	orderToSave.DoStatus = "pending" 
	orderToSave.Ownership = orderToSave.SupplierID 
	if orderToSave.ItemName == "Pipe" {
		orderToSave.Standard = append(orderToSave.Standard, 3.068,3.5,0.216,1.41) 
	}else {
		orderToSave.Standard = append(orderToSave.Standard, 40)
	}
	orderJson, _:=json.Marshal(orderToSave)


	//save the purchase order

	_TraceingLogger.Infof("order.PONumber..........", orderToSave.PONumber)

	erre := stub.PutState(orderToSave.PONumber, orderJson)

	if erre !=nil {
		_TraceingLogger.Errorf("Unable to save with PONumber " + orderToSave.PONumber)
		return shim.Error("Unable to save with PONumber " + orderToSave.PONumber)
	}
	_TraceingLogger.Infof("SupplierOrderReceive : PutState Success : " + string(orderJson))
	erer := stub.SetEvent(_CreateEvent, orderJson)

	if erer != nil {
		_TraceingLogger.Errorf("Event not generated for event : PO_ORDER")
		// return shim.Error("{\"error\" : \"Unable to generate Purchase Order Event.\"}")
	}
	resultData := map[string]interface{}{
		"trxID" : stub.GetTxID(),
		"POID"  : orderToSave.PONumber,
		"message" : "Purchase Order created by the buyer successfully",
		"Order" : orderToSave,
		"status" : true,
	}

	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

func (tr *Trace) supplierRecOrderSts(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("SupplierOrderReceive")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1{
		_TraceingLogger.Errorf("SupplierOrderReceive : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var stsOrder PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &stsOrder)
	if err != nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(stsOrder.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("SupplierOrderReceive : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
    if existingOrder.PoStatus == "create" {

		existingOrder.PoStatus = stsOrder.PoStatus   //Inprogress or reject
		existingOrder.UpdateTs = stsOrder.UpdateTs
	} else {
		_TraceingLogger.Errorf("SupplierOrderReceive : purchase order is not in create state")
		return shim.Error("purchase order is not in create state")
	}

	OrderBytes, err := json.Marshal(existingOrder)
	if err !=nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : Marshalling Error : " + string(err.Error()))
		return shim.Error("SupplierOrderReceive : Marshalling Error : " + string(err.Error()))
	}
	_TraceingLogger.Infof("SupplierOrderReceive : saving the Create Order : " + existingOrder.PONumber)

	errorr :=stub.PutState(stsOrder.PONumber, OrderBytes)

	if errorr != nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : Put State Failed Error : " + string(errorr.Error())) 
		return shim.Error("Put State Failed Error : " + string(errorr.Error()))
	}

	_TraceingLogger.Infof("SupplierOrderReceive : PutState Success : " + string(OrderBytes))
	err2 := stub.SetEvent(_SupplyRecOrderSts, OrderBytes)
	if err2 != nil {
		_TraceingLogger.Errorf("SupplierOrderReceive : Event not generating for : " + _SupplyRecOrderSts)
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       stsOrder.PONumber,
		"message":       "Supplier receive the Order Successfully.",
		"Header":        stsOrder,
	}

		respJSON,_ := json.Marshal(resultData)
		return shim.Success(respJSON)
		
}
	 
func (tr *Trace) supplierOrderCorrection(stub shim.ChaincodeStubInterface) peer.Response{
	_TraceingLogger.Infof("supplierOrderCorrection")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("supplierOrderCorrection : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var logisticOrder PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &logisticOrder)
	if err != nil {
		_TraceingLogger.Errorf("supplierOrderCorrection : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(logisticOrder.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("supplierOrderCorrection : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("supplierOrderCorrection : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
    if existingOrder.GRStatus == "backorder" {

		existingOrder.PoStatus = logisticOrder.PoStatus   // inProgress
		existingOrder.UpdateTs = logisticOrder.UpdateTs
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("supplierOrderCorrection : Marshalling Error : " + string(err.Error()))
			return shim.Error("supplierOrderCorrection : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("supplierOrderCorrection : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(logisticOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("supplierOrderCorrection : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("supplierOrderCorrection : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_LogiscticApproval, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("supplierOrderCorrection : Event not generating for : " + _LogiscticApproval)
		}
	}else {
		_TraceingLogger.Errorf("supplierOrderCorrection : GRStatus is not in backorder state")
		return shim.Error("GRStatus is not in backorder state")
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       logisticOrder.PONumber,
		"message":       "Supplier Order corrected successfully.",
		"Order":        logisticOrder,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)

}
	// Creating Delivery Order by supplier 

func (tr *Trace) createOrderBySupplier(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("OrderBySupplier")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("OrderBySupplier : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var Order PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &Order)
	if err != nil {
		_TraceingLogger.Errorf("OrderBySupplier : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(Order.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("OrderBySupplier : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("OrderBySupplier : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
    if existingOrder.PoStatus == "inProgress" {

		existingOrder.CarrierId = Order.CarrierId
		existingOrder.ShipmentId = Order.ShipmentId
		existingOrder.Truckno = Order.Truckno
		existingOrder.RegulatorId = Order.RegulatorId
		existingOrder.Ownership  = Order.CarrierId  // ownership change to carrierId
		existingOrder.DoStatus = Order.DoStatus   // expecting confirmation from regulator
		existingOrder.UpdateTs = Order.UpdateTs
		existingOrder.GTIN   = Order.GTIN
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("OrderBySupplier : Marshalling Error : " + string(err.Error()))
			return shim.Error("OrderBySupplier : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("OrderBySupplier : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(Order.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("OrderBySupplier : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("OrderBySupplier : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_CreateOrderEvent, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("OrderBySupplier : Event not generating for : " + _CreateOrderEvent)
		}
	}else if existingOrder.PoStatus == "create" {
		_TraceingLogger.Errorf("OrderBySupplier :first order need to be received by the supplier")
		return shim.Error("OrderBySupplier : first order need to be received by the supplier")
	}else if existingOrder.PoStatus == "rejected"{
		_TraceingLogger.Errorf("OrderBySupplier : order is already rejected by the supplier")
		return shim.Error("order is already rejected by the supplier")
	}else {
		_TraceingLogger.Errorf("OrderBySupplier : order is already inStock")
		return shim.Error("order is already inStock")
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":        Order.PONumber,
		"message":       "Delivery Order created Successfully by supplier",
		"Order":         Order,
	}

		respJSON,_ := json.Marshal(resultData)
		return shim.Success(respJSON)

}

// Logistics Approval

func (tr *Trace) logisticApproval(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("logiscticsApproval")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("logiscticsApproval : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var logisticOrder PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &logisticOrder)
	if err != nil {
		_TraceingLogger.Errorf("logiscticsApproval : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(logisticOrder.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("logiscticsApproval : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("logiscticsApproval : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
    if existingOrder.DoStatus == "expecting confirmation from regulator" {

		existingOrder.DoStatus = logisticOrder.DoStatus   // shipped or dispute
		existingOrder.UpdateTs = logisticOrder.UpdateTs
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("logiscticsApproval : Marshalling Error : " + string(err.Error()))
			return shim.Error("logiscticsApproval : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("logiscticsApproval : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(logisticOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("logiscticsApproval : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("logiscticsApproval : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_LogiscticApproval, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("logiscticsApproval : Event not generating for : " + _LogiscticApproval)
		}
	}else {
		_TraceingLogger.Errorf("logiscticsApproval : order is not in expecting confirmation from regulator state")
		return shim.Error("order is not in expecting confirmation from regulator state")
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       logisticOrder.PONumber,
		"message":       "Delivery Order approved by the regulator successfully.",
		"Order":        logisticOrder,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)

}

//Inventory Manager Creats Good Receipt

func (tr *Trace) InventoryManagerReceipt(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("InventoryManagerReciept")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("InventoryManagerReciept : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var receiptOrder PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &receiptOrder)
	if err != nil {
		_TraceingLogger.Errorf("InventoryManagerReciept : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(receiptOrder.PONumber)
	if len(RecordBytes) <= 0 {
		_TraceingLogger.Errorf("InventoryManagerReciept : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("InventoryManagerReciept : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
	 
	if existingOrder.DoStatus == "pending" {
		_TraceingLogger.Errorf("InventoryManagerReciept : Order is not shipped yet")
		return shim.Error("Order is not shipped yet")
	}
	if existingOrder.DoStatus == "expecting confirmation from regulator" {
		_TraceingLogger.Errorf("InventoryManagerReciept : order is waiting confirmation from regulator")
		return shim.Error("order is waiting confirmation from regulator")
	}
	if existingOrder.DoStatus == "dispute" {
		_TraceingLogger.Errorf("InventoryManagerReciept : It should go back to supplier for further correction")
		return shim.Error("It should go back to supplier for further correction")	
	}
	if existingOrder.DoStatus == "shipped" {

		existingOrder.InvMngId = receiptOrder.InvMngId
		existingOrder.ExpDate  = receiptOrder.ExpDate
		existingOrder.StockLocation = receiptOrder.StockLocation
		existingOrder.GoodReceipt = receiptOrder.GoodReceipt
		existingOrder.GRStatus = receiptOrder.GRStatus   //  pending
		existingOrder.UpdateTs = receiptOrder.UpdateTs
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("InventoryManagerReciept : Marshalling Error : " + string(err.Error()))
			return shim.Error("InventoryManagerReciept : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("InventoryManagerReciept : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(receiptOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("InventoryManagerReciept : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("logiscticsApproval : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_InventoryManager, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("logiscticsApproval : Event not generating for : " + _InventoryManager)
		}
	}else {
		_TraceingLogger.Errorf("InventoryManagerReciept : Order is not shipped yet")
		return shim.Error("Order is not shipped yet")	
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       receiptOrder.PONumber,
		"message":       "Inventory Manager Created the Receipt Successfully.",
		"Order":        receiptOrder,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)

}

//Inventory Approval

func (tr *Trace) inventoryApproval(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("inventoryApproval")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("inventoryApproval : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var inventoryOrder PurchaseOrder
	err := json.Unmarshal([]byte(args[0]), &inventoryOrder)
	if err != nil {
		_TraceingLogger.Errorf("inventoryApproval : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(inventoryOrder.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("inventoryApproval : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("inventoryApproval : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
    if existingOrder.GRStatus == "pending" && existingOrder.ItemName == "Pipe"{
        if existingOrder.Standard[1] == inventoryOrder.Outerdia && existingOrder.Standard[0] == inventoryOrder.Innerdia && existingOrder.Standard[2] == inventoryOrder.Wallwidth {
            if existingOrder.Quantity * existingOrder.Standard[3] == inventoryOrder.StadBatchWeght {
				existingOrder.DoStatus = inventoryOrder.DoStatus  // arrived 
				existingOrder.PoStatus = inventoryOrder.PoStatus   // inStock
				existingOrder.GRStatus = inventoryOrder.GRStatus   // received 
				existingOrder.UpdateTs = inventoryOrder.UpdateTs
				existingOrder.Ownership = existingOrder.InvMngId
			} else {
				existingOrder.GRStatus = "backorder"   //backorder(should go back to supplier)
			}
		} else {
			existingOrder.GRStatus = "backorder"   //backorder(should go back to supplier)
		}
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("inventoryApproval : Marshalling Error : " + string(err.Error()))
			return shim.Error("inventoryApproval : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("inventoryApproval : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(inventoryOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("inventoryApproval : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("inventoryApproval : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_InventoryApproval, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("inventoryApproval : Event not generating for : " + _InventoryApproval)
		}
	}else if existingOrder.GRStatus == "pending" && existingOrder.ItemName == "Cement" {
		if existingOrder.Quantity * existingOrder.Standard[0] == inventoryOrder.StadBatchWeght {
			existingOrder.DoStatus = inventoryOrder.DoStatus  // arrived 
			existingOrder.PoStatus = inventoryOrder.PoStatus   // inStock
			existingOrder.GRStatus = inventoryOrder.GRStatus   // received 
			existingOrder.UpdateTs = inventoryOrder.UpdateTs
			existingOrder.Ownership = existingOrder.InvMngId
		} else {
			existingOrder.GRStatus = "backorder"   //backorder(should go back to supplier)
		}
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("inventoryApproval : Marshalling Error : " + string(err.Error()))
			return shim.Error("inventoryApproval : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("inventoryApproval : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(inventoryOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("inventoryApproval : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("inventoryApproval : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_InventoryApproval, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("inventoryApproval : Event not generating for : " + _InventoryApproval)
		}
	}else {
		_TraceingLogger.Errorf("inventoryApproval : good receipt status is not in pending state or itemName is not provided") 
		return shim.Error("good receipt status is not in pending state or itemName is not provided")
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       inventoryOrder.PONumber,
		"message":       "Inventory receipt approved by the regulator successfully.",
		"Order":         inventoryOrder,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)
}

//Create Consumption Order

func (tr *Trace) FormenConsumption(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("formenConsumption")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("formenConsumption : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var formenOrder ForemenType
	err := json.Unmarshal([]byte(args[0]), &formenOrder)
	if err != nil {
		_TraceingLogger.Errorf("formenConsumption : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	if formenOrder.PONumber == "" {
		jsonResp := "PO number is missing: PONumber"
		_TraceingLogger.Errorf(jsonResp)
		return shim.Error(jsonResp)
	}
	ponumber := formenOrder.PONumber
	RecordBytes, _:= stub.GetState(formenOrder.PONumber)
	if len(RecordBytes) <= 0 {
		_TraceingLogger.Errorf("formenConsumption : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("formenConsumption : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
    // if existingOrder.Description == "cement" {
		if existingOrder.Quantity < formenOrder.PQuantity {
			_TraceingLogger.Errorf("formenConsumption : quantity is not available")
			return shim.Error("quantity is not available")
		}
	// }
    if existingOrder.GRStatus == "received" {	
		for _, j := range existingOrder.ForemenUpdate {
			if j.CONumber ==  formenOrder.CONumber {
			    return shim.Error("consumption order already exist : " + formenOrder.CONumber)
			}
		}
		existingOrder.Quantity = existingOrder.Quantity - formenOrder.PQuantity
		formenOrder.PONumber = ""
		existingOrder.ForemenUpdate = append(existingOrder.ForemenUpdate,formenOrder)
		// existingOrder.ForemenId = formenOrder.ForemenId 
		// existingOrder.Purpuse = formenOrder.Purpuse   // for  making concreate
		// existingOrder.CCOrder = formenOrder.CCOrder  // created
		// existingOrder.UpdateTs = formenOrder.UpdateTs
		// existingOrder.ForDesc = existingOrder.Description
		// existingOrder.CONumber = formenOrder.CONumber
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("formenConsumption : Marshalling Error : " + string(err.Error()))
			return shim.Error("formenConsumption : Marshalling Error : " + string(err.Error()))	
		}
		_TraceingLogger.Infof("formenConsumption : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(ponumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("formenConsumption : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("formenConsumption : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_Foremen, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("inventoryApproval : Event not generating for : " + _Foremen)
		}
	}else if existingOrder.GRStatus == "backorder" {
		_TraceingLogger.Errorf("formenConsumption : Goods order has been sent back to supplier") 
		return shim.Error("Goods order has been sent back to supplier") 
	}else {
		_TraceingLogger.Errorf("formenConsumption : Goods order is not received yet") 
		return shim.Error("Goods order is not received yet") 
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       ponumber,
		"message":       "Consumption Order created successfully by foremen",
		"Order":         formenOrder,
	}

	respJSON,_ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

// Inventory Manager stock release
func (tr *Trace) stockRelease(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("stockRelease")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("stockRelease : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var stockRelease ForemenType
	err := json.Unmarshal([]byte(args[0]), &stockRelease)
	if err != nil {
		_TraceingLogger.Errorf("stockRelease : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(stockRelease.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("stockRelease : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("stockRelease : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
	// if existingOrder.Quantity <= 0 && existingOrder.PQuantity <= existingOrder.Quantity{
	// 	_TraceingLogger.Errorf("stockRelease : there is no quantity available")
	// 	return shim.Error("there is no quantity available")		
	// }
	var response = false
	var indx = 0
	for index, j := range existingOrder.ForemenUpdate {
            if j.CONumber == stockRelease.CONumber {
				if j.CCOrder == "created" {
					response = true
					indx = index
				}else {
		            _TraceingLogger.Errorf("stockRelease : consumtion order not in created state")
					return shim.Error("consumtion order not in created state")
				}
		    }
		}

    if response {
		existingOrder.ForemenUpdate[indx].CCOrder = stockRelease.CCOrder   // expecting confirmation from regulator
		// existingOrder.ForemenUpdate[indx].futs = stockRelease.UpdateTs
		// existingOrder.PQuantity = stockRelease.PQuantity
		// existingOrder.Quantity =  existingOrder.Quantity - stockRelease.PQuantity
		existingOrder.ForemenUpdate[indx].BatchId = stockRelease.BatchId 
	}else {
		_TraceingLogger.Errorf("stockRelease : consumption order does not exist")
		return shim.Error("consumption order does not exist")
	}

	OrderBytes, err := json.Marshal(existingOrder)
	if err !=nil {
		_TraceingLogger.Errorf("stockRelease : Marshalling Error : " + string(err.Error()))
	    return shim.Error("stockRelease : Marshalling Error : " + string(err.Error()))
	}
	_TraceingLogger.Infof("stockRelease : saving the Create Order : " + existingOrder.PONumber)

	errorr :=stub.PutState(stockRelease.PONumber, OrderBytes)

	if errorr != nil {
		_TraceingLogger.Errorf("stockRelease : Put State Failed Error : " + string(errorr.Error())) 
		return shim.Error("Put State Failed Error : " + string(errorr.Error()))
	}

	_TraceingLogger.Infof("stockRelease : PutState Success : " + string(OrderBytes))
	err2 := stub.SetEvent(_StockRelease, OrderBytes)
	if err2 != nil {
		_TraceingLogger.Errorf("inventoryApproval : Event not generating for : " + _StockRelease)
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       stockRelease.PONumber,
		"message":       "Stock Release by Inventory Manager Successfully.",
		"Order":         stockRelease,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)
}

// Consumption Approval by Regulator for meterial release to foremen
func (tr *Trace) consumptionApproval(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("consumptionApproval")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("consumptionApproval : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var consumtionOrder ForemenType
	err := json.Unmarshal([]byte(args[0]), &consumtionOrder)
	if err != nil {
		_TraceingLogger.Errorf("consumptionApproval : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(consumtionOrder.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("consumptionApproval : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("consumptionApproval : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}

	var response = false
	var indx = 0
	for index, j := range existingOrder.ForemenUpdate {
            if j.CONumber == consumtionOrder.CONumber {
				if j.CCOrder == "expecting confirmation from regulator"{
					response = true
					indx = index
				}else{
			        _TraceingLogger.Errorf("consumptionApproval :consumption order not in waiting state for regulator approval")
					return shim.Error("consumption order not in waiting state for regulator approval")
				}
		    }
		}
    if response {

		existingOrder.ForemenUpdate[indx].CCOrder = consumtionOrder.CCOrder   // ready to use 
		existingOrder.ForemenUpdate[indx].BatchWeight = consumtionOrder.BatchWeight
		existingOrder.ForemenUpdate[indx].UpdateTs = consumtionOrder.UpdateTs
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("consumptionApproval : Marshalling Error : " + string(err.Error()))
			return shim.Error("consumptionApproval : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("consumptionApproval : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(consumtionOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("consumptionApproval : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("consumptionApproval : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_ConsumtionOrder, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("consumptionApproval : Event not generating for : " + _ConsumtionOrder)
		}
	}else {
		_TraceingLogger.Errorf("consumptionApproval : consumption order does not exist") 
		return shim.Error("consumption order does not exist")
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       consumtionOrder.PONumber,
		"message":       "Consumtion Order is approved by the Regulator successfully.",
		"Order":         consumtionOrder,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)
}

// // Display Order Status
func (tr *Trace) displayOrderStatus(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("foremenOrderStatus")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("foremenOrderStatus : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var displayOrder ForemenType
	err := json.Unmarshal([]byte(args[0]), &displayOrder)
	if err != nil {
		_TraceingLogger.Errorf("foremenOrderStatus : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(displayOrder.PONumber)
	if len(RecordBytes) <= 0{
		_TraceingLogger.Errorf("foremenOrderStatus : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("foremenOrderStatus : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
	var indx = 0
	var response = false
	for index, j := range existingOrder.ForemenUpdate {
		if j.CONumber == displayOrder.CONumber {
			if j.CCOrder == "ready to use" {
				response = true
				indx = index
			}else{
		        _TraceingLogger.Errorf("foremenOrderStatus : consumption order not in ready to use state")
				return shim.Error("consumption order not in ready to use state")
			}
		}
	}
    if response {

		existingOrder.ForemenUpdate[indx].CCOrder = displayOrder.CCOrder   // expecting confirmation from regulator for pouring
		existingOrder.ForemenUpdate[indx].UpdateTs = displayOrder.UpdateTs
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("foremenOrderStatus : Marshalling Error : " + string(err.Error()))
			return shim.Error("foremenOrderStatus : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("foremenOrderStatus : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(displayOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("foremenOrderStatus : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("foremenOrderStatus : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_DisplayOrder, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("foremenOrderStatus : Event not generating for : " + _DisplayOrder)
		}
		
	}else {
		_TraceingLogger.Errorf("foremenOrderStatus : order does not exist")
		return shim.Error("order does not exist")
	}
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       displayOrder.PONumber,
		"message":       "Expecting confirmation from Regulator for pouring",
		"Order":         displayOrder,
	}

	 respJSON,_ := json.Marshal(resultData)
	 return shim.Success(respJSON)
}


// // Consumption Approval for pouring
func (tr *Trace) consumptionApprovalForPouring(stub shim.ChaincodeStubInterface) peer.Response {

	_TraceingLogger.Infof("approvalForPouring")
	_, args:= stub.GetFunctionAndParameters()
	if len(args) < 1 {
		_TraceingLogger.Errorf("approvalForPouring : Invalid number of argument provided")
		return shim.Error("Invalid number of argument provided")
	}
	var approvalOrder ForemenType
	err := json.Unmarshal([]byte(args[0]), &approvalOrder)
	if err != nil {
		_TraceingLogger.Errorf("approvalForPouring : Invalid json provided as input")
		return shim.Error("Invalid json provided as input")
	}
	authorize, _ := tr.getInvokerIdentity(stub)

	if authorize == false {
		return shim.Error("Unauthorize access")
	}
	RecordBytes, _:= stub.GetState(approvalOrder.PONumber)
	if len(RecordBytes) <= 0 {
		_TraceingLogger.Errorf("approvalForPouring : Order does not exist")
		return shim.Error("Order does not exist")
	}
	var existingOrder PurchaseOrder 
	errOrder := json.Unmarshal([]byte(RecordBytes), &existingOrder)

	if errOrder != nil {
		_TraceingLogger.Errorf("approvalForPouring : existing Order Unmarshaling Error")
		return shim.Error("existing Order Unmarshaling Error")
	}
	var response = false
	var indx = 0;
	for index, j := range existingOrder.ForemenUpdate {
		if j.CONumber == approvalOrder.CONumber {
			if j.CCOrder == "expecting confirmation from regulator for pouring" {
				indx = index
				response = true
			}else {
		        _TraceingLogger.Errorf("approvalForPouring : consumption order not in pouring state")
				return shim.Error("consumption order not in pouring state")
			}
		}
	}
    if response {

		existingOrder.ForemenUpdate[indx].CCOrder = approvalOrder.CCOrder   // ready to be poured
		existingOrder.ForemenUpdate[indx].Density = approvalOrder.Density
		existingOrder.ForemenUpdate[indx].UpdateTs = approvalOrder.UpdateTs
	
		OrderBytes, err := json.Marshal(existingOrder)
		if err !=nil {
			_TraceingLogger.Errorf("approvalForPouring : Marshalling Error : " + string(err.Error()))
			 return shim.Error("approvalForPouring : Marshalling Error : " + string(err.Error()))
		}
		_TraceingLogger.Infof("approvalForPouring : saving the Create Order : " + existingOrder.PONumber)
	
		errorr :=stub.PutState(approvalOrder.PONumber, OrderBytes)
	
		if errorr != nil {
			_TraceingLogger.Errorf("approvalForPouring : Put State Failed Error : " + string(errorr.Error())) 
			return shim.Error("Put State Failed Error : " + string(errorr.Error()))
		}
	
		_TraceingLogger.Infof("approvalForPouring : PutState Success : " + string(OrderBytes))
		err2 := stub.SetEvent(_ApprovalForPouring, OrderBytes)
		if err2 != nil {
			_TraceingLogger.Errorf("approvalForPouring : Event not generating for : " + _ApprovalForPouring)
		}
	}else {
		_TraceingLogger.Errorf("approvalForPouring : consumtion Order does not exist")
		return shim.Error("consumtion Order does not exist")
	}
	
	resultData := map[string]interface{}{
		"trxnID":        stub.GetTxID(),
		"OPNumer":       approvalOrder.PONumber,
		"message":       "Consumtion order approved by the Regulator for pouring",
		"Order":         approvalOrder,
	}

	respJSON,_ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

// Customer query for Material by GTIN

func (tr *Trace) materialQuery(stub shim.ChaincodeStubInterface, args []string)  peer.Response {
    if len(args) < 1 {
		shim.Error("Invalid number of argument provided")
	}
	var searchCriteria = `{
		"obj" : "PurchaseOrder",
		"gtin" : "%s"
		}`

		orders := tr.retriveMaterial(stub, fmt.Sprintf(searchCriteria, args[0]))

		order,_ := json.Marshal(orders)
		return shim.Success(order)     
}

func (tr *Trace) retriveMaterial(stub shim.ChaincodeStubInterface, criteria string, indexs ...string) []PurchaseOrder {

	var finalSelector string
	records := make([]PurchaseOrder, 0)

	if len(indexs) == 0 {
		finalSelector = fmt.Sprintf("{\"selector\":%s }", criteria)

	} else {
		finalSelector = fmt.Sprintf("{\"selector\":%s , \"use_index\" :\"%s\" }", criteria, indexs[0])
	}

	_TraceingLogger.Infof("Query Selector : %s", finalSelector)
	resultsIterator, _ := stub.GetQueryResult(finalSelector)
	for resultsIterator.HasNext() {
		order := PurchaseOrder{}
		recordBytes, _ := resultsIterator.Next()
		err := json.Unmarshal(recordBytes.Value, &order)
		if err != nil {
			_TraceingLogger.Infof("Unable to unmarshal Order retrived:: %v", err)
		}
		records = append(records, order)
	}
	return records
}



//Returns the complete identity in the format
//Certitificate issuer orgs's domain name
//Returns string Unkown if not able to parse the invoker certificate
func (tr *Trace) getInvokerIdentity(stub shim.ChaincodeStubInterface) (bool, string) {
	//Following id comes in the format X509::<Subject>::<Issuer>>
	enCert, err := cid.GetX509Certificate(stub)
	if err != nil {
		return false, "Unknown."
	}

	issuersOrgs := enCert.Issuer.Organization
	if len(issuersOrgs) == 0 {
		return false, "Unknown.."
	}
	isOK, msg := isValidDomainName(issuersOrgs[0])
	if !isOK {
		return false, msg
	}
	return true, fmt.Sprintf("%s", issuersOrgs[0])

}
