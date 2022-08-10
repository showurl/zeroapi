package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/showurl/zeroapi/xhttp"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

var (
	ParamErr = fmt.Errorf("param error")
)

// NewRequestParser creates a new request parser from the given http.Request and resolver.
func NewRequestParser(r *http.Request, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	return NewJSONRequestParser(r, resolver)
}

type JSONRequestParser struct {
	r            *http.Request
	dec          *json.Decoder
	unmarshaler  jsonpb.Unmarshaler
	requestCount int
	selfId       string
	platform     string
	ip           string
}

func NewJSONRequestParser(r *http.Request, resolver jsonpb.AnyResolver) (*JSONRequestParser, error) {
	var in io.Reader

	vars := pathvar.Vars(r)
	params, err := xhttp.GetFormValues(r)
	if err != nil {
		return nil, err
	}
	for k, v := range vars {
		params[k] = v
	}
	if len(params) == 0 {
		in = r.Body
	} else {
		if r.ContentLength == 0 {
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(params); err != nil {
				logx.WithContext(r.Context()).Errorf("Encode error: %s", err)
				return nil, err
			}
			in = &buf
		} else {
			m := make(map[string]interface{})
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				logx.WithContext(r.Context()).Errorf("Decode error: %s", err)
				return nil, err
			}

			for k, v := range params {
				m[k] = v
			}
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(params); err != nil {
				logx.WithContext(r.Context()).Errorf("Encode error: %s", err)
				return nil, err
			}
			in = &buf
		}
	}

	return &JSONRequestParser{
		r:           r,
		dec:         json.NewDecoder(in),
		unmarshaler: jsonpb.Unmarshaler{AnyResolver: resolver},
		selfId:      r.Header.Get("user_id"),
		platform:    r.Header.Get("platform"),
		ip:          xhttp.GetRequestIP(r),
	}, nil
}

func (f *JSONRequestParser) Next(m proto.Message) error {
	if dm, ok := m.(*dynamic.Message); ok {
		_ = dm.TrySetFieldByName("selfId", f.selfId)
		_ = dm.TrySetFieldByName("selfID", f.selfId)
		_ = dm.TrySetFieldByName("SelfID", f.selfId)
		_ = dm.TrySetFieldByName("platform", f.platform)
		_ = dm.TrySetFieldByName("Platform", f.platform)
		_ = dm.TrySetFieldByName("ip", f.ip)
		_ = dm.TrySetFieldByName("IP", f.ip)
		_ = dm.TrySetFieldByName("Ip", f.ip)
	}
	var msg json.RawMessage
	if err := f.dec.Decode(&msg); err != nil {
		return err
	}
	f.requestCount++
	err := f.unmarshaler.Unmarshal(bytes.NewReader(msg), m)
	if err != nil {
		logx.WithContext(f.r.Context()).Errorf("unmarshal error: %s", err)
		return ParamErr
	}
	return err
}

func (f *JSONRequestParser) NumRequests() int {
	return f.requestCount
}
