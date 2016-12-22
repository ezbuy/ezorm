package test

import "time"

var _ time.Time

type UserLocation struct {
	Key       string  `db:"key" json:"key"`
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Value     int32   `db:"value" json:"value"`
	isNew     bool
}

func (p *UserLocation) GetNameSpace() string {
	return "people"
}

func (p *UserLocation) GetClassName() string {
	return "UserLocation"
}
func (p *UserLocation) GetStoreType() string {
	return "geo"
}

func (p *UserLocation) GetPrimaryKey() string {
	return ""
}

func (p *UserLocation) GetIndexes() []string {
	idx := []string{}
	return idx
}

type _UserLocationMgr struct {
}

var UserLocationMgr *_UserLocationMgr

func (m *_UserLocationMgr) NewUserLocation() *UserLocation {
	rval := new(UserLocation)
	return rval
}
