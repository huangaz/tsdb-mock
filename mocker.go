package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	cpool "github.com/open-falcon/falcon/transfer/sender/conn_pool"
)

const (
	ConnTimeout = 1000
	CallTimeout = 5000
)

type MetricValue struct {
	Endpoint  string  `json:"endpoint"`
	Metric    string  `json:"metric"`
	Value     float64 `json:"value"`
	Step      int64   `json:"step"`
	Type      string  `json:"counterType"`
	Tags      string  `json:"tags"`
	Timestamp int64   `json:"timestamp"`
}

type RpcParams struct {
	Req    interface{}
	Resp   interface{}
	Method string
}

type Mocker struct {
	Serv    string
	Addr    string
	RpcType string
	Port    int

	Endpoint string
	Interval int
	Batch    int
	Multi    int

	Stat      *Stat
	Pool      *cpool.SafeRpcConnPools
	GetParams func([]*MetricValue, string) (interface{}, interface{}, string)
}

func NewMocker(opts *CmdOpts) *Mocker {
	addrAndPort := fmt.Sprintf("%s:%d", opts.address, opts.port)
	var addrs = []string{addrAndPort}
	var p *cpool.SafeRpcConnPools

	switch opts.rpcType {
	case "rpc":
		p = cpool.CreateSafeRpcConnPools(opts.multi*2, opts.multi, ConnTimeout, CallTimeout, addrs)
	case "grpc":
		p = cpool.CreateSafeGrpcConnPools(opts.multi*2, opts.multi, ConnTimeout, CallTimeout, addrs)
	case "jsonrpc":
		p = cpool.CreateSafeJsonRpcConnPools(opts.multi*2, opts.multi, ConnTimeout, CallTimeout, addrs)
	}

	f := GetParamsTsdb

	mocker := &Mocker{
		Addr:    opts.address,
		Port:    opts.port,
		RpcType: opts.rpcType,

		Endpoint: opts.endpoint,
		Interval: opts.interval,
		Batch:    opts.batch,
		Multi:    opts.multi,

		Stat:      NewStat(),
		Pool:      p,
		GetParams: f,
	}

	return mocker
}

func (m *Mocker) Mock() {
	log.Printf("start mocker, estimate rate: %0.2f/s", float64(1e3*m.Batch*m.Multi)/float64(m.Interval))
	for i := 0; i < m.Multi; i++ {
		go m.mock(i)
	}
}

// non-block
func (m *Mocker) mock(id int) {
	ticker := time.NewTicker(time.Duration(m.Interval) * time.Millisecond)
	for {
		c := <-ticker.C
		endpoint := fmt.Sprintf("%s%02d.%02d.%03d",
			m.Endpoint,
			id,
			c.Second(),
			c.Nanosecond()/1e6/m.Interval,
		)

		items := GenItems(endpoint, m.Batch)
		addrAndPort := fmt.Sprintf("%s:%d", m.Addr, m.Port)
		req, resp, method := m.GetParams(items, m.RpcType)

		err := m.Pool.Call(addrAndPort, method, req, resp)
		if err == nil {
			m.Stat.Incr(int64(m.Batch))
		} else {
			fmt.Printf("%v, %v\n", err, resp)
		}
	}
}

// block
func (m *Mocker) Stats() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		m.Stat.Stats()
	}
}

func GenItems(endpoint string, itemCnt int) []*MetricValue {

	items := make([]*MetricValue, itemCnt, itemCnt)

	for i := 0; i < itemCnt; i++ {
		val := rand.Intn(10000)
		tags := fmt.Sprintf("tag=tag-%03d", i)
		items[i] = &MetricValue{
			Endpoint:  endpoint,
			Metric:    "fortest",
			Timestamp: time.Now().Unix(),
			Step:      int64(60),
			Value:     float64(val),
			Type:      "GAUGE",
			Tags:      tags,
		}
	}
	return items
}
