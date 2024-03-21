package uestatus

type UeStatus struct {
	/* Ue Identity*/
	UEid   string
	Status bool
}

func NewAmfUe(UEid string, Active bool) *UeStatus {
	return &UeStatus{
		UEid:   UEid,
		Status: Active,
	}
}

// AmfUeMap maps UEid to corresponding AmfUe instances.
var UeStatusMap map[string]*UeStatus

func init() {
	UeStatusMap = make(map[string]*UeStatus)
}

func StoreAmfUe(ue *UeStatus) {
	UeStatusMap[ue.UEid] = ue
}

// Get UE status value corresponding to the given UEid.
func GetUEStatus(UEid string) bool {
	ue, _ := UeStatusMap[UEid]
	if ue.Status {
		return false
	}
	return true
}

func CheckUEStatus(UEid string) bool {
	ue, ok := UeStatusMap[UEid]
	if !ok {
		return false
	}
	if ue.Status == false {
		// ue is un-active
		return false
	}
	return true
}
