package nimble

import (
	"code.google.com/p/nimble-cube/core"
)

type Chan1 struct {
	*Info
	slice Slice // TODO: rename buffer
	mutex *rwMutex
}

func makeChan1(tag, unit string, m *Mesh, memType MemType, bufBlocks int) Chan1 {

	N := -666
	if bufBlocks < 1 { // means auto
		N = idiv(m.NCell(), m.BlockLen()) // buffer all, just to be sure
	} else {
		N = m.BlockLen() * bufBlocks
	}

	return asChan1(MakeSlice(N, memType), tag, unit, m)

}

func asChan1(buffer Slice, tag, unit string, m *Mesh) Chan1 {
	AddQuant(tag)
	info := NewInfo(tag, unit, m)
	return Chan1{info, buffer, newRWMutex(buffer.Len(), tag)}
}

// WriteDone() signals a slice obtained by WriteNext() is fully
// written and can be sent down the Chan.
func (c Chan1) WriteDone() {
	c.mutex.WriteDone()
}

func (c Chan1) WriteNext(n int) Slice {
	c.mutex.WriteNext(n)
	a, b := c.mutex.WRange()
	return c.slice.Slice(a, b)
}

func (c Chan1) NComp() int           { return 1 }
func (c Chan1) BufLen() int          { return c.slice.Len() }
func (c Chan1) NBufferedBlocks() int { return idiv(c.NCell(), c.slice.Len()) }

//func (c *Chan1) WriteDelta(Δstart, Δstop int) []float32 {
//	c.mutex.WriteDelta(Δstart, Δstop)
//	a, b := c.mutex.WRange()
//	return c.slice.list[a:b]
//}

func idiv(a, b int) int {
	core.Assert(a%b == 0)
	return a / b
}
