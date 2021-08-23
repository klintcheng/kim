package kim

import (
	"errors"
	"fmt"
	"sync"

	"github.com/klintcheng/kim/logger"
	"github.com/klintcheng/kim/wire/pkt"
)

var ErrSessionLost = errors.New("err:session lost")

// Router defines
type Router struct {
	middlewares []HandlerFunc
	handlers    *FuncTree
	pool        sync.Pool
}

// NewRouter NewRouter
func NewRouter() *Router {
	r := &Router{
		handlers:    NewTree(),
		middlewares: make([]HandlerFunc, 0),
	}
	r.pool.New = func() interface{} {
		return BuildContext()
	}
	return r
}

func (r *Router) Use(handlers ...HandlerFunc) {
	r.middlewares = append(r.middlewares, handlers...)
}

// Handle regist a commond handler
func (r *Router) Handle(commond string, handlers ...HandlerFunc) {
	r.handlers.Add(commond, r.middlewares...)
	r.handlers.Add(commond, handlers...)
}

// Serve a packet from client
func (r *Router) Serve(packet *pkt.LogicPkt, dispather Dispather, cache SessionStorage, session Session) error {
	if dispather == nil {
		return fmt.Errorf("dispather is nil")
	}
	if cache == nil {
		return fmt.Errorf("cache is nil")
	}
	ctx := r.pool.Get().(*ContextImpl)
	ctx.reset()
	ctx.request = packet
	ctx.Dispather = dispather
	ctx.SessionStorage = cache
	ctx.session = session

	r.serveContext(ctx)
	// Put Context to Pool
	r.pool.Put(ctx)
	return nil
}

func (r *Router) serveContext(ctx *ContextImpl) {
	chain, ok := r.handlers.Get(ctx.Header().Command)
	if !ok {
		ctx.handlers = []HandlerFunc{handleNoFound}
		ctx.Next()
		return
	}
	logger.Info("---- ", chain)
	ctx.handlers = chain
	ctx.Next()
}

func handleNoFound(ctx Context) {
	_ = ctx.Resp(pkt.Status_NotImplemented, &pkt.ErrorResp{Message: "NotImplemented"})
}

// FuncTree is a tree structure
type FuncTree struct {
	nodes map[string]HandlersChain
}

// NewTree NewTree
func NewTree() *FuncTree {
	return &FuncTree{nodes: make(map[string]HandlersChain, 10)}
}

// Add a handler to tree
func (t *FuncTree) Add(path string, handlers ...HandlerFunc) {
	if t.nodes[path] == nil {
		t.nodes[path] = HandlersChain{}
	}

	t.nodes[path] = append(t.nodes[path], handlers...)
}

// Get a handler from tree
func (t *FuncTree) Get(path string) (HandlersChain, bool) {
	f, ok := t.nodes[path]
	return f, ok
}
