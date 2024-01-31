package message

import "github.com/free5gc/nas/nasType"

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
