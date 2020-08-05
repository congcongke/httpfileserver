package lock

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Lock struct {
	mutex        sync.Mutex
	resources    map[string]*LockResource
	lockInterval time.Duration
}

type LockResource struct {
	lockTime time.Time
	subj     string
}

func NewLock() *Lock {
	return &Lock{
		resources:    map[string]*LockResource{},
		lockInterval: 5 * time.Minute,
	}
}

func (l *Lock) Lock(filename, subj string) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	res, ok := l.resources[filename]
	if !ok {
		l.resources[filename] = &LockResource{
			lockTime: time.Now(),
			subj:     subj,
		}
		return true
	}

	if res.lockTime.Add(l.lockInterval).Before(time.Now()) {
		res.lockTime = time.Now()
		res.subj = subj
		return true
	}

	return false
}

func (l *Lock) Unlock(filename, subj string, unlockFunc func() error) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	res, ok := l.resources[filename]
	if !ok {
		return false
	}

	if res.subj != subj {
		return false
	}

	delete(l.resources, filename)

	if err := unlockFunc(); err != nil {
		logrus.Error("failed for the unlock func: %v", err)
		return false
	}

	return true
}
