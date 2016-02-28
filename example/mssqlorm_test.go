package customer

import (
	"testing"
)

func TestFindOne(t *testing.T) {
	item, err := CustomerMgr.FindOne("NickName = ? AND CustomerId = ?", "jerry123", 384381)

	if err != nil {
		t.Error(err.Error())
	}

	if item.CustomerId != 384381 {
		t.Error(item)
	}
}

func TestFind(t *testing.T) {
	items, err := CustomerMgr.Find("CustomerId < ?", 10)
	if err != nil {
		t.Error(err.Error())
	}

	for _, v := range items {
		if v.CustomerId >= 10 {
			t.Error(v)
		}
	}
}

func TestFindAll(t *testing.T) {
	items, err := CustomerMgr.FindAll()

	if err != nil {
		t.Error(err.Error())
	}
	if len(items) != 371552 {
		t.Error(len(items))
	}
}

func TestFindWithOffset(t *testing.T) {
	items, err := CustomerMgr.FindWithOffset("CatalogCode = ?", 1, 10, "SG")
	if err != nil {
		t.Error(err.Error())
	}

	if len(items) != 10 {
		t.Error(len(items))
	}
}
