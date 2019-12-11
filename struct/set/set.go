package set

type ISet interface {
	Add(interface{})
	Remove(interface{})
	Contains(interface{}) bool
	Size() int
	IsEmpty() bool
}
