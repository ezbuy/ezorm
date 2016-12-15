package test

import "time"

var _ time.Time

type UserLocation struct {
	RegionId  int32   `bson:"RegionId" json:"RegionId"`
	Longitude float64 `bson:"Longitude" json:"Longitude"`
	Latitude  float64 `bson:"Latitude" json:"Latitude"`
	UserId    int32   `bson:"UserId" json:"UserId"`
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
