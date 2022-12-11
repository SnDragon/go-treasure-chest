package shorturl

import (
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/middleware"
	"github.com/SnDragon/go-treasure-chest/internal/app/shorturl/resource"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (app *App) Initialize() error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// 初始化资源
	if err := resource.InitResource(); err != nil {
		return err
	}
	app.Router = mux.NewRouter()
	mw := &middleware.MiddleWare{}
	app.Router.Use(mw.LogMiddleware, mw.RecoverMiddleware)
	app.InitializeRouter()
	return nil
}

func (app *App) InitializeRouter() {
	// 将长地址转成短地址
	app.Router.Handle("/api/shorten", http.HandlerFunc(app.createShortUrl)).Methods(http.MethodPost)
	app.Router.Handle("/api/info", http.HandlerFunc(app.getUrlInfo)).Methods(http.MethodGet)
	app.Router.Handle("/{sid:[A-Za-z0-9]+}", http.HandlerFunc(app.redirect)).Methods(http.MethodGet)
	//mw := &MiddleWare{}
	//m := alice.New(mw.CounterMiddleware)
	//app.Router.Handle("/api/counter", m.ThenFunc(app.counter)).Methods(http.MethodGet)

}

func (app *App) Run(addr string) {
	log.Printf("App run in %v...\n", addr)
	log.Fatalln(http.ListenAndServe(addr, app.Router))
}
