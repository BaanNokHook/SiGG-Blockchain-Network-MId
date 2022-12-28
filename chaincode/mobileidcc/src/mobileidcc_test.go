package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
	fmt.Println(string(res.Payload))
}

func checkInvokeFail(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	fmt.Println("-------------------------------------------")
	fmt.Println("------------------Failed-------------------")
	fmt.Println("-------------------------------------------")
	if res.Status == shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvokeWithValue(t *testing.T, stub *shim.MockStub, args [][]byte, value string) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println(string(res.Payload))
		fmt.Println("Invoke value was not", value, "as expected")
		t.FailNow()
	} else {
		fmt.Println(string(res.Payload))
	}
}

func getTimeString() string {
	timeLocation, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(timeLocation).Format("2006-01-02T15:04:05")
}

/*
*****************************************************************************************************************************************
**********************************************************NORMAL*************************************************************************
*****************************************************************************************************************************************
 */

func TestMobileID_Init(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_Init\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_EnrollMobileID(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_EnrollMobileID\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_UpdateMobileID(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_UpdateMobileID\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Update MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "NEC01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_UpdateMobileID_Suspend(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_UpdateMobileID_Suspend\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Suspend MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "S", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_UpdateMobileID_Resume(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_UpdateMobileID_Resume\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Suspend MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "S", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Resume MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_UpdateMobileID_Terminate(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_UpdateMobileID_Terminate\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Terminate MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "T", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_RecordConsent(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_RecordConsent\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_RetrieveMobileId(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_RetrieveMobileId\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	mbIdData := `{"dg1":{"mobile_no":"0891234567","issuer":"AIS","mobile_id_sn":"030891234567","ial":"2.1","status":"A","face_engine_id":"YIT01234","timestamp":"` + getTimeString() + `"},"dg2":"A","sod":{"hdg1":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","hdg2":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","hface_template":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","dig":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","sig":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","cert":"1"}}`
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(mbIdData)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeWithValue(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)}, mbIdData)

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_RecordVerificationResult(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_RecordVerificationResult\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_GetConsentLog(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_GetConsentLog\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***Get Consent Log By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Get Consent Log By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

/*
*****************************************************************************************************************************************
**********************************************************Membership*********************************************************************
*****************************************************************************************************************************************
 */

func TestMobileID_ListMembers(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_ListMembers\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***List Members from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("listMembers"), []byte(`{}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_GetMember(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_GetMember\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Get Member (AIS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getMember"), []byte(`{"member_code": "AIS"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_CreateMember(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_CreateMember\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Create Member (ABS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_UpdateMember(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_UpdateMember\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Create Member (ABS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})

	fmt.Printf("\n***Update Member (ABS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("updateMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "T","serial_no_prefix": "99"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

/*
*****************************************************************************************************************************************
**********************************************************HealthCheck********************************************************************
*****************************************************************************************************************************************
 */

func TestMobileID_InvokeHealthCheck(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_InvokeHealthCheck\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Invoke Health Check from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("invokeHealthCheck"), []byte(``)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_ListHealthCheck(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_ListHealthCheck\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Invoke Health Check from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("invokeHealthCheck"), []byte(``)})

	fmt.Printf("\n***Invoke Health Check from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("invokeHealthCheck"), []byte(``)})

	fmt.Printf("\n***Invoke Health Check from BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("invokeHealthCheck"), []byte(``)})

	fmt.Printf("\n***List Health Check from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("listHealthCheck"), []byte(``)})

	fmt.Printf("\n***List Health Check from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("listHealthCheck"), []byte(``)})

	fmt.Printf("\n***List Health Check from BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("listHealthCheck"), []byte(``)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

/*
*****************************************************************************************************************************************
**********************************************************AUDIT**************************************************************************
*****************************************************************************************************************************************
 */

func TestMobileID_getMobileIdForAudit(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_getMobileIdForAudit\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***getMobileIdForAudit from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getMobileIdForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_getConsentLogForAudit(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_getConsentLogForAudit\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***getConsentLogForAudit By NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Failed Case
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func TestMobileID_getMobileIdForAudit_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_getMobileIdForAudit\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***getMobileIdForAudit from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("getMobileIdForAudit"), []byte(`{"mobile_no": "089"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("getMobileIdForAudit"), []byte(`{"mobile_no": "0891234568"`)})
	checkInvoke(t, stub, [][]byte{[]byte("getMobileIdForAudit"), []byte(`{"mobile_no": "0891234567"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_getConsentLogForAudit_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_getConsentLogForAudit\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***getConsentLogForAudit By NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac8"}`)})
	checkInvoke(t, stub, [][]byte{[]byte("getConsentLogForAudit"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_listConsentLogForAudit_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_listConsentLogForAudit\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***listConsentLogForAudit By NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogForAudit"), []byte(`{"verifier": "BBL", "records": "100", "retrieve_type": "NORMAL", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogForAudit"), []byte(`{"verifier": "BB", "records": "100", "retrieve_type": "NORMAL", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogForAudit"), []byte(`{"verifier": "BBL", "records": "0", "retrieve_type": "NORMAL", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogForAudit"), []byte(`{"verifier": "BBL", "records": "100", "retrieve_type": "SHUFFLE", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogForAudit"), []byte(`{"verifier": "BBL", "records": "100", "retrieve_type": "NORMAL", "start_date": "2019-10-01", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogForAudit"), []byte(`{"verifier": "BBL", "records": "100", "retrieve_type": "NORMAL", "start_date": "2019-10-01T15:04:05", "end_date": "2019-10-01"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Failed Case!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_listConsentLogByIssuer_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_listConsentLogByIssuer\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***listConsentLogByIssuer By NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "records": "0", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByIssuer"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "2019-10-01"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Failed Case!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_listConsentLogByVerifier_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_listConsentLogByVerifier\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***listConsentLogByVerifier By NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "0891234567", "verifier": "BBL", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "089", "verifier": "BBL", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "0891234567", "verifier": "BB", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "0891234567", "verifier": "BBL", "mobile_id_sn": "03", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "0891234567", "verifier": "BBL", "mobile_id_sn": "030891234567", "records": "0", "start_date": "2019-10-01T15:04:05", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "0891234567", "verifier": "BBL", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01", "end_date": "` + getTimeString() + `"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("listConsentLogByVerifier"), []byte(`{"mobile_no": "0891234567", "verifier": "BBL", "mobile_id_sn": "030891234567", "records": "100", "start_date": "2019-10-01T15:04:05", "end_date": "2019-10-01"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Failed Case!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_getConsentHistory_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_getConsentHistory\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***getConsentHistory By NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentHistory"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Failed Case!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_EnrollMobileID_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_EnrollMobileID\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "BBL", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Enroll MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "BBL", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.10", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "0.3", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "S", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "R", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT1", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "2019-10-01T15:04:05"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "2019-10-01"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "A", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "A", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "A", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "A", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "A", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": ""}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "020891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "S", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "T", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "031891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_updateMobileId_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_updateMobileId\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "T", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Enroll MobileID from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "T", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "BBL", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "031891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.10", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "0.3", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT1", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "2019-10-01T15:04:05"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "2019-10-01"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "A", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "A", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "A", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "A", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "A", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": ""}}`)})

	checkInvoke(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "T", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_RecordConsent_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileIDCC_RecordConsent\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "031891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "BBL", "mobile_id_sn": "031891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_RetrieveMobileId_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_RetrieveMobileId\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	mbIdData := `{"dg1":{"mobile_no":"0891234567","issuer":"AIS","mobile_id_sn":"030891234567","ial":"2.1","status":"A","face_engine_id":"YIT01234","timestamp":"` + getTimeString() + `"},"dg2":"A","sod":{"hdg1":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","hdg2":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","hface_template":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","dig":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","sig":"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA","cert":"1"}}`
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(mbIdData)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "031891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "AIS", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BB", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "0.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.30", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "012345678901234567890123456789"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac0", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_RecordVerificationResult_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_RecordVerificationResult\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "AIS", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "089", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "03", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BB", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "A", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "012345678910"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "-1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac0", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "031891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "AIS", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_GetConsentLog_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_GetConsentLog\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Enroll MobileID***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("enrollMobileId"), []byte(`{"dg1": {"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "ial": "2.1", "status": "A", "face_engine_id": "YIT01234", "timestamp": "` + getTimeString() + `"}, "dg2": "A", "sod": {"hdg1": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hdg2": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "hface_template": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "dig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "sig": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", "cert": "1"}}`)})

	fmt.Printf("\n***Record Consent By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Retrieve MobileId By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("retrieveMobileIdAndUpdateConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "aal": "2.3", "tx_type": "A", "ref1": "ref1", "ref2": "ref2"}`)})

	fmt.Printf("\n***Record Verification Result By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("recordVerificationResult"), []byte(`{"mobile_no": "0891234567", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9", "verifier": "BBL", "verified": "Y", "face_score_verified": "1"}`)})

	fmt.Printf("\n***Get Consent Log By AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	fmt.Printf("\n***Get Consent Log By BBL***\n")
	sid = &msp.SerializedIdentity{Mspid: "BBLMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvoke(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "089", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AI", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "03", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac9"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19"}`)})

	checkInvokeFail(t, stub, [][]byte{[]byte("getConsentLog"), []byte(`{"mobile_no": "0891234567", "issuer": "AIS", "mobile_id_sn": "030891234567", "cid": "4054ec19-e0b0-41e9-8cc5-2deff4937ac0"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_ListMembers_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_ListMembers\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	// fmt.Printf("\n***Init***\n")
	// checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***List Members from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("listMembers"), []byte(`{}`)})

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***List Members from AIS***\n")
	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("listMembers"), []byte(`{}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_GetMember_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_GetMember\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Get Member (AIS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("getMember"), []byte(`{"member_code": "ABS"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("getMember"), []byte(`{"member_code": "AI"}`)})

	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("getMember"), []byte(`{"member_code": "AIS"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_CreateMember_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_CreateMember\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})

	fmt.Printf("\n***Create Member (ABS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "AB","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ADMIN","status": "A","serial_no_prefix": "99"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "B","serial_no_prefix": "99"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "000"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "T","serial_no_prefix": "99"}`)})
	checkInvoke(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_UpdateMember_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_UpdateMember\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***Create Member (ABS) from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)

	checkInvokeFail(t, stub, [][]byte{[]byte("updateMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "T","serial_no_prefix": "99"}`)})
	checkInvoke(t, stub, [][]byte{[]byte("createMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "A","serial_no_prefix": "99"}`)})
	checkInvokeFail(t, stub, [][]byte{[]byte("updateMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "T","serial_no_prefix": "999"}`)})

	sid = &msp.SerializedIdentity{Mspid: "AISMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("updateMember"), []byte(`{"member_code": "ABS","member_name": "ABS","member_role": "ISSUER,VERIFIER","status": "T","serial_no_prefix": "99"}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}

func TestMobileID_General_FailedCase(t *testing.T) {
	fmt.Printf("\n\n##############################################\n")
	fmt.Printf("TestMobileID_General_FailedCase\n")
	fmt.Printf("##############################################\n")
	scc := new(MobileIdChaincode)
	stub := shim.NewMockStub("mobileidcc", scc)
	sid := &msp.SerializedIdentity{}
	var b []byte
	var err error

	fmt.Printf("\n***Init***\n")
	checkInit(t, stub, [][]byte{[]byte("init")})

	fmt.Printf("\n***List Members from NBTC***\n")
	sid = &msp.SerializedIdentity{Mspid: "NBTCMSP", IdBytes: []byte("test")}
	b, err = proto.Marshal(sid)
	if err != nil {
		t.FailNow()
	}
	stub.SetCreator(b)
	checkInvokeFail(t, stub, [][]byte{[]byte("listMember"), []byte(`{}`)})

	fmt.Printf("##############################################\n")
	fmt.Printf("Success!\n")
	fmt.Printf("##############################################\n")
}
