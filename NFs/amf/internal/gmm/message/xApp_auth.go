package message

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/milenage"
	"github.com/free5gc/util/ueauth"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type XAppAuthenticationRequest struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.AuthenticationRequestMessageIdentity
	nasType.SpareHalfOctetAndNgksi
	nasType.ABBA
	*nasType.AuthenticationParameterRAND
	//*nasType.AuthenticationParameterAUTN
	xAppAuthenticationParameterRAND
}

type xAppAuthenticationParameterRAND struct {
	Iei   int
	Octet []uint8
}

const (
	SqnMAx    int64 = 0xFFFFFFFFFFFF
	ind       int64 = 32
	keyStrLen int   = 32
	opStrLen  int   = 32
	opcStrLen int   = 32
)

const (
	authenticationRejected string = "AUTHENTICATION_REJECTED"
	resyncAMF              string = "0000"
)

func (r xAppAuthenticationParameterRAND) SetxAppRANDValue(rANDValue []uint8) {
	copy(r.Octet[0:16], rANDValue[:])
}

//func (r *xAppAuthenticationParameterRAND) SetxAppIei(iei uint8) {
//	r.Iei = iei
//}

func NewXAppAuthenticationRequest(iei uint8) (authenticationRequest *XAppAuthenticationRequest) {
	authenticationRequest = &XAppAuthenticationRequest{}
	//xAppAuthenticationParameterRAND.SetxAppIei(iei)
	return authenticationRequest
}

func aucSQN(opc, k, auts, rand []byte) ([]byte, []byte) {
	AK, SQNms := make([]byte, 6), make([]byte, 6)
	macS := make([]byte, 8)
	ConcSQNms := auts[:6]
	AMF, err := hex.DecodeString(resyncAMF)
	if err != nil {
		return nil, nil
	}

	logger.GmmLog.Infof("aucSQN: ConcSQNms=[%x]", ConcSQNms)

	err = milenage.F2345(opc, k, rand, nil, nil, nil, nil, AK)
	if err != nil {
		logger.GmmLog.Infof("aucSQN milenage F2345 err:", err)
	}

	for i := 0; i < 6; i++ {
		SQNms[i] = AK[i] ^ ConcSQNms[i]
	}

	logger.GmmLog.Infof("aucSQN: opc=[%x], k=[%x], rand=[%x], AMF=[%x], SQNms=[%x]\n", opc, k, rand, AMF, SQNms)
	// The AMF used to calculate MAC-S assumes a dummy value of all zeros
	err = milenage.F1(opc, k, rand, SQNms, AMF, nil, macS)
	if err != nil {
		logger.GmmLog.Infof("aucSQN milenage F1 err:", err)
	}
	logger.GmmLog.Infof("aucSQN: macS=[%x]\n", macS)
	return SQNms, macS
}

func strictHex(s string, n int) string {
	l := len(s)
	if l < n {
		return fmt.Sprintf(strings.Repeat("0", n-l) + s)
	} else {
		return s[l-n : l]
	}
}

func XAppAKAGenerateAUTH() (response *models.AuthenticationVector, err error) {
	var authInfoRequest models.AuthenticationInfoRequest
	authInfoRequest.ServingNetworkName = "5G:mnc093.mcc208.3gppnetwork.org"
	authInfoRequest.ResynchronizationInfo = nil
	authInfoRequest.SupportedFeatures = ""

	logger.GmmLog.Infof("In GenerateAuthDataProcedure")

	rand.Seed(time.Now().UnixNano())
	//var err error

	/*
		K, RAND, CK, IK: 128 bits (16 bytes) (hex len = 32)
		SQN, AK: 48 bits (6 bytes) (hex len = 12) TS33.102 - 6.3.2
		AMF: 16 bits (2 bytes) (hex len = 4) TS33.102 - Annex H
	*/

	hasK, hasOP, hasOPC := false, false, false
	var kStr, opStr, opcStr string
	var k, op, opc []byte

	//------------------------ Terry Modify Start --------------------------//
	//	Goals: Generate a OUT-X nas packet.                                 //
	//  Method:                                                             //
	//     1. Connect to mongo DB                                           //
	//     2. Get the basic op & K* in the mongo DB                         //
	//     3. Generate the parameter with Authentication                    //
	//----------------------------------------------------------------------// 8baf473f2f8fd09487cccbd7097c6862

	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//
	//err = client.Connect(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer client.Disconnect(ctx)
	//err = client.Ping(ctx, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Connected to MongoDB!")

	// kStr should be got from mongoDB
	kStr = "8baf473f2f8fd09487cccbd7097c6862"
	if len(kStr) == keyStrLen {
		k, err = hex.DecodeString(kStr)
		if err != nil {
			logger.GmmLog.Infof("err:", err)
		} else {
			hasK = true
		}
	}

	// opcstr should be got from mongoDB
	opcStr = "8e27b6af0e692e750f32667a3b14605d"
	if len(opcStr) == opcStrLen {
		opc, err = hex.DecodeString(opcStr)
		if err != nil {
			logger.GmmLog.Infof("err:", err)
		} else {
			hasOPC = true
		}
	} else {
		logger.GmmLog.Infof("opcStr length is ", len(opcStr))
	}

	// opstr should be got from mongoDB
	opStr = "0"
	if len(opStr) == opStrLen {
		op, err = hex.DecodeString(opStr)
		if err != nil {
			logger.GmmLog.Infof("err:", err)
		} else {
			hasOP = true
		}

	} else {
		logger.GmmLog.Infof("opStr length is ", len(opStr))
	}

	if !hasOPC && !hasOP {
		return nil, err
	}
	fmt.Println("opcStr: ", opcStr)

	if !hasOPC {
		if hasK && hasOP {
			opc, err = milenage.GenerateOPC(k, op)
			if err != nil {
				logger.GmmLog.Infof("milenage GenerateOPC err:", err)
			}
		} else {
			logger.GmmLog.Infof("Unable to derive OPC")
			return nil, err
		}
	}

	sqnStr := "16f3b3f70ff2"
	sqn, err := hex.DecodeString(sqnStr)
	if err != nil {
		logger.GmmLog.Infof("err:", err)
		return nil, err
	}
	fmt.Println("sqnStr: ", sqnStr)
	logger.GmmLog.Infof("K=[%x], sqn=[%x], OP=[%x], OPC=[%x]", k, sqn, op, opc)

	RAND := make([]byte, 16)
	_, err = cryptoRand.Read(RAND)
	if err != nil {
		logger.GmmLog.Infof("err:", err)
		return nil, err
	}

	hexString := "8000"
	//AMF, err := hex.DecodeString(authSubs.AuthenticationManagementField)
	AMF, err := hex.DecodeString(hexString)
	if err != nil {
		logger.GmmLog.Infof("err:", err)
		return nil, err
	}
	fmt.Println("AMF: ", AMF)
	logger.GmmLog.Infof("RAND=[%x], AMF=[%x]", RAND, AMF)

	// re-synchronization
	if authInfoRequest.ResynchronizationInfo != nil {
		logger.GmmLog.Infof("Authentication re-synchronization")

		Auts, deCodeErr := hex.DecodeString(authInfoRequest.ResynchronizationInfo.Auts)
		if deCodeErr != nil {
			logger.GmmLog.Infof("err:", deCodeErr)
			return nil, err
		}

		randHex, deCodeErr := hex.DecodeString(authInfoRequest.ResynchronizationInfo.Rand)
		if deCodeErr != nil {
			logger.GmmLog.Infof("err:", deCodeErr)
			return nil, err
		}

		SQNms, macS := aucSQN(opc, k, Auts, randHex)
		if reflect.DeepEqual(macS, Auts[6:]) {
			_, err = cryptoRand.Read(RAND)
			if err != nil {
				logger.GmmLog.Infof("err:", deCodeErr)
				return nil, err
			}

			// increment sqn authSubs.SequenceNumber
			bigSQN := big.NewInt(0)
			sqnStr = hex.EncodeToString(SQNms)
			logger.GmmLog.Infof("SQNstr=[%s]", sqnStr)
			bigSQN.SetString(sqnStr, 16)

			bigInc := big.NewInt(ind + 1)

			bigP := big.NewInt(SqnMAx)
			bigSQN = bigInc.Add(bigSQN, bigInc)
			bigSQN = bigSQN.Mod(bigSQN, bigP)
			sqnStr = fmt.Sprintf("%x", bigSQN)
			sqnStr = strictHex(sqnStr, 12)
		} else {
			//logger.GmmLog.Infof("Re-Sync MAC failed ", supi)
			logger.GmmLog.Infof("MACS ", macS)
			logger.GmmLog.Infof("Auts[6:] ", Auts[6:])
			logger.GmmLog.Infof("Sqn ", SQNms)
			return nil, err
		}
	}

	// increment sqn
	bigSQN := big.NewInt(0)
	sqn, err = hex.DecodeString(sqnStr)
	if err != nil {
		logger.GmmLog.Infof("err:", err)
		return nil, err
	}

	bigSQN.SetString(sqnStr, 16)

	bigInc := big.NewInt(1)
	bigSQN = bigInc.Add(bigSQN, bigInc)

	SQNheStr := fmt.Sprintf("%x", bigSQN)
	SQNheStr = strictHex(SQNheStr, 12)
	//patchItemArray := []models.PatchItem{
	//	{
	//		Op:    models.PatchOperation_REPLACE,
	//		Path:  "/sequenceNumber",
	//		Value: SQNheStr,
	//	},
	//}

	//var rsp *http.Response
	//rsp, err = client.AuthenticationDataDocumentApi.ModifyAuthentication(
	//	context.Background(), supi, patchItemArray)
	//if err != nil {
	//	problemDetails = &models.ProblemDetails{
	//		Status: http.StatusForbidden,
	//		Cause:  "modification is rejected ",
	//		Detail: err.Error(),
	//	}
	//
	//	logger.GmmLog.Infof("update sqn error:", err)
	//	return nil, problemDetails
	//}
	//defer func() {
	//	if rspCloseErr := rsp.Body.Close(); rspCloseErr != nil {
	//		logger.SdmLog.Errorf("ModifyAuthentication response body cannot close: %+v", rspCloseErr)
	//	}
	//}()

	// Run milenage
	macA, macS := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)

	// Generate macA, macS
	err = milenage.F1(opc, k, RAND, sqn, AMF, macA, macS)
	if err != nil {
		logger.GmmLog.Infof("milenage F1 err:", err)
	}

	// Generate RES, CK, IK, AK, AKstar
	// RES == XRES (expected RES) for server
	err = milenage.F2345(opc, k, RAND, RES, CK, IK, AK, AKstar)
	if err != nil {
		logger.GmmLog.Infof("milenage F2345 err:", err)
	}
	logger.GmmLog.Infof("milenage RES=[%s]", hex.EncodeToString(RES))

	// Generate AUTN
	logger.GmmLog.Infof("SQN=[%x], AK=[%x]", sqn, AK)
	logger.GmmLog.Infof("AMF=[%x], macA=[%x]", AMF, macA)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(sqn); i++ {
		SQNxorAK[i] = sqn[i] ^ AK[i]
	}
	logger.GmmLog.Infof("SQN xor AK=[%x]", SQNxorAK)
	AUTN := append(append(SQNxorAK, AMF...), macA...)
	logger.GmmLog.Infof("AUTN=[%x]", AUTN)

	var av models.AuthenticationVector

	// derive XRES*
	key := append(CK, IK...)
	FC := ueauth.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(authInfoRequest.ServingNetworkName)
	P1 := RAND
	P2 := RES

	kdfValForXresStar, err := ueauth.GetKDFValue(
		key, FC, P0, ueauth.KDFLen(P0), P1, ueauth.KDFLen(P1), P2, ueauth.KDFLen(P2))
	if err != nil {
		logger.GmmLog.Infof("Get kdfValForXresStar err: %+v", err)
	}
	xresStar := kdfValForXresStar[len(kdfValForXresStar)/2:]
	logger.GmmLog.Infof("xresStar=[%x]", xresStar)

	// derive Kausf
	FC = ueauth.FC_FOR_KAUSF_DERIVATION
	P0 = []byte(authInfoRequest.ServingNetworkName)
	P1 = SQNxorAK
	kdfValForKausf, err := ueauth.GetKDFValue(key, FC, P0, ueauth.KDFLen(P0), P1, ueauth.KDFLen(P1))
	if err != nil {
		logger.GmmLog.Infof("Get kdfValForKausf err: %+v", err)
	}
	logger.GmmLog.Infof("Kausf=[%x]", kdfValForKausf)

	// Fill in rand, xresStar, autn, kausf
	av.Rand = hex.EncodeToString(RAND)
	av.XresStar = hex.EncodeToString(xresStar)
	av.Autn = hex.EncodeToString(AUTN)
	av.Kausf = hex.EncodeToString(kdfValForKausf)
	av.AvType = models.AvType__5_G_HE_AKA

	response = &av
	//response.Supi = supi
	fmt.Println("av: ", av)
	return response, err
}
