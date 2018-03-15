package main

import (
	tsdbPb "github.com/huangaz/tsdb/protobuf"

	"strings"
)

func GetParamsTsdb(items []*MetricValue, rpcType string) (req, resp interface{}, method string) {

	switch rpcType {
	case "grpc":
		pbTsdbItems := make([]*tsdbPb.DataPoint, len(items), len(items))
		for i, item := range items {
			pbTsdbItems[i] = convert2PbTsdbItem(item)
		}

		req = &tsdbPb.PutRequest{Datas: pbTsdbItems}
		resp = &tsdbPb.PutResponse{}
		method = "/tsdbPb.Tsdb/Put"
	case "rpc":
		TsdbItems := make([]*tsdbPb.DataPoint, len(items), len(items))
		for i, item := range items {
			TsdbItems[i] = convert2PbTsdbItem(item)
		}

		req = &tsdbPb.PutRequest{Datas: TsdbItems}
		resp = &tsdbPb.PutResponse{}
		method = "RpcServer.Put"
	case "jsonrpc":
		TsdbItems := make([]*tsdbPb.DataPoint, len(items), len(items))
		for i, item := range items {
			TsdbItems[i] = convert2PbTsdbItem(item)
		}

		req = &tsdbPb.PutRequest{Datas: TsdbItems}
		resp = &tsdbPb.PutResponse{}
		method = "JsonRpcServer.Put"
	}

	return
}

func convert2PbTsdbItem(m *MetricValue) *tsdbPb.DataPoint {
	// s := []string{m.Endpoint, m.Metric, m.Type, m.Tags}
	s := []string{"for", "test"}
	newKey := strings.Join(s, "/")

	return &tsdbPb.DataPoint{
		Key: &tsdbPb.Key{
			Key:     []byte(newKey),
			ShardId: 1,
		},
		Value: &tsdbPb.TimeValuePair{
			Timestamp: m.Timestamp,
			Value:     m.Value,
		},
	}
}
