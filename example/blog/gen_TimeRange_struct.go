package test

import "time"

var _ time.Time

type TimeRange struct {
	From []int64 `bson:"From" json:"From"`
	To   []int64 `bson:"To" json:"To"`
}

const (
	TimeRangeMgoFieldFrom = "From"
	TimeRangeMgoFieldTo   = "To"
)

func (p *TimeRange) GetNameSpace() string {
	return "blog"
}

func (p *TimeRange) GetClassName() string {
	return "TimeRange"
}

type _TimeRangeMgr struct {
}

var TimeRangeMgr *_TimeRangeMgr

// Get_TimeRangeMgr returns the orm manager in case of its name starts with lower letter
func Get_TimeRangeMgr() *_TimeRangeMgr { return TimeRangeMgr }

func (m *_TimeRangeMgr) NewTimeRange() *TimeRange {
	rval := new(TimeRange)
	return rval
}
