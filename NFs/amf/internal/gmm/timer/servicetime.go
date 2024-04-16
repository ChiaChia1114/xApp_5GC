package timer

import "time"

// ServiceTimer holds the start and end timestamps.
type ServiceTimer struct {
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

func SetEndTime(UEinfor string, ET time.Time) bool {
	ue, _ := ServiceTimerMap[UEinfor]
	ue.StartTime = ET
	if GetStartTime(UEinfor) != ET {
		return false
	}
	return true
}

// CalculateServiceTime calculates the service time using the start and end timestamps.
func CalculateServiceTime(ST time.Time, ET time.Time) time.Duration {
	return ET.Sub(ST)
}
