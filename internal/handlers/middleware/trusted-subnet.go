package middleware

import (
	"net"
	"net/http"
	"strings"
)

func TrustedSubnet(subnet *net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if subnet == nil {
				next.ServeHTTP(w, r)
				return
			}

			contentType := r.Header.Get("Content-Type")

			remoteAddr := r.Header.Get("X-Real-IP")
			if remoteAddr == "" {
				ra := strings.Split(r.RemoteAddr, ":")
				if len(ra) == 2 {
					remoteAddr = ra[0]
				}
			}

			ip := net.ParseIP(remoteAddr)

			if !subnet.Contains(ip) {
				if strings.Contains(contentType, "application/json") {
					JSONError(w, "The request from this ip-address was rejected.", http.StatusForbidden)
				} else {
					http.Error(w, "The request from this ip-address was rejected.", http.StatusForbidden)
				}
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
