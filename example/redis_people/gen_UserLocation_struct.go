package test

import "time"

var _ time.Time

type UserLocation struct {
	Key       string  `db:"Key"	json:"Key"`
	Longitude float64 `db:"Longitude"	json:"Longitude"`
	Latitude  float64 `db:"Latitude"	json:"Latitude"`
	Value     int32   `db:"Value"	json:"Value"`
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
