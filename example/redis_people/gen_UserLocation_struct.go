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

const (
	UserLocationMysqlFieldKey       = "key"
	UserLocationMysqlFieldLongitude = "longitude"
	UserLocationMysqlFieldLatitude  = "latitude"
	UserLocationMysqlFieldValue     = "value"
)

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

// Get_UserLocationMgr returns the orm manager in case of its name starts with lower letter
func Get_UserLocationMgr() *_UserLocationMgr { return UserLocationMgr }

func (m *_UserLocationMgr) NewUserLocation() *UserLocation {
	rval := new(UserLocation)
	return rval
}
