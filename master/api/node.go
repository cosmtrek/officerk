package api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/cosmtrek/officerk/models"
	"github.com/cosmtrek/officerk/services"
	"github.com/cosmtrek/officerk/utils/api"
)

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
