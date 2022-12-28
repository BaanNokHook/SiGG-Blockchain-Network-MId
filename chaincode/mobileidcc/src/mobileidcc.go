// @author NBTC, developed by Chanwanich Co., Ltd.
// @version 1.2.0 remove shim.Logger

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/msp"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/shopspring/decimal"

	"github.com/op/go-logging"
)

// *********************************************
// MAIN
// *********************************************

func main() {
	mob := new(MobileIdChaincode)
	//for set Development mode
	mob.devMode = false
	err := shim.Start(mob)
	if err != nil {
		fmt.Printf("Error starting MobileId Chaincode: %s", err)
	}
}

// *********************************************
// Chaincode
// *********************************************

type MobileIdChaincode struct {
	devMode bool
}

// Init Method
func (c *MobileIdChaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_loguuid := uuid.New().String()
	logger := newShimLogger(_loguuid, "Init Chaincode", LogInfo)
	logger.Info("Start")

	logger.Info("Init Chaincode...")

	_memberInfo := Member{}
	_memberInfo.MemberRole = MemberRegulator
	_memberInfo.Status = MemberActive
	logger.Info("createMember NBTC...")
	var args []string
	args = append(args, `{"member_code": "NBT","member_name": "NBTC","member_role": "REGULATOR","status": "A"}`)
	c.createMember(stub, _memberInfo, args, _loguuid)
	logger.Info("createMember NBTC...Success!")

	logger.Info("createMember AIS")
	args[0] = `{"member_code": "AIS","member_name": "Advanced Info Service","member_role": "ISSUER,VERIFIER", "serial_no_prefix": "03","status": "A"}`
	c.createMember(stub, _memberInfo, args, _loguuid)
	logger.Info("createMember AIS...Success!")

	logger.Info("createMember BBL")
	args[0] = `{"member_code": "BBL","member_name": "Bangkok Bank","member_role": "VERIFIER","status": "A"}`
	c.createMember(stub, _memberInfo, args, _loguuid)
	logger.Info("createMember BBL...Success!")

	if c.devMode {
		logger.Info("createMember Org1MSP")
		args[0] = `{"member_code": "Org1MSP","member_name": "Org1 MSP","member_role": "ISSUER,VERIFIER,REGULATOR","serial_no_prefix": "01","status": "A"}`
		c.createMember(stub, _memberInfo, args, _loguuid)
		logger.Info("createMember Org1MSP...Success!")

		logger.Info("createMember Org2MSP")
		args[0] = `{"member_code": "Org2MSP","member_name": "Org2 MSP","member_role": "ISSUER,VERIFIER","serial_no_prefix": "02","status": "A"}`
		c.createMember(stub, _memberInfo, args, _loguuid)
		logger.Info("createMember Org2MSP...Success!")
	}

	return shimSuccess(logger, nil, MessageSuccess)
}

// Invoke Method
func (c *MobileIdChaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	_loguuid := uuid.New().String()
	function, args := stub.GetFunctionAndParameters()

	logger := newShimLogger(_loguuid, "Invoke Chaincode: "+function, LogInfo)
	logger.Debug("Start")

	_memberInfo, err := c.getMspInfo(stub, _loguuid)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	logger.Info(fmt.Sprintf("%s %s %s %s %s\n", _memberInfo.MemberCode,
		_memberInfo.MemberName, _memberInfo.MemberRole,
		_memberInfo.Status, _memberInfo.Registered))

	logger.Info(fmt.Sprintf("Access by %s\n", _memberInfo.MemberCode))

	// mobileId function lists
	switch function {
	// ******* HealthCheck Function
	case "invokeHealthCheck":
		return c.invokeHealthCheck(stub, _memberInfo, _loguuid)
	case "listHealthCheck":
		return c.listHealthCheck(stub, _loguuid)
	// ******* Membership Function
	case "getMember":
		return c.getMember(stub, _memberInfo, args, _loguuid)
	case "listMembers":
		return c.listMembers(stub, _memberInfo, _loguuid)
	case "createMember":
		return c.createMember(stub, _memberInfo, args, _loguuid)
	case "updateMember":
		return c.updateMember(stub, _memberInfo, args, _loguuid)
	// ******* MobileId Function
	case "enrollMobileId":
		return c.enrollMobileId(stub, _memberInfo, args, _loguuid)
	case "updateMobileId":
		return c.updateMobileId(stub, _memberInfo, args, _loguuid)
	case "retrieveMobileIdAndUpdateConsentLog":
		return c.retrieveMobileIdAndUpdateConsentLog(stub, _memberInfo, args, _loguuid)
	// ******* Consent Function
	case "recordConsentLog":
		return c.recordConsentLog(stub, _memberInfo, args, _loguuid)
	case "recordVerificationResult":
		return c.recordVerificationResult(stub, _memberInfo, args, _loguuid)
	case "getConsentLog":
		return c.getConsentLog(stub, _memberInfo, args, _loguuid)
	// ******* Audit Function
	case "getMobileIdForAudit":
		return c.getMobileIdForAudit(stub, _memberInfo, args, _loguuid)
	case "getConsentLogForAudit":
		return c.getConsentLogForAudit(stub, _memberInfo, args, _loguuid)
	case "listConsentLogForAudit":
		return c.listConsentLogForAudit(stub, _memberInfo, args, _loguuid)
	case "listConsentLogByIssuer":
		return c.listConsentLogByIssuer(stub, _memberInfo, args, _loguuid)
	case "listConsentLogByVerifier":
		return c.listConsentLogByVerifier(stub, _memberInfo, args, _loguuid)
	case "getConsentHistory":
		return c.getConsentHistory(stub, _memberInfo, args, _loguuid)
	// ******* Enhance-1 Function
	case "getMobileId":
		return c.getMobileId(stub, _memberInfo, args, _loguuid)
	case "countMobileIdByStatus":
		return c.countMobileIdByStatus(stub, _memberInfo, args, _loguuid)
	case "listConsentLogByDate":
		return c.listConsentLogByDate(stub, _memberInfo, args, _loguuid)
	// ******* Enhance-2 add revoke method
	case "revokeConsentLog":
		return c.revokeConsentLog(stub, _memberInfo, args, _loguuid)
	// ******* MobileId 1.5 Extend
	case "retrieveMobileIdIssuer":
		return c.retrieveMobileIdIssuer(stub, _memberInfo, args, _loguuid)
	case "listRequestLogByDate":
		return c.listRequestLogByDate(stub, _memberInfo, args, _loguuid)
	case "getMobileIdRequestLog":
		return c.getMobileIdRequestLog(stub, _memberInfo, args, _loguuid)
	default:
		return shimError(logger, StatusCodeNotFound, "Invalid/Unknown invoke function name: "+function)
	}
}

//*******************************************************
// MobileId - Get Key Utils
//*******************************************************

func createMobileIdKey(stub shim.ChaincodeStubInterface, req MobileId) (string, error) {
	mbIdKey, err := stub.CreateCompositeKey(MobileIdKeyObjectType,
		[]string{
			"mobile_no", req.Dg1.MobileNo,
		})
	return mbIdKey, err
}

func createConsentIdKey(stub shim.ChaincodeStubInterface, req ConsentLog) (string, error) {
	csIDKey, err := stub.CreateCompositeKey(MobileIdConsentKeyObjectType,
		[]string{
			"mobile_no", req.MobileNo,
			"mobile_id_sn", req.MobileIdSn,
			"cid", req.Cid,
		})
	return csIDKey, err
}

func createMemberKey(stub shim.ChaincodeStubInterface, req Member) (string, error) {
	mbIDKey, err := stub.CreateCompositeKey(MemberKeyObjectType,
		[]string{
			"member_code", req.MemberCode,
		})
	return mbIDKey, err
}

func createHealthCheckKey(stub shim.ChaincodeStubInterface, req HealthCheck) (string, error) {
	mbIDKey, err := stub.CreateCompositeKey(MemberKeyObjectType,
		[]string{
			"member_healthcheck", req.MemberCode,
		})
	return mbIDKey, err
}

func createMobileIdRequestKey(stub shim.ChaincodeStubInterface, req MobileIdRequestLog) (string, error) {
	csIDKey, err := stub.CreateCompositeKey(MobileIdRequestKeyObjectType,
		[]string{
			"mobile_no", req.MobileNo,
			"rp_id", req.RpId,
			"tx_id", req.TxId,
		})
	return csIDKey, err
}

// *********************************************
// HealthCheck functions
// *********************************************

// Invoke HealthCheck
func (c *MobileIdChaincode) invokeHealthCheck(stub shim.ChaincodeStubInterface, _member Member, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "invokeHealthCheck", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator, MemberIssuer, MemberVerifier}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	_hcReq := HealthCheck{}
	err := errors.New("")
	_hcReq.Timestamp, err = getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_hcReq.MemberCode = _member.MemberCode

	hcKey, err := createHealthCheckKey(stub, _hcReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	_hcByte, err := json.Marshal(_hcReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(hcKey, _hcByte)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	logger.Debug("MemberCode: " + _hcReq.MemberCode)
	logger.Debug("Timestamp: " + _hcReq.Timestamp)
	// Return response
	return shimSuccess(logger, _hcByte, MessageSuccess)
}

// List HealthCheck
func (c *MobileIdChaincode) listHealthCheck(stub shim.ChaincodeStubInterface, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listHealthCheck", LogInfo)
	logger.Debug("Start")

	iterator, err := stub.GetStateByPartialCompositeKey(MemberKeyObjectType, []string{"member_healthcheck"})
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	var _hcs []HealthCheck
	defer iterator.Close()
	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _hcInfo HealthCheck
		err = json.Unmarshal(kv.Value, &_hcInfo)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_hcs = append(_hcs, _hcInfo)
	}
	_memberBytes, err := json.Marshal(_hcs)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _memberBytes, MessageSuccess)
}

// *********************************************
// Member Management functions
// *********************************************

// Gets Member in the system (Internal Use)
func (c *MobileIdChaincode) getMspInfo(stub shim.ChaincodeStubInterface, loguuid string) (Member, error) {
	logger := newShimLogger(loguuid, "getMspInfo", LogInfo)
	logger.Debug("Start")

	_memberInfo := Member{}

	creator, err := stub.GetCreator()
	if err != nil {
		return _memberInfo, err
	}

	si := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creator, si)
	if err != nil {
		return _memberInfo, err
	}
	logger.Info(fmt.Sprintf("MSPName: %s", si.Mspid))

	//begin mockup data for Node Org1MSP,Org2MSP in devMode
	//si.Mspid = "NBT"
	_memberString := fmt.Sprintf(`{"member_code": "%s"}`, si.Mspid)
	if !c.devMode {
		if len(si.Mspid) >= 3 {
			si.Mspid = si.Mspid[0:3]
		}
		_memberString = fmt.Sprintf(`{"member_code": "%s"}`, strings.ToUpper(si.Mspid))
	}
	//end mockup data for Node Org1MSP,Org2MSP in devMode

	// Unmarshal and validate request.
	_memberReq := Member{}
	err = json.Unmarshal([]byte(_memberString), &_memberReq)
	if err != nil {
		return _memberInfo, err
	}

	// get mobileId member by composite key
	_memberKey, err := createMemberKey(stub, _memberReq)
	if err != nil {
		return _memberInfo, err
	}

	// get world state by key
	_memberBytes, err := stub.GetState(_memberKey)
	if err != nil {
		return _memberInfo, err
	}
	if _memberBytes == nil {
		return _memberInfo, errors.New("Member does not exist")
	}

	err = json.Unmarshal(_memberBytes, &_memberInfo)
	if err != nil {
		return _memberInfo, err
	}

	return _memberInfo, nil
}

// Gets Member in the system
func (c *MobileIdChaincode) getMember(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getMember", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_memberReq := Member{}
	err := json.Unmarshal([]byte(args[0]), &_memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Validate member code field.
	if err := validateMemberCode("member_code", _memberReq.MemberCode, c.devMode); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// get mobileId member by composite key.
	_memberKey, err := createMemberKey(stub, _memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// get world state by key.
	_memberBytes, err := stub.GetState(_memberKey)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	if _memberBytes == nil {
		return shimError(logger, StatusCodeNotFound, "Member is not exist")
	}

	// Return response
	return shimSuccess(logger, _memberBytes, MessageSuccess)
}

// Lists Members in the system
func (c *MobileIdChaincode) listMembers(stub shim.ChaincodeStubInterface, _member Member, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listMembers", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	iterator, err := stub.GetStateByPartialCompositeKey(MemberKeyObjectType, []string{"member_code"})
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	var _members []Member
	defer iterator.Close()
	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _memberInfo Member
		err = json.Unmarshal(kv.Value, &_memberInfo)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_members = append(_members, _memberInfo)
	}
	_memberBytes, err := json.Marshal(_members)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _memberBytes, MessageSuccess)
}

// Creates Member and saves into the system
func (c *MobileIdChaincode) createMember(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "createMember", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_memberReq := Member{}
	err := json.Unmarshal([]byte(args[0]), &_memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMember(_memberReq, c.devMode); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Validate serial no prefix field.
	if strings.Contains(_memberReq.MemberRole, MemberIssuer) {
		if err := validateSerialNoPrefix("serial_no_prefix", _memberReq.SerialNoPrefix); err != nil {
			return shimError(logger, StatusCodeInvalidArgument, err.Error())
		}
	}
	if _memberReq.Status != MemberActive {
		return shimError(logger, StatusCodeInvalidArgument, "Member status must active.")
	}

	// Precondition: Member must not exist
	memKey, err := createMemberKey(stub, _memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	b, err := stub.GetState(memKey)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	if b != nil {
		return shimError(logger, StatusCodeInvalidOperation, "Member must not exist.")
	}

	// Update Member and saves into the system
	_memberReq.Registered, err = getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	logger.Info(fmt.Sprintf("%s %s %s %s %s", _memberReq.MemberCode, _memberReq.MemberName,
		_memberReq.MemberRole, _memberReq.Status, _memberReq.Registered))

	_memberByte, err := json.Marshal(_memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(memKey, _memberByte)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, nil, MessageSuccess)
}

// Updates Member in the system
func (c *MobileIdChaincode) updateMember(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "updateMember", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_memberReq := Member{}
	err := json.Unmarshal([]byte(args[0]), &_memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateMember(_memberReq, c.devMode); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Validate serial no prefix field.
	if strings.Contains(_memberReq.MemberRole, MemberIssuer) {
		if err := validateSerialNoPrefix("serial_no_prefix", _memberReq.SerialNoPrefix); err != nil {
			return shimError(logger, StatusCodeInvalidArgument, err.Error())
		}
	}

	// Precondition: Member must exist.
	memKey, err := createMemberKey(stub, _memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	b, err := stub.GetState(memKey)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	if b == nil {
		return shimError(logger, StatusCodeInvalidOperation, "Member must exist.")
	}

	//
	_memberInfo := Member{}
	err = json.Unmarshal(b, &_memberInfo)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_memberInfo.MemberCode = _memberReq.MemberCode
	_memberInfo.MemberName = _memberReq.MemberName
	_memberInfo.MemberRole = _memberReq.MemberRole
	_memberInfo.SerialNoPrefix = _memberReq.SerialNoPrefix
	_memberInfo.Status = _memberReq.Status
	_memberInfo.ServiceUrl = _memberReq.ServiceUrl

	logger.Info(fmt.Sprintf("%s %s %s %s %s %s", _memberInfo.MemberCode, _memberInfo.MemberName,
		_memberInfo.MemberRole, _memberInfo.Status, _memberInfo.Registered, _memberInfo.SerialNoPrefix))

	_memberByte, err := json.Marshal(_memberInfo)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(memKey, _memberByte)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	return shimSuccess(logger, nil, MessageSuccess)
}

// *********************************************
// MobileId functions
// *********************************************

// Enrolls MobileId data into the system (for Issuer).
func (c *MobileIdChaincode) enrollMobileId(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "enrollMobileId", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_mobileIDReq := MobileId{}
	err := json.Unmarshal([]byte(args[0]), &_mobileIDReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdFormat(_mobileIDReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	logger.Debug(_mobileIDReq.Dg1.MobileNo)

	// Precondition: The issuer field in the request must be equal to member's code.
	if _mobileIDReq.Dg1.Issuer != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request must be equal to member's code (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Precondition: The status field in the request must be active.
	if _mobileIDReq.Dg1.Status != MobileIdActive {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The status field in the request must be active (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Precondition: The timestamp field must be in 24hr range.
	currentTimeStr, err := getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	difSec, err := compareTimeString(currentTimeStr, _mobileIDReq.Dg1.Timestamp)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if int64(math.Abs(float64(difSec))) > TimestampRange {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The timestamp field must be in 24hr range (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Precondition: In case there is existing MobileId in the system...
	_mobileIDKey, err := createMobileIdKey(stub, _mobileIDReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if mobileIdExists(stub, _mobileIDReq.Dg1.MobileNo) {
		_mobileID, err := getMobileIDByNo(stub, _mobileIDReq.Dg1.MobileNo)
		if err != nil {
			return shimError(logger, StatusCodeNotFound, err.Error())
		}

		// Precondition: The existing MobileId must be terminated.
		if _mobileID.Dg1.Status != MobileIdTerminated {
			return shimErrorf(logger, StatusCodeInvalidOperation, "The existing MobileId must be terminated. (mobile no: %s, status: %s)", _mobileID.Dg1.MobileNo, _mobileID.Dg1.Status)
		}

		// Precondition: The MobileId to enroll must have not been enrolled before (the serial number must not present in blockchain history).
		historyIter, err := stub.GetHistoryForKey(_mobileIDKey)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		defer historyIter.Close()
		for historyIter.HasNext() {
			_keyModification, err := historyIter.Next()
			if err != nil {
				return shimError(logger, StatusCodeInternalError, err.Error())
			}
			var _mobileIDTmp MobileId
			err = json.Unmarshal(_keyModification.Value, &_mobileIDTmp)
			if err != nil {
				return shimError(logger, StatusCodeInternalError, err.Error())
			}
			if _mobileIDReq.Dg1.MobileIdSn == _mobileIDTmp.Dg1.MobileIdSn {
				return shimErrorf(logger, StatusCodeInvalidOperation, "The MobileId to enroll must have not been enrolled before (mobile no: %s).", _mobileID.Dg1.MobileNo)
			}
		}
	}

	/* skip for case hash MobileIdSn
	// Precondition: The serial number in the request must start with the last two digits of issuer's MNC.
	if _mobileIDReq.Dg1.MobileIdSn[0:2] != _member.SerialNoPrefix {
		return shimErrorf(logger, StatusCodeInvalidArgument, "The serial number in the request must start with the last two digits of issuer's MNC (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}*/

	// Marshal and save MobileId data into the system.
	_mobileIDBytes, err := json.Marshal(_mobileIDReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(_mobileIDKey, _mobileIDBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response.
	_responseObj := getMessageSuccessResponse(stub.GetTxID())
	return shimSuccess(logger, _responseObj, MessageSuccess)
}

// Updates MobileId data in the system (for Issuer).
func (c *MobileIdChaincode) updateMobileId(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "updateMobileId", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_mobileIDReq := MobileId{}
	err := json.Unmarshal([]byte(args[0]), &_mobileIDReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateMobileIdFormat(_mobileIDReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}
	logger.Debug(_mobileIDReq.Dg1.MobileNo)

	// Precondition: The issuer field in the request must be equal to member's code.
	if _mobileIDReq.Dg1.Issuer != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request must be equal to member's code (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Precondition: The timestamp field must be in 24hr range.
	currentTimeStr, err := getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	difSec, err := compareTimeString(currentTimeStr, _mobileIDReq.Dg1.Timestamp)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if int64(math.Abs(float64(difSec))) > TimestampRange {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The timestamp field must be in 24hr range (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Precondition: MobileId must exist.
	_mobileIDKey, err := createMobileIdKey(stub, _mobileIDReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_mobileID, err := getMobileIDByNo(stub, _mobileIDReq.Dg1.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	// Precondition: The existing MobileId must not be terminated.
	if _mobileID.Dg1.Status == MobileIdTerminated {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The existing MobileId must not be terminated (mobile no: %s).", _mobileID.Dg1.MobileNo)
	}

	// Precondition: The issuer field in the request must be equal to the issuer field in the existing MobileId in the system.
	if _mobileIDReq.Dg1.Issuer != _mobileID.Dg1.Issuer {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request must be equal to the issuer field in the existing MobileId in the system (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Precondition: The serial number field in the request must be equal to the serial number field in the existing MobileId in the system.
	if _mobileIDReq.Dg1.MobileIdSn != _mobileID.Dg1.MobileIdSn {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The serial number field in the request must be equal to the serial number field in the existing MobileId in the system (mobile no: %s).", _mobileIDReq.Dg1.MobileNo)
	}

	// Marshal and save MobileId data into the system.
	_mobileIDBytes, err := json.Marshal(_mobileIDReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(_mobileIDKey, _mobileIDBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	_responseObj := getMessageSuccessResponse(stub.GetTxID())
	return shimSuccess(logger, _responseObj, MessageSuccess)
}

// Records Consent Log into the system (for Issuer).
func (c *MobileIdChaincode) recordConsentLog(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "recordConsentLog", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRecordOrGetConsentLog(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Precondition: ConsentLog must not exist in the system.
	_consentKey, err := createConsentIdKey(stub, _consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Get ConsentLog in the system.
	logger.Debug("Get ConsentLog in the system.")
	if b, _ := stub.GetState(_consentKey); b != nil {
		return shimErrorf(logger, StatusCodeInvalidOperation, "ConsentLog must not exist in the system (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	// Precondition: The specified MobileId must exist in the system.
	_mobileID, err := getMobileIDByNo(stub, _consentReq.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	// Precondition: The existing MobileId must be active.
	if _mobileID.Dg1.Status != MobileIdActive {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The existing MobileId must be active (mobile no: %s).", _mobileID.Dg1.MobileNo)
	}
	// Precondition: The serial no field in the request must be equal to the serial no field in the existing MobileId.
	if _consentReq.MobileIdSn != _mobileID.Dg1.MobileIdSn {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The serial no field in the request must be equal to the serial no field in the existing MobileId (mobile no: %s).", _consentReq.MobileNo)
	}
	// Precondition: The issuer field in the request, the issuer field in the existing MobileId, and the member's code must be the same.
	if _consentReq.Issuer != _member.MemberCode || _consentReq.Issuer != _mobileID.Dg1.Issuer || _mobileID.Dg1.Issuer != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request, the issuer field in the existing MobileId, and the member's code must be the same (mobile no: %s).", _consentReq.MobileNo)
	}

	logger.Info(fmt.Sprintf("[%s,%s,%s,%s]", _consentReq.MobileNo, _consentReq.MobileIdSn, _consentReq.Issuer, _consentReq.Cid))

	// Marshal ans save ConsentLog into the system.
	_consent := ConsentLog{}
	_consent.MobileNo = _consentReq.MobileNo
	_consent.MobileIdSn = _consentReq.MobileIdSn
	_consent.Issuer = _consentReq.Issuer
	_consent.Cid = _consentReq.Cid
	_consent.Created, err = getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, err := json.Marshal(_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	err = stub.PutState(_consentKey, _consentBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response.
	return shimSuccess(logger, _consentBytes, MessageSuccess)
}

// Retrieves MobileId and updates ConsentLog (for Verifier).
func (c *MobileIdChaincode) retrieveMobileIdAndUpdateConsentLog(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "retrieveMobileIdAndUpdateConsentLog", LogInfo)
	logger.Debug("Start retrieveMobileIdAndUpdateConsentLog")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberVerifier}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRetriveMobileId(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Precondition: ConsentLog must exist in the system.
	_consentKey, err := createConsentIdKey(stub, _consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, _ := stub.GetState(_consentKey)
	if _consentBytes == nil {
		return shimErrorf(logger, StatusCodeNotFound, "ConsentLog must exist in the system (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}
	_consent := ConsentLog{}
	err = json.Unmarshal(_consentBytes, &_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Precondition: The verifier field in the existing ConsentLog must be empty.
	if _consent.Verifier != "" {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The verifier field in the existing ConsentLog must be empty (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	// Precondition: The specified MobileId must exist in the system.
	_mobileID, err := getMobileIDByNo(stub, _consentReq.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	// Precondition: The existing MobileId must be active.
	if _mobileID.Dg1.Status != MobileIdActive {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The existing MobileId must be active (mobile no: %s).", _mobileID.Dg1.MobileNo)
	}
	// Precondition: The serial no field in the request must be equal to the serial no field in the existing MobileId.
	if _consentReq.MobileIdSn != _mobileID.Dg1.MobileIdSn {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The serial no field in the request must be equal to the serial no field in the existing MobileId (mobile no: %s).", _consentReq.MobileNo)
	}
	// Precondition: The issuer field in the request must be equal to the issuer field in the existing MobileId.
	if _consentReq.Issuer != _mobileID.Dg1.Issuer {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request must be equal to the issuer field in the existing MobileId (mobile no: %s).", _consentReq.MobileNo)
	}
	// Precondition: The verifier field in the request must be equal to the member's code.
	if _consentReq.Verifier != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The verifier field in the request must be equal to the member's code (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	logger.Info(fmt.Sprintf("Update Verifier Consent Log[%s, %s, %s, %s, %s, %s]", _consentReq.Verifier, _consentReq.TxType, _consentReq.Aal, _consentReq.Ref1, _consentReq.Ref2, _consentReq.RpId))

	// Update, marshal, and save ConsentLog into the system.
	_consent.Verifier = _consentReq.Verifier
	_consent.RpId = _consentReq.RpId
	_consent.Aal = _consentReq.Aal
	_consent.TxType = _consentReq.TxType
	_consent.Ref1 = _consentReq.Ref1
	_consent.Ref2 = _consentReq.Ref2
	_consent.Used, err = getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, err = json.Marshal(_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(_consentKey, _consentBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Get MobileId in the system
	logger.Info("Retrieve MobileId Data [" + _consentReq.MobileNo + "," + _consentReq.MobileIdSn + "," + _consentReq.Issuer + "]")
	_mobileIDBytes, err := json.Marshal(_mobileID)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _mobileIDBytes, MessageSuccess)
}

// Records verification result in Consent Log (for Verifier)
func (c *MobileIdChaincode) recordVerificationResult(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "recordVerificationResult", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberVerifier}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRecordResult(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Precondition: ConsentLog must exist in the system.
	_consent := ConsentLog{}
	_consent.MobileNo = _consentReq.MobileNo
	_consent.MobileIdSn = _consentReq.MobileIdSn
	_consent.Cid = _consentReq.Cid
	_consentKey, err := createConsentIdKey(stub, _consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, _ := stub.GetState(_consentKey)
	if _consentBytes == nil {
		return shimErrorf(logger, StatusCodeNotFound, "ConsentLog must exist in the system (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}
	err = json.Unmarshal(_consentBytes, &_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Precondition: The verifier field in the request must be equal to the verifier field in the existing ConsentLog.
	if _consent.Verifier != _consentReq.Verifier {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The verifier field in the request must be equal to the verifier field in the existing ConsentLog (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	// Precondition: The verified and face score fields in the existing ConsentLog must be empty.
	if _consent.Verified != "" || _consent.FaceScoreVerified != "" {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The verified and face score fields in the existing ConsentLog must be empty (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	// Precondition: The specified MobileId must exist in the system.
	_mobileID, err := getMobileIDByNo(stub, _consentReq.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	// Precondition: The existing MobileId must be active.
	if _mobileID.Dg1.Status != MobileIdActive {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The existing MobileId must be active (mobile no: %s).", _mobileID.Dg1.MobileNo)
	}
	// Precondition: The serial no field in the request must be equal to the serial no field in the existing MobileId.
	if _consentReq.MobileIdSn != _mobileID.Dg1.MobileIdSn {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The serial no field in the request must be equal to the serial no field in the existing MobileId (mobile no: %s).", _consentReq.MobileNo)
	}
	// Precondition: The verifier field in the request must be equal to the member's code.
	if _consentReq.Verifier != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The verifier field in the request must be equal to the member's code (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	// Update verification result
	_consent.Verified = _consentReq.Verified
	_consent.FaceScoreVerified = _consentReq.FaceScoreVerified
	_consentBytes, err = json.Marshal(_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(_consentKey, _consentBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	_responseObj := getMessageSuccessResponse(stub.GetTxID())
	return shimSuccess(logger, _responseObj, MessageSuccess)
}

// Gets Consent Log (For Issuer and Verifier)
func (c *MobileIdChaincode) getConsentLog(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getConsentLog", LogInfo)
	logger.Debug("Start getConsentLog")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer, MemberVerifier}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRecordOrGetConsentLog(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Precondition: ConsentLog must exist in the system.
	_consentKey, err := createConsentIdKey(stub, _consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, _ := stub.GetState(_consentKey)
	if _consentBytes == nil {
		return shimErrorf(logger, StatusCodeNotFound, "ConsentLog must exist in the system (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	logger.Info("[" + _consentReq.MobileNo + "," + _consentReq.Cid + "]")

	_consent := ConsentLog{}
	err = json.Unmarshal(_consentBytes, &_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// The verifier field or the issuer field in the ConsentLog must be equal to the member's code.
	if _consent.Verifier != _member.MemberCode && _consent.Issuer != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The verifier field or the issuer field in the ConsentLog must be equal to the member's code (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	// Return response
	return shimSuccess(logger, _consentBytes, MessageSuccess)
}

// *********************************************
// Audit Method
// *********************************************

// Gets MobileId for audit.
func (c *MobileIdChaincode) getMobileIdForAudit(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getMobileIdForAudit", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_mobileIDKeyReq := MobileIdKey{}
	err := json.Unmarshal([]byte(args[0]), &_mobileIDKeyReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdKeyFormat(_mobileIDKeyReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	logger.Info("Get getMobileId >> [" + _mobileIDKeyReq.MobileNo + "]")

	// Get MobileId in the system
	_mobileID, err := getMobileIDByNo(stub, _mobileIDKeyReq.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	_mobileIDBytes, err := json.Marshal(_mobileID)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _mobileIDBytes, MessageSuccess)
}

// Gets ConsentLog for audit.
func (c *MobileIdChaincode) getConsentLogForAudit(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getConsentLogForAudit", LogInfo)
	logger.Info("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRecordOrGetConsentLog(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// check if ConsentLog exists
	_consentKey, err := createConsentIdKey(stub, _consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, _ := stub.GetState(_consentKey)
	if _consentBytes == nil {
		return shimErrorf(logger, StatusCodeNotFound, "ConsentLog not found (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}

	_consentInfo := ConsentLog{}
	err = json.Unmarshal(_consentBytes, &_consentInfo)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if _consentReq.Issuer != _consentInfo.Issuer {
		return shimError(logger, StatusCodeInvalidArgument, "The issuer field in the requested data must be as same as the issuer field in the consent log.")
	}

	logger.Info("[" + _consentReq.MobileNo + "," + _consentReq.Cid + "]")

	// Return response
	return shimSuccess(logger, _consentBytes, MessageSuccess)
}

// Lists ConsentLog for audit.
func (c *MobileIdChaincode) listConsentLogForAudit(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listConsentLogForAudit", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentLogAuditQueryReq := ConsentLogAuditQuery{}
	err := json.Unmarshal([]byte(args[0]), &_consentLogAuditQueryReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateConsentLogAuditQueryFormat(_consentLogAuditQueryReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{"
	queryString += "\"_id\":{\"$regex\":\"mobileid_consent_key\"},"
	queryString += fmt.Sprintf("\"verifier\":\"%s\",", _consentLogAuditQueryReq.Verifier)
	queryString += fmt.Sprintf("\"used\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _consentLogAuditQueryReq.StartDate, _consentLogAuditQueryReq.EndDate)
	queryString += "}}"

	var _consents []ConsentLog
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	defer iterator.Close()
	for iterator.HasNext() {
		_kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _cs ConsentLog
		err = json.Unmarshal(_kv.Value, &_cs)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_consents = append(_consents, _cs)
	}

	pageSize, err := strconv.ParseInt(_consentLogAuditQueryReq.Records, 10, 32)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if len(_consents) > int(pageSize) {
		if _consentLogAuditQueryReq.RetrieveType == RetrieveTypeRandom {
			_idxList := makeRetrieveList(len(_consents), int(pageSize), RetrieveTypeRandom)
			tmpConsents := _consents
			_consents = nil
			for _, v := range _idxList {
				_consents = append(_consents, tmpConsents[v])
			}
		} else {
			_consents = _consents[0:int(pageSize)]
		}
	}

	_consentsBytes, err := json.Marshal(_consents)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _consentsBytes, MessageSuccess)
}

// Lists ConsentLog by issuer.
func (c *MobileIdChaincode) listConsentLogByIssuer(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listConsentLogByIssuer", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentLogQueryReq := ConsentLogQuery{}
	err := json.Unmarshal([]byte(args[0]), &_consentLogQueryReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateConsentLogQueryIssuerFormat(_consentLogQueryReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{"
	queryString += "\"_id\":{\"$regex\":\"mobileid_consent_key\"},"
	queryString += fmt.Sprintf("\"mobile_no\":\"%s\",", _consentLogQueryReq.MobileNo)
	queryString += fmt.Sprintf("\"mobile_id_sn\":\"%s\",", _consentLogQueryReq.MobileIdSn)
	queryString += fmt.Sprintf("\"issuer\":\"%s\",", _consentLogQueryReq.Issuer)
	queryString += fmt.Sprintf("\"created\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _consentLogQueryReq.StartDate, _consentLogQueryReq.EndDate)
	queryString += "}}"

	fmt.Println(queryString)

	var _consents []ConsentLog
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	defer iterator.Close()
	for iterator.HasNext() {
		_kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _cs ConsentLog
		err = json.Unmarshal(_kv.Value, &_cs)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_consents = append(_consents, _cs)
	}

	pageSize, err := strconv.ParseInt(_consentLogQueryReq.Records, 10, 32)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if len(_consents) > int(pageSize) {
		_consents = _consents[0:int(pageSize)]
	}

	_consentsBytes, err := json.Marshal(_consents)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _consentsBytes, MessageSuccess)
}

// Lists ConsentLog by verifier.
func (c *MobileIdChaincode) listConsentLogByVerifier(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listConsentLogByVerifier", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentLogQueryReq := ConsentLogQuery{}
	err := json.Unmarshal([]byte(args[0]), &_consentLogQueryReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateConsentLogQueryVerifierFormat(_consentLogQueryReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{"
	queryString += "\"_id\":{\"$regex\":\"mobileid_consent_key\"},"
	queryString += fmt.Sprintf("\"mobile_no\":\"%s\",", _consentLogQueryReq.MobileNo)
	queryString += fmt.Sprintf("\"mobile_id_sn\":\"%s\",", _consentLogQueryReq.MobileIdSn)
	queryString += fmt.Sprintf("\"verifier\":\"%s\",", _consentLogQueryReq.Verifier)
	queryString += fmt.Sprintf("\"used\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _consentLogQueryReq.StartDate, _consentLogQueryReq.EndDate)
	queryString += "}}"

	fmt.Println(queryString)

	var _consents []ConsentLog
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	defer iterator.Close()
	for iterator.HasNext() {
		_kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _cs ConsentLog
		err = json.Unmarshal(_kv.Value, &_cs)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_consents = append(_consents, _cs)
	}

	pageSize, err := strconv.ParseInt(_consentLogQueryReq.Records, 10, 32)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if len(_consents) > int(pageSize) {
		_consents = _consents[0:int(pageSize)]
	}

	_consentsBytes, err := json.Marshal(_consents)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _consentsBytes, MessageSuccess)
}

// Gets ConsentLog history.
func (c *MobileIdChaincode) getConsentHistory(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getConsentHistory", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRecordOrGetConsentLog(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	_consentKey, err := createConsentIdKey(stub, _consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	historyIter, err := stub.GetHistoryForKey(_consentKey)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	defer historyIter.Close()
	var _consents []ConsentLog
	for historyIter.HasNext() {
		_kv, err := historyIter.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _cs ConsentLog
		err = json.Unmarshal(_kv.Value, &_cs)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}

		if _consentReq.Issuer != _cs.Issuer {
			return shimError(logger, StatusCodeInvalidArgument, "The issuer field in the requested data must be as same as the issuer field in the consent log.")
		}

		logger.Info("TxId: " + _kv.TxId)
		logger.Info("Timestamp: " + _kv.Timestamp.String())
		logger.Info("Data: " + string(_kv.Value))
		_consents = append(_consents, _cs)
	}
	_consentsBytes, err := json.Marshal(&_consents)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return response
	return shimSuccess(logger, _consentsBytes, MessageSuccess)
}

// *********************************************
// Enhance-1
// *********************************************

// Gets MobileId for Issuer
func (c *MobileIdChaincode) getMobileId(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getMobileId", LogInfo)
	logger.Debug("Start getMobileId")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_mobileIDSnRequest := MobileIdSnRequest{}
	err := json.Unmarshal([]byte(args[0]), &_mobileIDSnRequest)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdSnRequestFormat(_mobileIDSnRequest); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{\"dg1\":{"
	queryString += fmt.Sprintf("\"issuer\":\"%s\"", _member.MemberCode)
	if _mobileIDSnRequest.MobileNo != "" {
		queryString += fmt.Sprintf(",\"mobile_no\":\"%s\"", _mobileIDSnRequest.MobileNo)
	}
	if _mobileIDSnRequest.MobileIdSn != "" {
		queryString += fmt.Sprintf(",\"mobile_id_sn\":\"%s\"", _mobileIDSnRequest.MobileIdSn)
	}
	queryString += "}},\"limit\":1}"

	logger.Info(queryString)

	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	var _mbs []MobileId
	defer iterator.Close()
	for iterator.HasNext() {
		_kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _mb MobileId
		err = json.Unmarshal(_kv.Value, &_mb)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_mbs = append(_mbs, _mb)
	}
	if len(_mbs) <= 0 {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	// Condition: The issuer field must be equal to member's code.
	if _mbs[0].Dg1.Issuer != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer of mobile no must be equal to member's code (mobile no: %s).", _mbs[0].Dg1.MobileNo)
	}

	// Lookup MobileId from world state
	_mobileID, err := getMobileIDByNo(stub, _mbs[0].Dg1.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	_mobileIDBytes, err := json.Marshal(_mobileID)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _mobileIDBytes, MessageSuccess)
}

// Count MobileId by Status for Issuer.
func (c *MobileIdChaincode) countMobileIdByStatus(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "countMobileIdByStatus", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_mobileIDStatusReq := MobileIdStatus{}
	err := json.Unmarshal([]byte(args[0]), &_mobileIDStatusReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdStatusFormat(_mobileIDStatusReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{\"dg1\":{"
	queryString += fmt.Sprintf("\"issuer\":\"%s\",", _member.MemberCode)
	queryString += fmt.Sprintf("\"status\":\"%s\"", _mobileIDStatusReq.Status)
	queryString += "}}}"

	logger.Info(queryString)

	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	_countMobileID := 0
	defer iterator.Close()
	for iterator.HasNext() {
		_, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_countMobileID++
	}

	var _returnResp CountMobileIdStatusResponse
	_returnResp.Status = _mobileIDStatusReq.Status
	_returnResp.Count = strconv.Itoa(_countMobileID)
	_returnRespBytes, err := json.Marshal(_returnResp)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _returnRespBytes, MessageSuccess)
}

// Lists ConsentLog for Issuer.
func (c *MobileIdChaincode) listConsentLogByDate(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listConsentLogByDate", LogInfo)
	logger.Debug("Start listConsentLogByDate")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer, MemberVerifier}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentLogQueryByDateReq := ConsentLogQueryByDate{}
	err := json.Unmarshal([]byte(args[0]), &_consentLogQueryByDateReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateConsentLogQueryByDateFormat(_consentLogQueryByDateReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{"
	queryString += "\"_id\":{\"$regex\":\"mobileid_consent_key\"}"
	if strings.Contains(_member.MemberRole, MemberIssuer) {
		queryString += fmt.Sprintf(",\"issuer\":\"%s\"", _member.MemberCode)
	} else if strings.Contains(_member.MemberRole, MemberVerifier) {
		queryString += fmt.Sprintf(",\"verifier\":\"%s\"", _member.MemberCode)
	} else {
		return shimError(logger, StatusCodeInternalError, "MemberRole failed")
	}

	_dtCondCount := 0
	var queryDateCond string
	if _consentLogQueryByDateReq.CreatedStartDate != "" && _consentLogQueryByDateReq.CreatedEndDate != "" {
		queryDateCond += fmt.Sprintf("\"created\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _consentLogQueryByDateReq.CreatedStartDate, _consentLogQueryByDateReq.CreatedEndDate)
		_dtCondCount++
	}
	if _consentLogQueryByDateReq.UsedStartDate != "" && _consentLogQueryByDateReq.UsedEndDate != "" {
		if _dtCondCount > 0 {
			if _dtCondCount == 1 {
				queryDateCond = fmt.Sprintf("{%s}", queryDateCond)
			}
			queryDateCond += fmt.Sprintf(",{\"used\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}}", _consentLogQueryByDateReq.UsedStartDate, _consentLogQueryByDateReq.UsedEndDate)
		} else {
			queryDateCond += fmt.Sprintf("\"used\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _consentLogQueryByDateReq.UsedStartDate, _consentLogQueryByDateReq.UsedEndDate)
		}
		_dtCondCount++
	}
	if _consentLogQueryByDateReq.RevokedStartDate != "" && _consentLogQueryByDateReq.RevokedEndDate != "" {
		if _dtCondCount > 0 {
			if _dtCondCount == 1 {
				queryDateCond = fmt.Sprintf("{%s}", queryDateCond)
			}
			queryDateCond += fmt.Sprintf(",{\"revoked\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}}", _consentLogQueryByDateReq.RevokedStartDate, _consentLogQueryByDateReq.RevokedEndDate)
		} else {
			queryDateCond += fmt.Sprintf("\"revoked\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _consentLogQueryByDateReq.RevokedStartDate, _consentLogQueryByDateReq.RevokedEndDate)
		}
		_dtCondCount++
	}
	if _dtCondCount > 1 {
		queryString += ",\"$or\": ["
		queryString += queryDateCond
		queryString += "]"
	} else {
		if _dtCondCount == 1 {
			queryString += "," + queryDateCond
		}
	}

	queryString += "}"
	// queryString += fmt.Sprintf(",\"limit\":%s", _consentLogQueryByDateReq.Records)
	queryString += "}"
	logger.Info(queryString)

	var _consents []ConsentLog
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	defer iterator.Close()
	for iterator.HasNext() {
		_kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _cs ConsentLog
		err = json.Unmarshal(_kv.Value, &_cs)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_consents = append(_consents, _cs)
	}

	// use limit option in query
	pageSize, err := strconv.ParseInt(_consentLogQueryByDateReq.Records, 10, 32)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if len(_consents) > int(pageSize) {
		_consents = _consents[0:int(pageSize)]
	}

	_consentsBytes, err := json.Marshal(_consents)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _consentsBytes, MessageSuccess)
}

// Revoke unused ConsentLog for Issuer.
func (c *MobileIdChaincode) revokeConsentLog(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "revokeConsentLog", LogInfo)
	logger.Debug("Start")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_consentReq := ConsentLog{}
	err := json.Unmarshal([]byte(args[0]), &_consentReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if err = validateConsentLogFormatForRevokeConsent(_consentReq); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Precondition: ConsentLog must exist in the system.
	_consent := ConsentLog{}
	_consent.MobileNo = _consentReq.MobileNo
	_consent.MobileIdSn = _consentReq.MobileIdSn
	_consent.Cid = _consentReq.Cid
	_consentKey, err := createConsentIdKey(stub, _consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_consentBytes, _ := stub.GetState(_consentKey)
	if _consentBytes == nil {
		return shimErrorf(logger, StatusCodeNotFound, "ConsentLog must exist in the system (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}
	err = json.Unmarshal(_consentBytes, &_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Precondition: The verifier field in the request must be equal to the verifier field in the existing ConsentLog.
	if _consent.Issuer != _consentReq.Issuer {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request must be equal to the issuer field in the existing ConsentLog (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}
	// Precondition: The verified and face score fields in the existing ConsentLog must be empty.
	if _consent.Used != "" || _consent.Verifier != "" {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The ConsentLog is in used (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}
	// Precondition: The specified MobileId must exist in the system.
	_mobileID, err := getMobileIDByNo(stub, _consentReq.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	// Precondition: The Issuer field in the mobile-id must be equal to the member's code.
	if _mobileID.Dg1.Issuer != _member.MemberCode {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The issuer field in the request must be equal to the member's code (cid: %s, mobile no: %s).", _consentReq.Cid, _consentReq.MobileNo)
	}
	// Precondition: The serial no field in the request must be equal to the serial no field in the existing MobileId.
	if _consentReq.MobileIdSn != _mobileID.Dg1.MobileIdSn {
		return shimErrorf(logger, StatusCodeInvalidOperation, "The serial no field in the request must be equal to the serial no field in the existing MobileId (mobile no: %s).", _consentReq.MobileNo)
	}

	// Update revoked date
	if _consentReq.Revoked != "" {
		_consent.Revoked = _consentReq.Revoked
	} else {
		_consent.Revoked, err = getCurrentDateTimeString()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
	}
	_consentBytes, err = json.Marshal(_consent)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(_consentKey, _consentBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	_responseObj := getMessageSuccessResponse(stub.GetTxID())
	return shimSuccess(logger, _responseObj, MessageSuccess)
}

// *********************************************
// MobileId 1.5 Extend
// *********************************************

// Retrieve MobileId Issuer for Verifier Node and Insert Request Logs
func (c *MobileIdChaincode) retrieveMobileIdIssuer(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "retrieveMobileIdIssuer", LogInfo)
	logger.Debug("Start retrieveMobileIdIssuer")

	// 1.Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberVerifier}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// 2.Unmarshal and validate request.
	_mobileIDIssuerRequest := MobileIdIssuerRequest{}
	err := json.Unmarshal([]byte(args[0]), &_mobileIDIssuerRequest)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdIssuerRequestFormat(_mobileIDIssuerRequest); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// 3.1.Lookup MobileId
	// Precondition: The specified MobileId must exist in the system.
	_mobileID, err := getMobileIDByNo(stub, _mobileIDIssuerRequest.MobileNo)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	//** return status to caller
	// Precondition: The existing MobileId must be active.
	//if _mobileID.Dg1.Status != MobileIdActive {
	//	return shimErrorf(logger, StatusCodeInvalidOperation, "The existing MobileId must be active (mobile no: %s).", _mobileID.Dg1.MobileNo)
	//}

	// 3.2.Lookup Member Information
	_memberReq := Member{}
	_memberReq.MemberCode = _mobileID.Dg1.Issuer

	// Precondition: Member must exist.
	memKey, err := createMemberKey(stub, _memberReq)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	b, err := stub.GetState(memKey)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}
	if b == nil {
		return shimError(logger, StatusCodeInvalidOperation, "Member must exist.")
	}
	_memberInfo := Member{}
	err = json.Unmarshal(b, &_memberInfo)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// 4.1.Prepare Model for Response
	_mobIssuer := MobileIdIssuer{}
	_mobIssuer.MobileNo = _mobileID.Dg1.MobileNo
	_mobIssuer.Issuer = _mobileID.Dg1.Issuer
	_mobIssuer.MobileIdSn = _mobileID.Dg1.MobileIdSn
	_mobIssuer.Status = _mobileID.Dg1.Status
	_mobIssuer.ServiceUrl = _memberInfo.ServiceUrl // find from issuer member

	// 4.2.Prepare Model for Record Request Log
	_mobRequestLog := MobileIdRequestLog{}
	_mobRequestLog.MobileNo = _mobileID.Dg1.MobileNo
	_mobRequestLog.Issuer = _mobileID.Dg1.Issuer
	_mobRequestLog.MobileIdSn = _mobileID.Dg1.MobileIdSn
	_mobRequestLog.Verifier = _member.MemberCode
	_mobRequestLog.RpId = _mobileIDIssuerRequest.RpId
	_mobRequestLog.TxId = _mobileIDIssuerRequest.TxId
	_mobRequestLog.Status = _mobileID.Dg1.Status
	_mobRequestLog.Aal = _mobileIDIssuerRequest.Aal
	_mobRequestLog.Created, err = getCurrentDateTimeString()
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// 5.Insert Request Consent Log
	// Precondition: RequestLog must exist in the system.
	_requestLogKey, err := createMobileIdRequestKey(stub, _mobRequestLog)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	_requestLogBytes, err := json.Marshal(_mobRequestLog)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	err = stub.PutState(_requestLogKey, _requestLogBytes)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// 6.Return MobileIdIssuer
	_mobIssuerBytes, err := json.Marshal(_mobIssuer)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _mobIssuerBytes, MessageSuccess)
}

// List RequestLog
func (c *MobileIdChaincode) listRequestLogByDate(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "listRequestLogByDate", LogInfo)
	logger.Debug("Start listRequestLogByDate")

	// Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer, MemberVerifier, MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// Unmarshal and validate request.
	_mbRequestLogQry := MobileIdRequestLogQueryByDate{}
	err := json.Unmarshal([]byte(args[0]), &_mbRequestLogQry)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdRequestLogQueryByDateFormat(_mbRequestLogQry); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	// Query
	var queryString string
	queryString += "{\"selector\":{"
	// mobileid_request_key
	queryString += "\"_id\": {\"$regex\": \"mobileid_request_key\" }"

	if strings.Contains(_member.MemberRole, MemberIssuer) {
		//queryString += "{\"selector\":{"
		queryString += fmt.Sprintf(",\"issuer\":\"%s\"", _member.MemberCode)
	} else if strings.Contains(_member.MemberRole, MemberVerifier) {
		//queryString += "{\"selector\":{"
		queryString += fmt.Sprintf(",\"verifier\":\"%s\"", _member.MemberCode)
	} else if strings.Contains(_member.MemberRole, MemberRegulator) {
		//
	} else {
		return shimError(logger, StatusCodeInternalError, "MemberRole failed")
	}

	if _mbRequestLogQry.StartDate != "" && _mbRequestLogQry.EndDate != "" {
		queryString += fmt.Sprintf(",\"created\":{\"$gte\": \"%s\",\"$lte\": \"%s\"}", _mbRequestLogQry.StartDate, _mbRequestLogQry.EndDate)
	}
	queryString += "}"
	// queryString += fmt.Sprintf(",\"limit\":%s", _mbRequestLogQry.Records)
	queryString += "}"
	logger.Info(queryString)

	var _mbRequestLogs []MobileIdRequestLog
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shimError(logger, StatusCodeNotFound, err.Error())
	}

	defer iterator.Close()
	for iterator.HasNext() {
		_kv, err := iterator.Next()
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		var _req MobileIdRequestLog
		err = json.Unmarshal(_kv.Value, &_req)
		if err != nil {
			return shimError(logger, StatusCodeInternalError, err.Error())
		}
		_mbRequestLogs = append(_mbRequestLogs, _req)
	}

	// use limit option in query
	pageSize, err := strconv.ParseInt(_mbRequestLogQry.Records, 10, 32)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	if len(_mbRequestLogs) > int(pageSize) {
		_mbRequestLogs = _mbRequestLogs[0:int(pageSize)]
	}

	_mbRequestLogsBytes, err := json.Marshal(_mbRequestLogs)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	// Return response
	return shimSuccess(logger, _mbRequestLogsBytes, MessageSuccess)
}

// get MobileId Request Log By Id
func (c *MobileIdChaincode) getMobileIdRequestLog(stub shim.ChaincodeStubInterface, _member Member, args []string, loguuid string) sc.Response {
	logger := newShimLogger(loguuid, "getMobileIdRequestLog", LogInfo)
	logger.Debug("Start getMobileIdRequestLog")

	// 1.Validate authorization.
	if err := validateMemberAuthorization(_member, []string{MemberIssuer, MemberVerifier, MemberRegulator}); err != nil {
		return shimError(logger, StatusCodeUnauthorized, err.Error())
	}

	// 2.Unmarshal and validate request.
	_midRequestLog := MobileIdRequestLogKey{}
	err := json.Unmarshal([]byte(args[0]), &_midRequestLog)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}

	if err = validateMobileIdRequestLogFormat(_midRequestLog); err != nil {
		return shimError(logger, StatusCodeInvalidArgument, err.Error())
	}

	_requestLog := MobileIdRequestLog{}
	_requestLog.MobileNo = _midRequestLog.MobileNo
	_requestLog.RpId = _midRequestLog.RpId
	_requestLog.TxId = _midRequestLog.TxId

	// Precondition: ConsentLog must exist in the system.
	_requestLogKey, err := createMobileIdRequestKey(stub, _requestLog)
	if err != nil {
		return shimError(logger, StatusCodeInternalError, err.Error())
	}
	_requestLogBytes, _ := stub.GetState(_requestLogKey)
	if _requestLogBytes == nil {
		return shimErrorf(logger, StatusCodeNotFound, "Request Log must exist in the system (rp: %s, mobile no: %s).", _midRequestLog.RpId, _midRequestLog.MobileNo)
	}

	logger.Info("MobileIdRequestLog[" + _midRequestLog.MobileNo + "," + _midRequestLog.RpId + "]")

	// Return response
	return shimSuccess(logger, _requestLogBytes, MessageSuccess)
}

// *********************************************
// Validate methods
// *********************************************

func validateMemberAuthorization(_member Member, allowedRoles []string) error {
	for _, v := range allowedRoles {
		if strings.Contains(_member.MemberRole, v) && _member.Status == MemberActive {
			return nil
		}
	}
	return errors.New(MessageUnauthorized)
}

func validateString(name string, str string, minLength int, maxLength int, regexPattern string) error {
	if len(str) < minLength {
		return errors.New(name + ": Too short")
	}
	if len(str) > maxLength {
		return errors.New(name + ": Too long")
	}
	if !regexp.MustCompile(regexPattern).MatchString(str) {
		return errors.New(name + ": Pattern mismatch")
	}
	return nil
}

func validateInt(name string, str string, min int, max int) error {
	num, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return errors.New(name + ": Parse Int Failed")
	}
	if int(num) < min {
		return errors.New(name + ": Too low")
	}
	if int(num) > max {
		return errors.New(name + ": Too high")
	}
	return nil
}

func validateRecords(name string, str string, min int, max int) error {
	if err := validateString(name, str, min, max, `^[0-9]*$`); err != nil {
		return err
	}
	num, err2 := strconv.ParseInt(str, 10, 32)
	if err2 != nil {
		return errors.New(name + ": Parse Int Failed")
	}
	if int(num) < min {
		return errors.New(name + ": Too low")
	}
	if int(num) > max {
		return errors.New(name + ": Too high")
	}
	return nil
}

func validateDecimal(name string, str string, min float64, max float64) error {
	numDec, err := decimal.NewFromString(str)
	if err != nil {
		return errors.New(name + ": Parse Float Failed")
	}
	minDec := decimal.NewFromFloat(min)
	maxDec := decimal.NewFromFloat(max)
	if numDec.LessThan(minDec) {
		return errors.New(name + ": Too low")
	}
	if numDec.GreaterThan(maxDec) {
		return errors.New(name + ": Too high")
	}
	return nil
}

// Validate mobile_no
func validateMobileNo(name string, val string) error {
	//minLength, maxLength = 9, 12
	//regexPattern = `^[0-9]*$`
	//validateString(name, val, 9, 12, `^[0-9]*$`);
	if err := validateString(name, val, 44, 44, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	return nil
}

// Validate mobile_id_sn
func validateMobileIDSN(name string, val string) error {
	//minLength, maxLength = 3, 12
	//regexPattern = `^[0-9]*$`
	//validateString(name, val, 3, 12, `^[0-9]*$`);
	if err := validateString(name, val, 44, 44, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	return nil
}

// Validate datetime value
func validateDateTimeValue(dtVal string) error {
	layout := "2006-01-02T15:04:05"
	_, err := time.Parse(layout, dtVal)
	if err != nil {
		//fmt.Println("ERROR:",err)
		return err
	}
	return nil
}

// Validate datetime format
func validateDateTimeFormat(name string, val string) error {
	//minLength, maxLength = 19, 19
	//regexPattern = `^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])T(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(\\.[0-9]+)?(Z)?$`
	//validateString(name, val, 19, 19, `^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])T(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(\\.[0-9]+)?(Z)?$`);
	if err := validateString(name, val, 19, 19, `^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])T(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(\\.[0-9]+)?(Z)?$`); err != nil {
		return err
	}
	if err := validateDateTimeValue(val); err != nil {
		return err
	}

	return nil
}

// Validate member_code
func validateMemberCode(name string, val string, _devMode bool) error {
	//minLength, maxLength = 3, 16
	//regexPattern = `^[0-9A-Za-z]{3,16}$`
	//validateString(name, val, 3, 16 `^[0-9A-Za-z]{3,16}$`);
	if _devMode {
		if err := validateString(name, val, 3, 16, `^[0-9A-Za-z]{3,16}$`); err != nil {
			return err
		}
	} else {
		if err := validateString(name, val, 3, 3, `^[0-9A-Za-z]{3}$`); err != nil {
			return err
		}
	}

	return nil
}

// Validate Serial_No_Prefix
func validateSerialNoPrefix(name string, val string) error {
	//minLength, maxLength = 2, 2,
	//regexPattern = `^[0-9]*$`
	//validateString(name, val, 2, 2, `^[0-9]*$`);
	if err := validateString(name, val, 2, 2, `^[0-9]*$`); err != nil {
		return err
	}
	return nil
}

func validateMobileIdFormat(_req MobileId) error {
	if err := validateMobileNo("dg1.mobile_no", _req.Dg1.MobileNo); err != nil {
		return err
	}
	if err := validateString("dg1.issuer", _req.Dg1.Issuer, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateMobileIDSN("dg1.mobile_id_sn", _req.Dg1.MobileIdSn); err != nil {
		return err
	}
	if err := validateString("dg1.ial", _req.Dg1.Ial, 3, 3, `^\d\.\d$`); err != nil {
		return err
	}
	if err := validateDecimal("dg1.ial", _req.Dg1.Ial, 1.0, 3.0); err != nil {
		return err
	}

	if err := validateString("dg1.status", _req.Dg1.Status, 1, 1, `A|S|T`); err != nil {
		return err
	}
	if err := validateString("dg1.face_engine_id", _req.Dg1.FaceEngineId, 8, 8, `^[0-9A-Z]+$`); err != nil {
		return err
	}

	if err := validateDateTimeFormat("dg1.timestamp", _req.Dg1.Timestamp); err != nil {
		return err
	}
	if err := validateString("dg2", _req.Dg2, 1, 1000000, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}

	if err := validateString("sod.hdg1", _req.Sod.Hdg1, 44, 44, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	if err := validateString("sod.hdg2", _req.Sod.Hdg2, 44, 44, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	if err := validateString("sod.hface_template", _req.Sod.HfaceTemplate, 44, 44, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	if err := validateString("sod.dig", _req.Sod.Dig, 44, 44, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	if err := validateString("sod.sig", _req.Sod.Sig, 344, 344, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	if err := validateString("sod.cert", _req.Sod.Cert, 1, 10000, `^[a-zA-Z0-9+/]+={0,2}$`); err != nil {
		return err
	}
	return nil
}

func validateConsentLogFormatForRecordOrGetConsentLog(_req ConsentLog) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateString("issuer", _req.Issuer, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
		return err
	}
	if err := validateString("cid", _req.Cid, 36, 36, `^(\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\}{0,1})$`); err != nil {
		return err
	}
	return nil
}

func validateConsentLogFormatForRetriveMobileId(_req ConsentLog) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateString("issuer", _req.Issuer, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
		return err
	}
	if err := validateString("cid", _req.Cid, 36, 36, `^(\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\}{0,1})$`); err != nil {
		return err
	}

	if err := validateString("verifier", _req.Verifier, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateString("aal", _req.Aal, 3, 3, `^\d\.\d$`); err != nil {
		return err
	}
	if err := validateDecimal("aal", _req.Aal, 1.0, 3.0); err != nil {
		return err
	}
	if err := validateString("tx_type", _req.TxType, 1, 5, `^[A-Z]+$`); err != nil {
		return err
	}
	if err := validateString("ref1", _req.Ref1, 1, 20, `.*`); err != nil {
		return err
	}
	if err := validateString("ref2", _req.Ref2, 0, 20, `.*`); err != nil {
		return err
	}
	return nil
}

func validateConsentLogFormatForRecordResult(_req ConsentLog) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
		return err
	}
	if err := validateString("cid", _req.Cid, 36, 36, `^(\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\}{0,1})$`); err != nil {
		return err
	}
	if err := validateString("verifier", _req.Verifier, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateString("verified", _req.Verified, 1, 1, `Y|N`); err != nil {
		return err
	}
	if err := validateString("face_score_verified", _req.FaceScoreVerified, 0, 10, `.`); err != nil {
		return err
	}
	if err := validateDecimal("face_score_verified", _req.FaceScoreVerified, 0, 9999999999); err != nil {
		return err
	}
	// if err := validateString("revoked", _req.Revoked, 19, 19, `^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])T(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(\\.[0-9]+)?(Z)?$`); err != nil {
	// 	return err
	// }
	return nil
}

func validateMember(_req Member, _devMode bool) error {
	if err := validateMemberCode("member_code", _req.MemberCode, _devMode); err != nil {
		return err
	}
	if err := validateString("member_name", _req.MemberName, 1, 100, `.`); err != nil {
		return err
	}
	if err := validateString("member_role", _req.MemberRole, 1, 25, `^((ISSUER|VERIFIER|REGULATOR),?){1,3}$`); err != nil {
		return err
	}
	if err := validateString("status", _req.Status, 1, 1, `A|T`); err != nil {
		return err
	}
	return nil
}

func validateMobileIdKeyFormat(_req MobileIdKey) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	return nil
}

func validateConsentLogAuditQueryFormat(_req ConsentLogAuditQuery) error {
	if err := validateString("verifier", _req.Verifier, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateRecords("records", _req.Records, 1, 100); err != nil {
		return err
	}
	if err := validateString("retrieve_type", _req.RetrieveType, 1, 10, `^RANDOM|NORMAL$`); err != nil {
		return err
	}
	if err := validateDateTimeFormat("start_date", _req.StartDate); err != nil {
		return err
	}
	if err := validateDateTimeFormat("end_date", _req.EndDate); err != nil {
		return err
	}
	return nil
}

func validateConsentLogQueryIssuerFormat(_req ConsentLogQuery) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateString("issuer", _req.Issuer, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
		return err
	}
	if err := validateRecords("records", _req.Records, 1, 100); err != nil {
		return err
	}
	if err := validateDateTimeFormat("start_date", _req.StartDate); err != nil {
		return err
	}
	if err := validateDateTimeFormat("end_date", _req.EndDate); err != nil {
		return err
	}
	return nil
}

func validateConsentLogQueryVerifierFormat(_req ConsentLogQuery) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
		return err
	}
	if err := validateString("verifier", _req.Verifier, 3, 3, `^[0-9A-Z]{3}$`); err != nil {
		return err
	}
	if err := validateRecords("records", _req.Records, 1, 100); err != nil {
		return err
	}
	if err := validateDateTimeFormat("start_date", _req.StartDate); err != nil {
		return err
	}
	if err := validateDateTimeFormat("end_date", _req.EndDate); err != nil {
		return err
	}
	return nil
}

func validateMobileIdSnRequestFormat(_req MobileIdSnRequest) error {
	_reqValid := false
	if _req.MobileNo != "" {
		if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
			return err
		}
		_reqValid = true
	}

	if _req.MobileIdSn != "" {
		if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
			return err
		}
		_reqValid = true
	}
	if !_reqValid {
		return errors.New("Incorrect request data")
	}
	return nil
}

func validateMobileIdStatusFormat(_req MobileIdStatus) error {
	if err := validateString("dg1.status", _req.Status, 1, 1, `A|S|T`); err != nil {
		return err
	}
	return nil
}

func validateConsentLogQueryByDateFormat(_req ConsentLogQueryByDate) error {
	_queryValid := false
	if _req.CreatedStartDate != "" && _req.CreatedEndDate != "" {
		if err := validateDateTimeFormat("created_start_date", _req.CreatedStartDate); err != nil {
			return err
		}
		if err := validateDateTimeFormat("created_end_date", _req.CreatedEndDate); err != nil {
			return err
		}
		_queryValid = true
	} else if _req.CreatedStartDate != "" || _req.CreatedEndDate != "" {
		return errors.New("created_start_date and created_end_date should be completed")
	}

	if _req.UsedStartDate != "" && _req.UsedEndDate != "" {
		if err := validateDateTimeFormat("used_start_date", _req.UsedStartDate); err != nil {
			return err
		}
		if err := validateDateTimeFormat("used_end_date", _req.UsedEndDate); err != nil {
			return err
		}
		_queryValid = true
	} else if _req.UsedStartDate != "" || _req.UsedEndDate != "" {
		return errors.New("used_start_date and used_end_date should be completed")
	}

	if _req.RevokedStartDate != "" && _req.RevokedEndDate != "" {
		if err := validateDateTimeFormat("revoked_start_date", _req.RevokedStartDate); err != nil {
			return err
		}
		if err := validateDateTimeFormat("revoked_end_date", _req.RevokedEndDate); err != nil {
			return err
		}
		_queryValid = true
	} else if _req.RevokedStartDate != "" || _req.RevokedEndDate != "" {
		return errors.New("revoked_start_date and revoked_end_date should be completed")
	}

	if !_queryValid {
		return errors.New("Invalid request data")
	}
	//page_size > records
	if err := validateRecords("records", _req.Records, 1, 100); err != nil {
		return err
	}
	return nil
}

func validateConsentLogFormatForRevokeConsent(_req ConsentLog) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateMobileIDSN("mobile_id_sn", _req.MobileIdSn); err != nil {
		return err
	}
	if err := validateString("cid", _req.Cid, 36, 36, `^(\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\}{0,1})$`); err != nil {
		return err
	}
	//check revoked date if it have value.
	if _req.Revoked != "" {
		if err := validateDateTimeFormat("revoked", _req.Revoked); err != nil {
			return err
		}
	}
	return nil
}

// validateMobileIdIssuerRequestFormat
func validateMobileIdIssuerRequestFormat(_req MobileIdIssuerRequest) error {
	_reqValid := false

	if _req.RpId != "" {
		if err := validateString("rp_id", _req.RpId, 3, 16, `^[0-9A-Z]{3,16}$`); err != nil {
			return err
		}
		_reqValid = true
	}
	if _req.TxId != "" {
		//if err := validateString("tx_id", _req.TxId, 3, 36, `^(\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\}{0,1})$`); err != nil {
		if err := validateString("tx_id", _req.TxId, 3, 36, `^[0-9a-zA-Z]{3,36}$`); err != nil {
			return err
		}
		_reqValid = true
	}
	if _req.MobileNo != "" {
		if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
			return err
		}
		_reqValid = true
	}
	if _req.Aal != "" {
		if err := validateDecimal("aal", _req.Aal, 1.0, 3.0); err != nil {
			return err
		}
		_reqValid = true
	}

	if !_reqValid {
		return errors.New("Invalid request data")
	}
	return nil
}

func validateMobileIdRequestLogQueryByDateFormat(_req MobileIdRequestLogQueryByDate) error {
	_queryValid := false
	if _req.StartDate != "" && _req.EndDate != "" {
		if err := validateDateTimeFormat("created_start_date", _req.StartDate); err != nil {
			return err
		}
		if err := validateDateTimeFormat("created_end_date", _req.EndDate); err != nil {
			return err
		}
		_queryValid = true
	} else if _req.StartDate != "" || _req.EndDate != "" {
		return errors.New("start date and end date should be completed")
	}

	if !_queryValid {
		return errors.New("Invalid date format")
	}
	//page_size > records
	if err := validateRecords("records", _req.Records, 1, 100); err != nil {
		return err
	}
	return nil
}

func validateMobileIdRequestLogFormat(_req MobileIdRequestLogKey) error {
	if err := validateMobileNo("mobile_no", _req.MobileNo); err != nil {
		return err
	}
	if err := validateString("rp_id", _req.RpId, 3, 16, `^[0-9A-Z]{3,16}$`); err != nil {
		return err
	}
	//if err := validateString("tx_id", _req.TxId, 3, 36, `^(\{{0,1}([0-9a-fA-F]){8}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){4}-([0-9a-fA-F]){12}\}{0,1})$`); err != nil {
	if err := validateString("tx_id", _req.TxId, 3, 36, `^[0-9a-zA-Z]{3,36}$`); err != nil {
		return err
	}
	return nil
}

// *********************************************
// Get data from world state
// *********************************************

func mobileIdExists(stub shim.ChaincodeStubInterface, mobileNo string) bool {
	_mobile := MobileId{}
	_dg1 := DataGroup1{}
	_dg1.MobileNo = mobileNo
	_mobile.Dg1 = _dg1

	_mbIDKey, err := createMobileIdKey(stub, _mobile)
	if err != nil {
		return false
	}

	_mbByte, err := stub.GetState(_mbIDKey)
	if err != nil {
		return false
	}
	if _mbByte == nil {
		return false
	}

	return true
}

// Change getMobileId to getMobileIDByNo
func getMobileIDByNo(stub shim.ChaincodeStubInterface, mobileNo string) (MobileId, error) {
	_mobile := MobileId{}
	_dg1 := DataGroup1{}
	_dg1.MobileNo = mobileNo
	_mobile.Dg1 = _dg1

	_mbIDKey, err := createMobileIdKey(stub, _mobile)
	if err != nil {
		return _mobile, err
	}

	_mbByte, err := stub.GetState(_mbIDKey)
	if err != nil {
		return _mobile, err
	}

	err = json.Unmarshal(_mbByte, &_mobile)
	if err != nil {
		return _mobile, err
	}

	return _mobile, err
}

// *********************************************
// Utils
// *********************************************
func newShimLogger(uuid string, name string, level LoggingLevel) *ChaincodeLogger {
	logger := NewLogger(fmt.Sprintf("%s %s", uuid, name))
	logger.SetLevel(level)
	return logger
}

func shimSuccess(logger *ChaincodeLogger, payload []byte, msg string) sc.Response {
	logger.Info(fmt.Sprintf("[Code %d] %s", StatusCodeSuccess, msg))
	return sc.Response{
		Status:  int32(StatusCodeSuccess),
		Payload: payload,
	}
}

func shimError(logger *ChaincodeLogger, statusCode int, msg string) sc.Response {
	msg = fmt.Sprintf("%s: %s", StatusCodeMap[statusCode], msg)
	logger.Error(fmt.Sprintf("[Code %d] %s", statusCode, msg))
	return sc.Response{
		Status:  int32(statusCode),
		Message: msg,
	}
}

func shimErrorf(logger *ChaincodeLogger, statusCode int, format string, args ...interface{}) sc.Response {
	msg := fmt.Sprintf(format, args...)
	msg = fmt.Sprintf("%s: %s", StatusCodeMap[statusCode], msg)
	logger.Error(fmt.Sprintf("[Code %d] %s", statusCode, msg))
	return sc.Response{
		Status:  int32(statusCode),
		Message: msg,
	}
}

func getMessageSuccessResponse(_code string) []byte {
	msgObj := OperationResponse{}
	msgObj.Code = _code
	msgObj.Message = MessageSuccess
	//return ([]byte(msgObj))
	//return msgObj
	byteObj, _ := json.Marshal(msgObj)
	return byteObj
}

func getCurrentDateTimeString() (string, error) {
	timeLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return "", err
	}
	return time.Now().In(timeLocation).Format(TimeFormat), nil
}

func compareTimeString(time1 string, time2 string) (int64, error) {
	ts1, err := time.Parse(TimeFormat, time1)
	if err != nil {
		return 0, err
	}
	ts2, err := time.Parse(TimeFormat, time2)
	if err != nil {
		return 0, err
	}

	return (ts1.Unix() - ts2.Unix()), nil
}

func makeRetrieveList(inputSize int, targetSize int, retrieveType string) []int {
	a := make([]int, inputSize)
	for i := range a {
		a[i] = i
	}

	if retrieveType == RetrieveTypeRandom {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	}
	b := a[0:targetSize]
	sort.Ints(b)
	return b
}

//*******************************************************
// Data Structure
//*******************************************************

// DataGroup1 is a type represent a DG1
type DataGroup1 struct {
	MobileNo     string `json:"mobile_no"`
	Issuer       string `json:"issuer"`
	MobileIdSn   string `json:"mobile_id_sn"`
	Ial          string `json:"ial"`
	Status       string `json:"status"`
	FaceEngineId string `json:"face_engine_id"`
	Timestamp    string `json:"timestamp"`
}

// SODDataGroup Data is a type represent a SOD
type SODDataGroup struct {
	Hdg1          string `json:"hdg1"`
	Hdg2          string `json:"hdg2"`
	HfaceTemplate string `json:"hface_template"`
	Dig           string `json:"dig"`
	Sig           string `json:"sig"`
	Cert          string `json:"cert"`
}

// MobileId is a type represent a MobileId Identifier Information in blockchain
type MobileId struct {
	Dg1 DataGroup1   `json:"dg1"`
	Dg2 string       `json:"dg2"`
	Sod SODDataGroup `json:"sod"`
}

// ConsentLog is a type represent a Consent Information in blockchain
type ConsentLog struct {
	MobileNo          string `json:"mobile_no"`
	Issuer            string `json:"issuer"`
	MobileIdSn        string `json:"mobile_id_sn"`
	Cid               string `json:"cid"`
	Created           string `json:"created"`
	Verifier          string `json:"verifier"`
	Used              string `json:"used"`
	Aal               string `json:"aal"`
	TxType            string `json:"tx_type"`
	Ref1              string `json:"ref1"`
	Ref2              string `json:"ref2"`
	Verified          string `json:"verified"`
	FaceScoreVerified string `json:"face_score_verified"`
	Revoked           string `json:"revoked"`
	RpId              string `json:"rp_id"`
}

type Member struct {
	MemberCode     string `json:"member_code"`
	MemberName     string `json:"member_name"`
	MemberRole     string `json:"member_role"`
	Status         string `json:"status"`
	Registered     string `json:"registered"`
	SerialNoPrefix string `json:"serial_no_prefix"`
	ServiceUrl     string `json:"service_url"`
}

type OperationResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MobileIdKey struct {
	MobileNo string `json:"mobile_no"`
}

type ConsentLogAuditQuery struct {
	Verifier     string `json:"verifier"`
	Records      string `json:"records"`
	RetrieveType string `json:"retrieve_type"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

type ConsentLogQuery struct {
	MobileNo   string `json:"mobile_no"`
	MobileIdSn string `json:"mobile_id_sn"`
	Issuer     string `json:"issuer"`
	Verifier   string `json:"verifier"`
	Records    string `json:"records"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

type HealthCheck struct {
	MemberCode string `json:"member_code"`
	Timestamp  string `json:"timestamp"`
}

type MobileIdSnRequest struct {
	MobileNo   string `json:"mobile_no"`
	MobileIdSn string `json:"mobile_id_sn"`
}

type MobileIdStatus struct {
	Status string `json:"status"`
}

type CountMobileIdStatusResponse struct {
	Status string `json:"status"`
	Count  string `json:"count"`
}

type ConsentLogQueryByDate struct {
	CreatedStartDate string `json:"created_start_date"`
	CreatedEndDate   string `json:"created_end_date"`
	UsedStartDate    string `json:"used_start_date"`
	UsedEndDate      string `json:"used_end_date"`
	RevokedStartDate string `json:"revoked_start_date"`
	RevokedEndDate   string `json:"revoked_end_date"`
	Records          string `json:"records"`
}

type MobileIdIssuerRequest struct {
	RpId     string `json:"rp_id"`
	TxId     string `json:"tx_id"`
	MobileNo string `json:"mobile_no"`
	Aal      string `json:"aal"`
}

type MobileIdIssuer struct {
	MobileNo   string `json:"mobile_no"`
	MobileIdSn string `json:"mobile_id_sn"`
	Issuer     string `json:"issuer"`
	Status     string `json:"status"`
	ServiceUrl string `json:"service_url"`
}

type MobileIdRequestLog struct {
	MobileNo   string `json:"mobile_no"`
	MobileIdSn string `json:"mobile_id_sn"`
	Issuer     string `json:"issuer"`
	Verifier   string `json:"verifier"`
	RpId       string `json:"rp_id"`
	TxId       string `json:"tx_id"`
	Status     string `json:"status"`
	Aal        string `json:"aal"`
	Created    string `json:"created"`
}

type MobileIdRequestLogKey struct {
	MobileNo string `json:"mobile_no"`
	RpId     string `json:"rp_id"`
	TxId     string `json:"tx_id"`
}

type MobileIdRequestLogQueryByDate struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Records   string `json:"records"`
}

//*******************************************************
// Define variable
//*******************************************************

const (
	TimestampRange = 24 * 60 * 60
	TimeFormat     = "2006-01-02T15:04:05"
)

const (
	MobileIdKeyObjectType        = "mobileid_key"
	MobileIdConsentKeyObjectType = "mobileid_consent_key"
	MemberKeyObjectType          = "member_key"
	MobileIdRequestKeyObjectType = "mobileid_request_key"
)

const (
	MobileIdActive = "A"
	// MobileIdSuspended  = "S"
	MobileIdTerminated = "T"
)

const (
	MemberActive     = "A"
	MemberTerminated = "T"
)

const (
	MemberIssuer    = "ISSUER"
	MemberVerifier  = "VERIFIER"
	MemberRegulator = "REGULATOR"
)

const (
	MessageProcessing         = "Operation is processing..."
	MessageSuccess            = "Operation is successful"
	MessageMemberAccessDenied = "Access denied: Membership role is not support this function"
	MessageUnauthorized       = "Unauthorized"
)

const (
	StatusCodeSuccess          = 200
	StatusCodeInternalError    = 500
	StatusCodeInvalidArgument  = 400
	StatusCodeInvalidOperation = 412
	StatusCodeUnauthorized     = 401
	StatusCodeNotFound         = 404
)

const (
	MessageCodeSuccess          = "Success"
	MessageCodeInternalError    = "Internal Error"
	MessageCodeInvalidArgument  = "Invalid Argument"
	MessageCodeInvalidOperation = "Invalid Operation"
	MessageCodeUnauthorized     = "Unauthorized"
	MessageCodeNotFound         = "Not found"
)

var StatusCodeMap = map[int]string{
	StatusCodeSuccess:          MessageCodeSuccess,
	StatusCodeInternalError:    MessageCodeInternalError,
	StatusCodeInvalidArgument:  MessageCodeInvalidArgument,
	StatusCodeInvalidOperation: MessageCodeInvalidOperation,
	StatusCodeUnauthorized:     MessageCodeUnauthorized,
	StatusCodeNotFound:         MessageCodeNotFound,
}

const (
	RetrieveTypeNormal = "NORMAL"
	RetrieveTypeRandom = "RANDOM"
)

// ------------- Logging Control and Chaincode Loggers ---------------

// As independent programs, Go language chaincodes can use any logging
// methodology they choose, from simple fmt.Printf() to os.Stdout, to
// decorated logs created by the author's favorite logging package. The
// chaincode "shim" interface, however, is defined by the Hyperledger fabric
// and implements its own logging methodology. This methodology currently
// includes severity-based logging control and a standard way of decorating
// the logs.
//
// The facilities defined here allow a Go language chaincode to control the
// logging level of its shim, and to create its own logs formatted
// consistently with, and temporally interleaved with the shim logs without
// any knowledge of the underlying implementation of the shim, and without any
// other package requirements. The lack of package requirements is especially
// important because even if the chaincode happened to explicitly use the same
// logging package as the shim, unless the chaincode is physically included as
// part of the hyperledger fabric source code tree it could actually end up
// using a distinct binary instance of the logging package, with different
// formats and severity levels than the binary package used by the shim.
//
// Another approach that might have been taken, and could potentially be taken
// in the future, would be for the chaincode to supply a logging object for
// the shim to use, rather than the other way around as implemented
// here. There would be some complexities associated with that approach, so
// for the moment we have chosen the simpler implementation below. The shim
// provides one or more abstract logging objects for the chaincode to use via
// the NewLogger() API, and allows the chaincode to control the severity level
// of shim logs using the SetLoggingLevel() API.

// LoggingLevel is an enumerated type of severity levels that control
// chaincode logging.
type LoggingLevel logging.Level

// These constants comprise the LoggingLevel enumeration
const (
	LogDebug    = LoggingLevel(logging.DEBUG)
	LogInfo     = LoggingLevel(logging.INFO)
	LogNotice   = LoggingLevel(logging.NOTICE)
	LogWarning  = LoggingLevel(logging.WARNING)
	LogError    = LoggingLevel(logging.ERROR)
	LogCritical = LoggingLevel(logging.CRITICAL)
)

var shimLoggingLevel = LogInfo // Necessary for correct initialization; See Start()

// SetLoggingLevel allows a Go language chaincode to set the logging level of
// its shim.
func SetLoggingLevel(level LoggingLevel) {
	shimLoggingLevel = level
	logging.SetLevel(logging.Level(level), "shim")
}

// LogLevel converts a case-insensitive string chosen from CRITICAL, ERROR,
// WARNING, NOTICE, INFO or DEBUG into an element of the LoggingLevel
// type. In the event of errors the level returned is LogError.
func LogLevel(levelString string) (LoggingLevel, error) {
	l, err := logging.LogLevel(levelString)
	level := LoggingLevel(l)
	if err != nil {
		level = LogError
	}
	return level, err
}

// ------------- Chaincode Loggers ---------------

// ChaincodeLogger is an abstraction of a logging object for use by
// chaincodes. These objects are created by the NewLogger API.
type ChaincodeLogger struct {
	logger *logging.Logger
}

// NewLogger allows a Go language chaincode to create one or more logging
// objects whose logs will be formatted consistently with, and temporally
// interleaved with the logs created by the shim interface. The logs created
// by this object can be distinguished from shim logs by the name provided,
// which will appear in the logs.
func NewLogger(name string) *ChaincodeLogger {
	return &ChaincodeLogger{logging.MustGetLogger(name)}
}

// SetLevel sets the logging level for a chaincode logger. Note that currently
// the levels are actually controlled by the name given when the logger is
// created, so loggers should be given unique names other than "shim".
func (c *ChaincodeLogger) SetLevel(level LoggingLevel) {
	logging.SetLevel(logging.Level(level), c.logger.Module)
}

// IsEnabledFor returns true if the logger is enabled to creates logs at the
// given logging level.
func (c *ChaincodeLogger) IsEnabledFor(level LoggingLevel) bool {
	return c.logger.IsEnabledFor(logging.Level(level))
}

// Debug logs will only appear if the ChaincodeLogger LoggingLevel is set to
// LogDebug.
func (c *ChaincodeLogger) Debug(args ...interface{}) {
	c.logger.Debug(args...)
}

// Info logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogInfo or LogDebug.
func (c *ChaincodeLogger) Info(args ...interface{}) {
	c.logger.Info(args...)
}

// Notice logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogNotice, LogInfo or LogDebug.
func (c *ChaincodeLogger) Notice(args ...interface{}) {
	c.logger.Notice(args...)
}

// Warning logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogWarning, LogNotice, LogInfo or LogDebug.
func (c *ChaincodeLogger) Warning(args ...interface{}) {
	c.logger.Warning(args...)
}

// Error logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogError, LogWarning, LogNotice, LogInfo or LogDebug.
func (c *ChaincodeLogger) Error(args ...interface{}) {
	c.logger.Error(args...)
}

// Critical logs always appear; They can not be disabled.
func (c *ChaincodeLogger) Critical(args ...interface{}) {
	c.logger.Critical(args...)
}

// Debugf logs will only appear if the ChaincodeLogger LoggingLevel is set to
// LogDebug.
func (c *ChaincodeLogger) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}

// Infof logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogInfo or LogDebug.
func (c *ChaincodeLogger) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}

// Noticef logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogNotice, LogInfo or LogDebug.
func (c *ChaincodeLogger) Noticef(format string, args ...interface{}) {
	c.logger.Noticef(format, args...)
}

// Warningf logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogWarning, LogNotice, LogInfo or LogDebug.
func (c *ChaincodeLogger) Warningf(format string, args ...interface{}) {
	c.logger.Warningf(format, args...)
}

// Errorf logs will appear if the ChaincodeLogger LoggingLevel is set to
// LogError, LogWarning, LogNotice, LogInfo or LogDebug.
func (c *ChaincodeLogger) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}

// Criticalf logs always appear; They can not be disabled.
func (c *ChaincodeLogger) Criticalf(format string, args ...interface{}) {
	c.logger.Criticalf(format, args...)
}
