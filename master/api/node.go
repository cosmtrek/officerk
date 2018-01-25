package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/cosmtrek/officerk/models"
	"github.com/cosmtrek/officerk/services"
	"github.com/cosmtrek/officerk/utils/api"
)

var nodeKey = ctxKey("node")

// NodeRequest ...
type NodeRequest struct {
	*models.Node
}

// Bind for post-processing NodeRequest
func (n *NodeRequest) Bind(r *http.Request) error {
	return nil
}

// NodeResponse ...
type NodeResponse struct {
	*models.Node
	Online bool `json:"online"`
}

// NewNodeResponse ...
func NewNodeResponse(node *models.Node, online bool) *NodeResponse {
	return &NodeResponse{Node: node, Online: online}
}

// Render for NodeResponse
func (nr NodeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// NodeListResponse ...
type NodeListResponse []*NodeResponse

// NewNodeListResponse ...
func NewNodeListResponse(nodes []*models.Node, onlineNodes []string) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, node := range nodes {
		online := false
		for _, n := range onlineNodes {
			if node.IP == n {
				online = true
			}
		}
		list = append(list, NewNodeResponse(node, online))
	}
	return list
}

// NodeCtx finds node
func NodeCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		id := chi.URLParam(r, "nodeID")
		if id == "" {
			render.Render(w, r, api.ErrNotFound)
			return
		}
		node := new(models.Node)
		if err = services.GetNode(db, id, node); err != nil {
			render.Render(w, r, api.ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), nodeKey, node)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ListNodes return nodes
func (h *Handler) ListNodes(w http.ResponseWriter, r *http.Request) {
	var err error
	nodes, err := services.GetNodes(db)
	if err != nil {
		render.Render(w, r, api.ErrNotFound)
		return
	}
	render.Render(w, r, api.OK(NewNodeListResponse(nodes, h.runtime.Nodes())))
}

// CreateNode creates node
func (h *Handler) CreateNode(w http.ResponseWriter, r *http.Request) {
	var err error
	data := &NodeRequest{}
	if err = render.Bind(r, data); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = services.CreateNode(db, data.Node); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, api.OK(NewNodeResponse(data.Node, h.runtime.IsOnline(data.Node.IP))))
}

// ListOnlineNodes return nodes live in etcd
func (h *Handler) ListOnlineNodes(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, api.OK(h.runtime.Nodes()))
}

// GetNode gets node
func (h *Handler) GetNode(w http.ResponseWriter, r *http.Request) {
	node := r.Context().Value(nodeKey).(*models.Node)
	render.Render(w, r, api.OK(NewNodeResponse(node, h.runtime.IsOnline(node.IP))))
}

// UpdateNode updates node
func (h *Handler) UpdateNode(w http.ResponseWriter, r *http.Request) {
	var err error
	node := r.Context().Value(nodeKey).(*models.Node)
	data := &NodeRequest{}
	if err = render.Bind(r, data); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = services.UpdateNode(db, node, data.Node); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	if err = services.GetNode(db, strconv.Itoa(int(node.ID)), data.Node); err != nil {
		render.Render(w, r, api.ErrNotFound)
		return
	}
	render.Render(w, r, api.OK(NewNodeResponse(data.Node, h.runtime.IsOnline(data.Node.IP))))
}

// DeleteNode deletes node
func (h *Handler) DeleteNode(w http.ResponseWriter, r *http.Request) {
	var err error
	node := r.Context().Value(nodeKey).(*models.Node)
	if err = services.DeleteNode(db, node); err != nil {
		render.Render(w, r, api.ErrInvalidRequest(err))
		return
	}
	// TODO: cancel the jobs running on this node
	render.Render(w, r, api.OK("{}"))
}
