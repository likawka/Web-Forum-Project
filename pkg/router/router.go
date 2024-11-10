package router

import (
	"context"
	"net/http"
	"regexp"
	"sync"

	"golang.org/x/time/rate"
)

type route struct {
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type Router struct {
	routes      []route
	limiter     *rate.Limiter
	limiterLock sync.Mutex
}

func (router *Router) SetRateLimiter(rps float64, burst int) {
	router.limiterLock.Lock()
	defer router.limiterLock.Unlock()
	router.limiter = rate.NewLimiter(rate.Limit(rps), burst)
}

func (router *Router) NewRoute(regexpString string, handler http.HandlerFunc) {
	regex := regexp.MustCompile("^" + regexpString + "$")
	router.routes = append(router.routes, route{
		regex,
		handler,
	})
}

func (router *Router) Serve(w http.ResponseWriter, r *http.Request) {
	router.limiterLock.Lock()
	limiter := router.limiter
	router.limiterLock.Unlock()

	if limiter != nil && !limiter.Allow() {
		http.Redirect(w, r, "/Too many requests", http.StatusFound)
		return
	}

	for _, v := range router.routes {
		matches := v.regex.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			matchMap := make(map[string]string)
			groupNames := v.regex.SubexpNames()

			for i := 1; i < len(matches); i++ {
				matchMap[groupNames[i]] = matches[i]
			}

			ctx := context.WithValue(r.Context(), struct{}{}, matchMap)
			v.handler(w, r.WithContext(ctx))
			return
		}
	}
}

func StartServer(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, nil)
}
