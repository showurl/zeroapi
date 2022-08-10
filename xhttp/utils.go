package xhttp

import "net/http"

const (
	xForwardedFor = "X-Forwarded-For"
	maxMemory     = 32 << 20 // 32MB
)

// GetFormValues returns the form values.
func GetFormValues(r *http.Request) (map[string]interface{}, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return nil, err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	return params, nil
}

// GetRemoteAddr returns the peer address, supports X-Forward-For.
func GetRemoteAddr(r *http.Request) string {
	v := r.Header.Get(xForwardedFor)
	if len(v) > 0 {
		return v
	}

	return r.RemoteAddr
}

func GetRequestIP(r *http.Request) string {
	addr := GetRemoteAddr(r)
	if addr == "" {
		realIp := r.Header.Get("X-Real-IP")
		if realIp != "" {
			return realIp
		}
	}
	return addr
}
