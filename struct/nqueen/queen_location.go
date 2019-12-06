package nqueen

import (
	"fmt"
	"sync"

	"github.com/Zeb-D/go-util/struct/list"
)

type Location struct {
	x, y int
}

func (l *Location) String() string {
	return fmt.Sprintf("[x:%d,y:%d]", l.x, l.x)
}

func NewLocation(x, y int) *Location {
	return &Location{
		x: x,
		y: y,
	}
}

// QueenList thread safe
type QueenList struct {
	QList *list.LinkedList
	RW    sync.RWMutex
}

func NewQueenList() *QueenList {
	return &QueenList{
		QList: &list.LinkedList{},
	}
}

func (ql *QueenList) GetList() *list.LinkedList {
	return ql.QList
}

func (ql *QueenList) AddFirst(loc *Location) {
	ql.RW.Lock()
	defer ql.RW.Unlock()
	ql.QList.AddFirst(loc)
}

func (ql *QueenList) RemoveFirst() {
	ql.RW.Lock()
	defer ql.RW.Unlock()
	ql.QList.RemoveFirst()
}
