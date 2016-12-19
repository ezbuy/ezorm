package test

import "time"

var _ time.Time

type UserLocation struct {
	RegionId  int32   `json:"region_id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
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
	return "RegionId"
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
