package test

import "time"

var _ time.Time

type User struct {
	Id          int32     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Mailbox     string    `db:"mailbox" json:"mailbox"`
	Sex         bool      `db:"sex" json:"sex"`
	Longitude   float64   `db:"longitude" json:"longitude"`
	Latitude    float64   `db:"latitude" json:"latitude"`
	Description string    `db:"description" json:"description"`
	Password    string    `db:"password" json:"password"`
	HeadUrl     string    `db:"head_url" json:"head_url"`
	Status      int32     `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	isNew       bool
}

const (
	UserMysqlFieldId          = "id"
	UserMysqlFieldName        = "name"
	UserMysqlFieldMailbox     = "mailbox"
	UserMysqlFieldSex         = "sex"
	UserMysqlFieldLongitude   = "longitude"
	UserMysqlFieldLatitude    = "latitude"
	UserMysqlFieldDescription = "description"
	UserMysqlFieldPassword    = "password"
	UserMysqlFieldHeadUrl     = "head_url"
	UserMysqlFieldStatus      = "status"
	UserMysqlFieldCreatedAt   = "created_at"
	UserMysqlFieldUpdatedAt   = "updated_at"
)

func (p *User) GetNameSpace() string {
	return "people"
}

func (p *User) GetClassName() string {
	return "User"
}
func (p *User) GetStoreType() string {
	return "hash"
}

func (p *User) GetPrimaryKey() string {
	return "Id"
}

func (p *User) GetIndexes() []string {
	idx := []string{
		"Name",
	}
	return idx
}

type _UserMgr struct {
}

var UserMgr *_UserMgr

// Get_UserMgr returns the orm manager in case of its name starts with lower letter
func Get_UserMgr() *_UserMgr { return UserMgr }

func (m *_UserMgr) NewUser() *User {
	rval := new(User)
	return rval
}
