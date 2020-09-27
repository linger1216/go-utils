package pattern

import (
	"fmt"
	"sync"
)

type CriticalSection struct {
	mux  sync.RWMutex
	data map[interface{}]chan struct{}
}

/*
RequestTag 请求标识 用于标识同一个资源
同一时刻只有一个请求能获取执行权限，获得执行权限的线程接下来需要执行具体的业务逻辑，
完成后调用release方法通知其他线程，操作已完成，获取资源即可
其他请求接下来需要调用wait方法
*/

func (u *CriticalSection) Require(key interface{}) bool {
	u.mux.Lock()
	if u.data == nil {
		u.data = make(map[interface{}]chan struct{})
	}

	_, ok := u.data[key]
	if ok {
		u.mux.Unlock()
		return false
	}

	u.data[key] = make(chan struct{})
	u.mux.Unlock()
	return true
}

/*RequestTag 请求标识 用于标识同一个资源
调用wait方法将处于阻塞状态，直到获得执行权限的线程处理完具体的业务逻辑，调用release方法来通知其他线程资源ok了
*/
func (u *CriticalSection) Wait(key interface{}) {
	u.mux.RLock()
	user, ok := u.data[key]
	u.mux.RUnlock()
	if !ok {
		return
	}
	select {
	case x := <-user:
		fmt.Println("recv ", x)
		return
	}
}

/*RequestTag 请求标识 用于标识同一个资源
获得执行权限的线程需要在执行完业务逻辑后调用该方法通知其他处于阻塞状态的线程
*/
func (u *CriticalSection) Release(key interface{}) {
	u.mux.Lock()
	defer u.mux.Unlock()
	if _, ok := u.data[key]; !ok {
		return
	}
	close(u.data[key])
	delete(u.data, key)
}
