package middleware

import (
	"log"
	"net/http"
	"time"
)

type MiddleWare struct {
	counter int
}

func (m *MiddleWare) LogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%v] %s - cost: %v\n", r.Method, r.URL, t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func (m *MiddleWare) RecoverMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recover from err: %+v\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (m *MiddleWare) CounterMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		m.counter++
		log.Printf("counter: %v", m.counter)
	}
	return http.HandlerFunc(fn)
}
