package timer

import "time"

// ServiceTimer holds the start and end timestamps.
type ServiceTimer struct {
	UEinfor   string
	StartTime time.Time
	EndTime   time.Time
}

// NORA ServiceTimer holds the start and end timestamps.
type NORAServiceTimer struct {
	UEinfor   string
	StartTime time.Time
	EndTime   time.Time
}

// NewServiceTimer creates a new ServiceTimer instance.
func NewServiceTimer(UEid string, StartTime time.Time) *ServiceTimer {
	return &ServiceTimer{
		UEinfor:   UEid,
		StartTime: StartTime,
	}
}

var ServiceTimerMap map[string]*ServiceTimer

func init() {
	ServiceTimerMap = make(map[string]*ServiceTimer)
}

func StoreTimeStamp(ue *ServiceTimer) {
	ServiceTimerMap[ue.UEinfor] = ue
}

func GetStartTime(UEinfor string) time.Time {
	ue, _ := ServiceTimerMap[UEinfor]
	return ue.StartTime
}

func SetStartTime(UEinfor string, ST time.Time) bool {
	ue, _ := ServiceTimerMap[UEinfor]
	ue.StartTime = ST
	if GetStartTime(UEinfor) != ST {
		return false
	}
	return true
}

//func SetEndTime(UEinfor string, ET time.Time) bool {
//	ue, _ := ServiceTimerMap[UEinfor]
//	ue.StartTime = ET
//	if GetStartTime(UEinfor) != ET {
//		return false
//	}
//	return true
//}

// CalculateServiceTime calculates the service time using the start and end timestamps.
func CalculateServiceTime(ST time.Time, ET time.Time) time.Duration {
	return ET.Sub(ST)
}

///// NORA-AKA time stamp calculation
//
//// NewServiceTimer creates a new ServiceTimer instance.
//func NORANewServiceTimer(UEid string) *NORAServiceTimer {
//	ST := time.Now()
//	return &NORAServiceTimer{
//		UEinfor:   UEid,
//		StartTime: ST,
//	}
//}
//
//var NORAServiceTimerMap map[string]*NORAServiceTimer
//
//func NORAinit() {
//	NORAServiceTimerMap = make(map[string]*NORAServiceTimer)
//}
//
//func NORAStoreTimeStamp(ue *NORAServiceTimer) {
//	NORAServiceTimerMap[ue.UEinfor] = ue
//}
//
//func NORAGetStartTime(UEinfor string) time.Time {
//	ue, _ := NORAServiceTimerMap[UEinfor]
//	return ue.StartTime
//}
//
//func NORACheckMap(UEinfor string) bool {
//	ue, _ := NORAServiceTimerMap[UEinfor]
//	if ue == nil {
//		return false
//	}
//	return true
//}
//
//func NORACheckStartTime(UEinfor string) bool {
//	ue, _ := NORAServiceTimerMap[UEinfor]
//	if ue.StartTime.IsZero() {
//		return false
//	}
//	return true
//}
//
//func NORASetStartTime(UEinfor string, ST time.Time) bool {
//	ue, _ := NORAServiceTimerMap[UEinfor]
//	ue.StartTime = ST
//	if NORAGetStartTime(UEinfor) != ST {
//		return false
//	}
//	return true
//}
//
////func NORASetEndTime(UEinfor string, ET time.Time) bool {
////	ue, _ := NORAServiceTimerMap[UEinfor]
////	ue.StartTime = ET
////	if GetStartTime(UEinfor) != ET {
////		return false
////	}
////	return true
////}
//
//// CalculateServiceTime calculates the service time using the start and end timestamps.
//func NORACalculateServiceTime(ST time.Time, ET time.Time) time.Duration {
//	return ET.Sub(ST)
//}
