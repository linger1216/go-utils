package jelly

type TimeRange struct {
	StartTime *Timestamp `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"startTime,omitempty"`
	EndTime   *Timestamp `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"endTime,omitempty"`
}

func (m *TimeRange) GetStartTime() *Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *TimeRange) GetEndTime() *Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}
