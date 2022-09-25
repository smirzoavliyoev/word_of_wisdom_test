package repo

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestGet(t *testing.T) {
	r := NewRepo()

	prevLen := len(r.Chalanges)
	r.Get()
	curLen := len(r.Chalanges)

	if prevLen-1 != curLen {
		t.FailNow()
	}
}

func TestSave(t *testing.T) {
	r := NewRepo()
	for i := 0; i < 100; i++ {
		r.Save(GenString(rand.Int() % 20))
	}

	fmt.Println(len(r.Chalanges))
	if len(r.Chalanges) != 100 {
		t.FailNow()
	}
}

func GenString(l int) string {
	alf := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	str := ""

	for i := 0; i < l; i++ {
		str += string(alf[rand.Int()%len(alf)])
	}
	return str
}
