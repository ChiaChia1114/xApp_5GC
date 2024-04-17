package timer

import "time"

// ServiceTimer holds the start and end timestamps.
type ServiceTimer struct {
	UEinfor   string
	StartTime time.Time
	EndTime   time.Time
}

// NORAServiceTimer holds the start and end timestamps for NORA-AKA.
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
	if ue, ok := ServiceTimerMap[UEinfor]; ok {
		return ue.StartTime
	}
	return time.Time{} // Return zero time if the key is not found
}

func SetStartTime(UEinfor string, ST time.Time) bool {
	if ue, ok := ServiceTimerMap[UEinfor]; ok {
		ue.StartTime = ST
		if GetStartTime(UEinfor) != ST {
			return false
		}
		return true
	}
	return false
}

// CalculateServiceTime calculates the service time using the start and end timestamps.
func CalculateServiceTime(ST time.Time, ET time.Time) time.Duration {
	return ET.Sub(ST)
}
