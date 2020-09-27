package sd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"time"
)

var (
	ErrKeyAlreadyExists  = fmt.Errorf("ErrKeyAlreadyExists")
	ErrEtcdLeaseNotFound = fmt.Errorf("ErrEtcdLeaseNotFound")
)

type EtcdConfig struct {
	Uri     string `yaml:"uri"`
	Timeout int    `json:"timeout" yaml:"timeout" `
}

type Etcd struct {
	endpoints []string
	client    *clientv3.Client
	kv        clientv3.KV
	timeout   time.Duration
	lease     clientv3.Lease
}

func NewEtcd(config *EtcdConfig) *Etcd {
	return newEtcd(strings.Split(config.Uri, ","), time.Duration(config.Timeout)*time.Second)
}

// create a etcd
func newEtcd(endpoints []string, timeout time.Duration) *Etcd {
	var client *clientv3.Client
	var err error

	conf := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	}

	if client, err = clientv3.New(conf); err != nil {
		panic(err)
	}
	etcd := &Etcd{
		endpoints: endpoints,
		client:    client,
		kv:        clientv3.NewKV(client),
		timeout:   timeout,
		lease:     clientv3.NewLease(client),
	}
	return etcd
}

func (e *Etcd) Close() error {
	_ = e.lease.Close()
	_ = e.client.Close()
	return nil
}

func (e *Etcd) TryLockWithTTL(key string, ttl int64) error {
	lease := clientv3.NewLease(e.client)
	grant, err := lease.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	defer func() { _ = lease.Close() }()
	ctx, cancelFunc := context.WithTimeout(context.Background(), e.timeout)
	defer cancelFunc()

	txn := e.client.Txn(ctx)
	etcdRes, err := txn.If(
		clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(grant.ID))).
		Commit()

	if err != nil {
		return err
	}
	if !etcdRes.Succeeded {
		return ErrKeyAlreadyExists
	}
	return nil
}

func (e *Etcd) InsertKVNoExisted(ctx context.Context, key, val string, leaseID clientv3.LeaseID) error {
	ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
	defer cancelFunc()

	var etcdRes *clientv3.TxnResponse
	var err error
	if leaseID != 0 {
		etcdRes, err = e.client.Txn(ctx).
			If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
			Then(clientv3.OpPut(key, val, clientv3.WithLease(leaseID))).
			Commit()
	} else {
		etcdRes, err = e.client.Txn(ctx).
			If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
			Then(clientv3.OpPut(key, val)).
			Commit()
	}
	if err != nil {
		return err
	}
	if !etcdRes.Succeeded {
		return ErrKeyAlreadyExists
	}
	return nil
}

func (e *Etcd) InsertKV(ctx context.Context, key, val string, leaseID clientv3.LeaseID) error {
	//ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
	//defer cancelFunc()
	//
	//var etcdRes *clientv3.TxnResponse
	//var err error
	//if leaseID != 0 {
	//	etcdRes, err = e.client.Txn(ctx).
	//		If(clientv3.Compare(clientv3.CreateRevision(key), ">", 0)).
	//		Then(clientv3.OpPut(key, val, clientv3.WithLease(leaseID))).
	//		//Else(clientv3.OpPut(key, val, clientv3.WithLease(leaseID))).
	//		Commit()
	//} else {
	//	etcdRes, err = e.client.Txn(ctx).
	//		If(clientv3.Compare(clientv3.CreateRevision(key), ">", 0)).
	//		Then(clientv3.OpPut(key, val)).
	//		//Else(clientv3.OpPut(key, val)).
	//		Commit()
	//}
	//if err != nil {
	//	return err
	//}
	//
	//if !etcdRes.Succeeded {
	//	return fmt.Errorf("InsertKV %s->%s", key, val)
	//}
	//return nil

	var err error
	if leaseID != 0 {
		_, err = e.kv.Put(ctx, key, val, clientv3.WithLease(leaseID))
	} else {
		_, err = e.kv.Put(ctx, key, val)
	}
	return err
}

func (e *Etcd) GetWithPrefixKey(ctx context.Context, prefix string) ([]string, []string, error) {
	var err error
	var resp *clientv3.GetResponse
	ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
	defer cancelFunc()

	if resp, err = e.kv.Get(ctx, prefix, clientv3.WithPrefix()); err != nil {
		return nil, nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil, nil
	}

	keys := make([]string, 0, len(resp.Kvs))
	values := make([]string, 0, len(resp.Kvs))
	for i := 0; i < len(resp.Kvs); i++ {
		keys = append(keys, string(resp.Kvs[i].Key))
		values = append(values, string(resp.Kvs[i].Value))
	}
	return keys, values, nil
}

func (e *Etcd) DelKey(ctx context.Context, key string) error {
	ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
	defer cancelFunc()
	_, err := e.kv.Delete(ctx, key)
	return err
}

func (e *Etcd) WatchWithPrefix(key string, cb func(*clientv3.Event) error) error {
	watcher := clientv3.NewWatcher(e.client)
	watchCh := watcher.Watch(context.Background(), key, clientv3.WithPrefix())
	go func() {
		for ch := range watchCh {
			if ch.Canceled {
				break
			}
			for _, event := range ch.Events {
				err := cb(event)
				if err != nil {
					panic(err)
				}
			}
		}
	}()
	return nil
}

func (e *Etcd) GrantLease(ttl int64) (clientv3.LeaseID, error) {
	grant, err := e.lease.Grant(context.Background(), ttl)
	if err != nil {
		return 0, err
	}
	return grant.ID, err
}

func (e *Etcd) RevokeLease(id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	return e.lease.Revoke(context.Background(), id)
}

func (e *Etcd) RenewLease(ctx context.Context, id clientv3.LeaseID) error {
	if id != 0 {
		ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
		defer cancelFunc()
		_, err := e.lease.KeepAliveOnce(ctx, id)
		return err
	}
	return ErrEtcdLeaseNotFound
}

func (e *Etcd) Get(ctx context.Context, key string) (value []byte, err error) {
	var getResponse *clientv3.GetResponse
	ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
	defer cancelFunc()
	if getResponse, err = e.kv.Get(ctx, key); err != nil {
		return
	}
	if len(getResponse.Kvs) == 0 {
		return
	}
	value = getResponse.Kvs[0].Value
	return
}

func (e *Etcd) KeepaliveWithTTL(ctx context.Context, key, value string, ttl int64) error {
	lease := clientv3.NewLease(e.client)
	grant, err := lease.Grant(ctx, ttl)
	if err != nil {
		return err
	}

	// 但这个grant会过期, 需要对其进行长时间的续租
	keepAliveResponseCh, err := lease.KeepAlive(ctx, grant.ID)
	if err != nil {
		return err
	}

	go func() {
		for keep := range keepAliveResponseCh {
			if keep == nil {
				panic(fmt.Sprintf("tx keepalive has lose key:%s", key))
			}
			//fmt.Printf("recv keepAliveResponse id:%d ttl:%d\n", keep.ID, keep.TTL)
		}
	}()

	// 执行事务需要超时
	ctx, cancelFunc := context.WithTimeout(ctx, e.timeout)
	defer cancelFunc()

	txn := e.client.Txn(ctx)
	resp, err := txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, value, clientv3.WithLease(grant.ID))).
		Else(clientv3.OpGet(key)).
		Commit()
	if err != nil {
		_ = lease.Close()
		return err
	}

	if !resp.Succeeded {
		_ = lease.Close()
		return fmt.Errorf("%s -> %s renew lease error\n", key, value)
	}

	return nil
}

//func (etcd *Etcd) Watch(key string, cb func(*clientv3.Event)) error {
//	watcher := clientv3.NewWatcher(etcd.client)
//	watchCh := watcher.Watch(context.Background(), key)
//	go func() {
//		for ch := range watchCh {
//			if ch.Canceled {
//				break
//			}
//			for _, event := range ch.Events {
//				cb(event)
//			}
//		}
//		fmt.Printf("the watcher lose for key:%s", key)
//	}()
//	return nil
//}
//

//
//func (etcd *Etcd) TxKeepaliveWithTTLNoExist(key, value string, ttl int64) error {
//	lease := clientv3.NewLease(etcd.client)
//	grant, err := lease.Grant(context.Background(), ttl)
//	if err != nil {
//		return err
//	}
//
//	// 但这个grant会过期, 需要对其进行长时间的续租
//	keepAliveResponseCh, err := lease.KeepAlive(context.Background(), grant.ID)
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		for keep := range keepAliveResponseCh {
//			if keep == nil {
//				panic(fmt.Sprintf("tx keepalive has lose key:%s", key))
//			}
//			//fmt.Printf("recv keepAliveResponse id:%d ttl:%d\n", keep.ID, keep.TTL)
//		}
//	}()
//
//	// 执行事务需要超时
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	txn := etcd.client.Txn(ctx)
//
//	/*
//		下面这一句，是构建了一个compare的条件，比较的是key的createRevision，
//		如果revision是0，则存入一个key，如果revision不为0，则读取这个key。
//		revision是etcd一个全局的序列号，每一个对etcd存储进行改动都会分配一个这个序号，
//		在v2中叫index，createRevision是表示这个key创建时被分配的这个序号。当key不存在时，createRivision是0。
//	*/
//	resp, err := txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
//		Then(clientv3.OpPut(key, value, clientv3.WithLease(grant.ID))).
//		Commit()
//	if err != nil {
//		_ = lease.Close()
//		return err
//	}
//
//	if !resp.Succeeded {
//		_ = lease.Close()
//		x := resp.Responses[0].GetResponseRange().Kvs[0]
//		return fmt.Errorf("%s -> %s not successed, cause %s -> %s\n", key, value, string(x.Key), string(x.Value))
//	}
//
//	fmt.Errorf("%s -> %s\n", key, value)
//	return nil
//}
//
//
//
//
//
//
//
//func (etcd *Etcd) TxKeepaliveWithTTL(key, value string, ttl int64) error {
//	lease := clientv3.NewLease(etcd.client)
//	grant, err := lease.Grant(context.Background(), ttl)
//	if err != nil {
//		return err
//	}
//
//
//	//lease.Revoke()
//
//	// 但这个grant会过期, 需要对其进行长时间的续租
//	keepAliveResponseCh, err := lease.KeepAlive(context.Background(), grant.ID)
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		for keep := range keepAliveResponseCh {
//			if keep == nil {
//				panic(fmt.Sprintf("tx keepalive has lose key:%s", key))
//			}
//			//fmt.Printf("recv keepAliveResponse id:%d ttl:%d\n", keep.ID, keep.TTL)
//		}
//	}()
//
//	// 执行事务需要超时
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	txn := etcd.client.Txn(ctx)
//
//	resp, err := txn.If(clientv3.Compare(clientv3.Version(key), "=", 0)).
//		Then(clientv3.OpPut(key, value, clientv3.WithLease(grant.ID))).
//		Else(clientv3.OpGet(key)).
//		Commit()
//	if err != nil {
//		_ = lease.Close()
//		return err
//	}
//
//	if !resp.Succeeded {
//		_ = lease.Close()
//		x := resp.Responses[0].GetResponseRange().Kvs[0]
//		return fmt.Errorf("%s -> %s not successed, cause %s -> %s\n", key, value, string(x.Key), string(x.Value))
//	}
//
//	fmt.Errorf("%s -> %s\n", key, value)
//	return nil
//}
//
//
////
////// get value  from a key
////func (etcd *Etcd) Get(key string) (value []byte, err error) {
////	var getResponse *clientv3.GetResponse
////	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
////	defer cancelFunc()
////	if getResponse, err = etcd.kv.Get(ctx, key); err != nil {
////		return
////	}
////	if len(getResponse.Kvs) == 0 {
////		return
////	}
////	value = getResponse.Kvs[0].Value
////	return
////}
////
//
//
//// get values from  prefixKey
//func (etcd *Etcd) GetWithPrefixKey(prefixKey string) (keys [][]byte, values [][]byte, err error) {
//	var getResponse *clientv3.GetResponse
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	if getResponse, err = etcd.kv.Get(ctx, prefixKey, clientv3.WithPrefix()); err != nil {
//		return
//	}
//	if len(getResponse.Kvs) == 0 {
//		return
//	}
//	keys = make([][]byte, 0)
//	values = make([][]byte, 0)
//	for i := 0; i < len(getResponse.Kvs); i++ {
//		keys = append(keys, getResponse.Kvs[i].Key)
//		values = append(values, getResponse.Kvs[i].Value)
//	}
//	return
//}

//
//// get values from  prefixKey limit
//func (etcd *Etcd) GetWithPrefixKeyLimit(prefixKey string, limit int64) (keys [][]byte, values [][]byte, err error) {
//	var getResponse *clientv3.GetResponse
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	if getResponse, err = etcd.kv.Get(ctx, prefixKey, clientv3.WithPrefix(), clientv3.WithLimit(limit)); err != nil {
//		return
//	}
//
//	if len(getResponse.Kvs) == 0 {
//		return
//	}
//
//	keys = make([][]byte, 0)
//	values = make([][]byte, 0)
//	for i := 0; i < len(getResponse.Kvs); i++ {
//		keys = append(keys, getResponse.Kvs[i].Key)
//		values = append(values, getResponse.Kvs[i].Value)
//	}
//	return
//}
//
//// put a key
//func (etcd *Etcd) Put(key, value string) (err error) {
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	if _, err = etcd.kv.Put(ctx, key, value); err != nil {
//		return
//	}
//	return
//}
//
//// put a key not exist
//func (etcd *Etcd) PutNotExist(key, value string) (success bool, oldValue []byte, err error) {
//
//	var (
//		txnResponse *clientv3.TxnResponse
//	)
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	txn := etcd.client.Txn(ctx)
//
//	txnResponse, err = txn.If(clientv3.Compare(clientv3.Version(key), "=", 0)).
//		Then(clientv3.OpPut(key, value)).
//		Else(clientv3.OpGet(key)).
//		Commit()
//
//	if err != nil {
//		return
//	}
//
//	if txnResponse.Succeeded {
//		success = true
//	} else {
//		oldValue = make([]byte, 0)
//		oldValue = txnResponse.Responses[0].GetResponseRange().Kvs[0].Value
//	}
//
//	return
//}
//
//
//func (etcd *Etcd) Update(key, value, oldValue string) (success bool, err error) {
//	var (
//		txnResponse *clientv3.TxnResponse
//	)
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//	txn := etcd.client.Txn(ctx)
//	txnResponse, err = txn.If(clientv3.Compare(clientv3.Value(key), "=", oldValue)).
//		Then(clientv3.OpPut(key, value)).
//		Commit()
//	if err != nil {
//		return
//	}
//	if txnResponse.Succeeded {
//		success = true
//	}
//	return
//}

//
//func (etcd *Etcd) Delete(key string) (err error) {
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	_, err = etcd.kv.Delete(ctx, key)
//
//	return
//}
//
//// delete the keys  with prefix key
//func (etcd *Etcd) DeleteWithPrefixKey(prefixKey string) (err error) {
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//	_, err = etcd.kv.Delete(ctx, prefixKey, clientv3.WithPrefix())
//	return
//}
//
// 自动续租kv

//
///*
//	//e := KeyChangeEvent{}
//	//e.Key = event.Kv.Key
//	//switch event.Type {
//	//case mvccpb.PUT:
//	//	if event.IsCreate() {
//	//		e.Type = KeyCreateChangeEvent
//	//	} else {
//	//		e.Type = KeyUpdateChangeEvent
//	//	}
//	//	e.Value = event.Kv.Value
//	//case mvccpb.DELETE:
//	//	e.Type = KeyDeleteChangeEvent
//	//}
//*/

//
//func (etcd *Etcd) Watch(key string, cb KeyChangeEventCallback) error {
//	watcher := clientv3.NewWatcher(etcd.client)
//	watchCh := watcher.Watch(context.Background(), key)
//	go func() {
//		for ch := range watchCh {
//			if ch.Canceled {
//				break
//			}
//			for _, event := range ch.Events {
//				cb(event)
//			}
//		}
//		fmt.Printf("the watcher lose for key:%s", key)
//	}()
//	return nil
//}

//
//// watch with prefix key
//func (etcd *Etcd) WatchWithPrefixKey(prefixKey string) (keyChangeEventResponse *WatchKeyChangeResponse) {
//
//	watcher := clientv3.NewWatcher(etcd.client)
//
//	watchChans := watcher.Watch(context.Background(), prefixKey, clientv3.WithPrefix())
//
//	keyChangeEventResponse = &WatchKeyChangeResponse{
//		Event:   make(chan *KeyChangeEvent, 250),
//		Watcher: watcher,
//	}
//
//	go func() {
//
//		for ch := range watchChans {
//
//			if ch.Canceled {
//				goto End
//			}
//			for _, event := range ch.Events {
//				etcd.handleKeyChangeEvent(event, keyChangeEventResponse.Event)
//			}
//		}
//
//	End:
//		log.Println("the watcher lose for prefixKey:", prefixKey)
//	}()
//
//	return
//}

//
//// handle the key change event
//func (etcd *Etcd) handleKeyChangeEvent(event *clientv3.Event, events chan *KeyChangeEvent) {
//
//	changeEvent := &KeyChangeEvent{
//		Key: string(event.Kv.Key),
//	}
//	switch event.Type {
//
//	case mvccpb.PUT:
//		if event.IsCreate() {
//			changeEvent.Type = KeyCreateChangeEvent
//		} else {
//			changeEvent.Type = KeyUpdateChangeEvent
//		}
//		changeEvent.Value = event.Kv.Value
//	case mvccpb.DELETE:
//
//		changeEvent.Type = KeyDeleteChangeEvent
//	}
//	events <- changeEvent
//
//}

//
//func (etcd *Etcd) TxWithTTL(key, value string, ttl int64) (txResponse *TxResponse, err error) {
//
//	var (
//		txnResponse *clientv3.TxnResponse
//		leaseID     clientv3.LeaseID
//		v           []byte
//	)
//	lease := clientv3.NewLease(etcd.client)
//
//	grantResponse, err := lease.Grant(context.Background(), ttl)
//
//	leaseID = grantResponse.ID
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	txn := etcd.client.Txn(ctx)
//	txnResponse, err = txn.If(
//		clientv3.Compare(clientv3.Version(key), "=", 0)).
//		Then(clientv3.OpPut(key, value, clientv3.WithLease(leaseID))).Commit()
//
//	if err != nil {
//		_ = lease.Close()
//		return
//	}
//
//	txResponse = &TxResponse{
//		LeaseID: leaseID,
//		Lease:   lease,
//	}
//	if txnResponse.Succeeded {
//		txResponse.Success = true
//	} else {
//		// close the lease
//		_ = lease.Close()
//		v, err = etcd.Get(key)
//		if err != nil {
//			return
//		}
//		txResponse.Success = false
//		txResponse.Key = key
//		txResponse.Value = string(v)
//	}
//	return
//}

//
//func (etcd *Etcd) TxKeepaliveWithTTLNoExist(key, value string, ttl int64) (txResponse *TxResponse, err error) {
//
//	var (
//		txnResponse    *clientv3.TxnResponse
//		leaseID        clientv3.LeaseID
//		aliveResponses <-chan *clientv3.LeaseKeepAliveResponse
//		v              []byte
//	)
//	lease := clientv3.NewLease(etcd.client)
//
//	grantResponse, err := lease.Grant(context.Background(), ttl)
//
//	leaseID = grantResponse.ID
//
//	if aliveResponses, err = lease.KeepAlive(context.Background(), leaseID); err != nil {
//
//		return
//	}
//
//	go func() {
//
//		for ch := range aliveResponses {
//
//			if ch == nil {
//				goto End
//			}
//
//		}
//
//	End:
//		log.Printf("the tx keepalive has lose key:%s", key)
//	}()
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	txn := etcd.client.Txn(ctx)
//	txnResponse, err = txn.If(
//		clientv3.Compare(clientv3.Version(key), "=", 0)).
//		Then(clientv3.OpPut(key, value, clientv3.WithLease(leaseID))).
//		Else(
//			clientv3.OpGet(key),
//		).Commit()
//
//	if err != nil {
//		_ = lease.Close()
//		return
//	}
//
//	txResponse = &TxResponse{
//		LeaseID: leaseID,
//		Lease:   lease,
//	}
//	if txnResponse.Succeeded {
//		txResponse.Success = true
//	} else {
//		// close the lease
//		_ = lease.Close()
//		txResponse.Success = false
//		if v, err = etcd.Get(key); err != nil {
//			return
//		}
//		txResponse.Key = key
//		txResponse.Value = string(v)
//	}
//	return
//}

//
//// transfer from  to with value
//func (etcd *Etcd) transfer(from string, to string, value string) (success bool, err error) {
//
//	var (
//		txnResponse *clientv3.TxnResponse
//	)
//
//	ctx, cancelFunc := context.WithTimeout(context.Background(), etcd.timeout)
//	defer cancelFunc()
//
//	txn := etcd.client.Txn(ctx)
//
//	txnResponse, err = txn.If(
//		clientv3.Compare(clientv3.Value(from), "=", value)).
//		Then(
//			clientv3.OpDelete(from),
//			clientv3.OpPut(to, value),
//		).Commit()
//
//	if err != nil {
//		return
//	}
//
//	success = txnResponse.Succeeded
//
//	return
//
//}
