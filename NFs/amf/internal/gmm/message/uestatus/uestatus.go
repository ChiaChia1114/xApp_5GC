package uestatus

type UeStatus struct {
	/* Ue Identity*/
	UEid   string
	Status bool
	count  int
}

func NewAmfUe(UEid string, Active bool) *UeStatus {
	return &UeStatus{
		UEid:   UEid,
		Status: Active,
		count:  0,
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
func GetUEStatus(Supi string) bool {
	ue, _ := UeStatusMap[Supi]
	if ue.Status {
		return false
	}
	return true
}

func CheckUEStatus(Supi string) bool {
	ue, ok := UeStatusMap[Supi]
	if !ok {
		return false
	}
	if ue.Status == false {
		// ue is un-active
		return false
	}
	return true
}

func GetCount(Supi string) int {
	ue, ok := UeStatusMap[Supi]
	if ok {
		Count := ue.count
		return Count
	}
	return 0
}

func DeleteAmfUe(Supi string) {
	delete(UeStatusMap, Supi)
}

func CountPlus(Supi string) bool {
	ue, ok := UeStatusMap[Supi]
	if ok {
		Count := ue.count
		Count++
		ue.count = Count
		StoreAmfUe(ue)
		CheckResult := GetCount(Supi)
		if CheckResult != Count {
			return false
		}
	}
	return true
}
