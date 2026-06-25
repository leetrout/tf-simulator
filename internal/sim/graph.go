package sim

import "github.com/leetrout/terraform-sim/internal/cloud"

// GraphNode is a jsongraph-shaped node. Extra fields (type, kind, status, ...) are
// passed through to the renderer, which colours nodes by status.
type GraphNode struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Type    string `json:"type"`
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address,omitempty"`
}

// GraphEdge is a jsongraph-shaped edge. kind="phantom" marks edges touching a
// phantom node (rendered dashed by the frontend).
type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label,omitempty"`
	Kind   string `json:"kind,omitempty"`
}

// GraphData is the {nodes, edges} payload consumed by jsongraph.
type GraphData struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// shortKind maps a Terraform type to the compact badge used in the graph.
func shortKind(t string) string {
	switch t {
	case "nimbus_network":
		return "NETWORK"
	case "nimbus_subnet":
		return "SUBNET"
	case "nimbus_network_interface":
		return "INTERFACE"
	case "nimbus_vm_instance":
		return "VM INSTANCE"
	case "nimbus_storage_bucket":
		return "BUCKET"
	}
	return t
}

// Graph builds a jsongraph view. view "cloud" includes resources present in the
// cloud, "state" includes resources present in state; anything else is the union.
func (s *Simulator) Graph(view string) GraphData {
	_, _, _, entries, _ := s.compute()

	include := func(e entry) bool {
		switch view {
		case "cloud":
			return e.InCloud
		case "state":
			return e.InState
		default:
			return true
		}
	}

	present := map[string]entry{}
	g := GraphData{Nodes: []GraphNode{}, Edges: []GraphEdge{}}
	for _, e := range entries {
		if !include(e) {
			continue
		}
		present[e.ID] = e
		label := e.Type + "." + e.Name
		if e.Address != "" {
			label = e.Address
		}
		g.Nodes = append(g.Nodes, GraphNode{
			ID:      e.ID,
			Label:   label,
			Type:    e.Type,
			Kind:    shortKind(e.Type),
			Name:    e.Name,
			Status:  e.Status,
			Address: e.Address,
		})
	}

	for _, e := range entries {
		if !include(e) {
			continue
		}
		for _, pf := range cloud.ParentFields(e.Type) {
			parentID := e.Attrs[pf.Field]
			if parentID == "" {
				continue
			}
			parent, ok := present[parentID]
			if !ok {
				continue
			}
			kind := ""
			if e.Status == StatusPhantom || parent.Status == StatusPhantom {
				kind = "phantom"
			}
			g.Edges = append(g.Edges, GraphEdge{
				Source: e.ID,
				Target: parentID,
				Label:  pf.Field,
				Kind:   kind,
			})
		}
	}
	return g
}
