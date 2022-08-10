package zeroapi

import (
	"encoding/json"
	"fmt"
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto" //lint:ignore SA1019 we have to import this because it appears in exported API
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type ResponseHandler func(in proto.Message) (code int, msg string, data interface{})

type handlerOption struct {
	responseHandler ResponseHandler
}

type OptionFunc func(*handlerOption)

func WithResponseHandler(responseHandler ResponseHandler) OptionFunc {
	return func(o *handlerOption) {
		if responseHandler != nil {
			o.responseHandler = responseHandler
		}
	}
}

type Handler struct {
	*grpcurl.DefaultEventHandler
	responseHandler ResponseHandler
}

func NewHandler(w http.ResponseWriter, source grpcurl.DescriptorSource, optionFs ...OptionFunc) *Handler {
	h := &Handler{
		DefaultEventHandler: &grpcurl.DefaultEventHandler{
			Out:       w,
			Formatter: grpcurl.NewJSONFormatter(true, grpcurl.AnyResolverFromDescriptorSource(source)),
		},
	}
	option := &handlerOption{responseHandler: h.defaultResponseHandler}
	for _, optionF := range optionFs {
		optionF(option)
	}
	h.responseHandler = option.responseHandler
	return h
}

func (h *Handler) OnReceiveResponse(resp proto.Message) {
	h.NumResponses++
	code, msg, data := h.responseHandler(resp)
	res := XResponse{
		Code:       int32(code),
		Msg:        msg,
		Data:       data,
		ServerTime: time.Now().UnixMilli(),
	}
	buf, _ := json.Marshal(res)
	fmt.Fprintln(h.Out, string(buf))
}

func (h *Handler) OnReceiveTrailers(stat *status.Status, md metadata.MD) {
	h.Status = stat
	if stat.Code() != codes.OK {
		respJsonStr := `{"msg":"服务繁忙，请稍后再试","code":-1}`
		fmt.Fprintln(h.Out, respJsonStr)
	}
}

func (h *Handler) defaultResponseHandler(in proto.Message) (code int, msg string, data interface{}) {
	if in == nil {
		return 0, "", nil
	} else {
		message, ok := in.(*dynamic.Message)
		if ok {
			if cod, exist := protoMessageValue(message, "errCode", 0); exist {
				code = int(InterfaceToInt64(cod))
			}
			if failedReason, exist := protoMessageValue(message, "failedReason", ""); exist && failedReason != "" {
				msg = failedReason.(string)
				if code == 0 {
					code = -1
				}
			}
			data = in
			return
		}
		return 0, "", in
	}
}

func protoMessageValue(message *dynamic.Message, key string, defaultValue interface{}) (interface{}, bool) {
	value, err := message.TryGetFieldByName(key)
	if err != nil {
		return defaultValue, false
	}
	if value == nil {
		return defaultValue, true
	}
	return value, true
}

func BuildHandler(resp interface{}) ResponseHandler {
	return func(in proto.Message) (code int, msg string, data interface{}) {
		if in == nil {
			return 0, "", nil
		}
		message, ok := in.(*dynamic.Message)
		if ok {
			if cod, exist := protoMessageValue(message, "errCode", 0); exist {
				code = int(InterfaceToInt64(cod))
			}
			if failedReason, exist := protoMessageValue(message, "failedReason", ""); exist && failedReason != "" {
				msg = failedReason.(string)
				if code == 0 {
					code = -1
				}
			}
		}
		bytes, err := proto.Marshal(in)
		if err != nil {
			return -1, err.Error(), nil
		}
		res := Copy(resp)
		if m, ok := res.(proto.Message); ok {
			err = proto.Unmarshal(bytes, m)
		} else {
			err = fmt.Errorf("response type %T is not proto.Message", res)
		}
		if err != nil {
			return -1, err.Error(), nil
		}
		return code, msg, res
	}
}

func WithBuildHandler(resp interface{}) OptionFunc {
	return func(o *handlerOption) {
		o.responseHandler = BuildHandler(resp)
	}
}
