package lock

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

// Outer represents the the thing calling the Lock() method
func outer() string {
	return inner()
}

// Inner represents the Lock() method calling caller()
func inner() string {
	return caller()
}

func TestCaller(t *testing.T) {
	RegisterTestingT(t)
	Ω(outer()).Should(Equal("github.com/bbengfort/x/lock.outer"))
}

type Lockable struct {
	MutexD
}

func (l *Lockable) Alpha() {
	l.Lock()
	defer l.Unlock()
	time.Sleep(time.Second * 1)
}

func (l *Lockable) Bravo() {
	l.Lock()
	defer l.Unlock()
	time.Sleep(time.Millisecond * 10)
}

func ExampleMutexD() {
	l := new(Lockable)
	go l.Alpha()
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 2; i++ {
		go l.Bravo()
	}
	fmt.Println(l.MutexD.String())
	// Output:
	// 1 locks requested by github.com/bbengfort/x/lock.(*Lockable).Alpha
	// 2 locks requested by github.com/bbengfort/x/lock.(*Lockable).Bravo
}

type RWLockable struct {
	RWMutexD
}

func (l *RWLockable) Alpha() {
	l.Lock()
	defer l.Unlock()
	time.Sleep(time.Second * 1)
}

func (l *RWLockable) Bravo() {
	l.RLock()
	defer l.RUnlock()
	time.Sleep(time.Millisecond * 100)
}

func ExampleRWMutexD() {
	l := new(RWLockable)
	go l.Alpha()
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 2; i++ {
		go l.Bravo()
	}
	fmt.Println(l.RWMutexD.String())
	// Output:
	// 1 locks requested by github.com/bbengfort/x/lock.(*Lockable).Alpha
	// 2 read locks requested by github.com/bbengfort/x/lock.(*Lockable).Bravo
}
