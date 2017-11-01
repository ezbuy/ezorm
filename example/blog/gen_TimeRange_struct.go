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

func (m *_TimeRangeMgr) NewTimeRange() *TimeRange {
	rval := new(TimeRange)
	return rval
}
