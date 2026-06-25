package sim

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// State is a parsed terraform.tfstate (format version 4), reduced to what the
// simulator needs.
type State struct {
	Version          int
	TerraformVersion string
	Serial           int
	Resources        []StateResource
}

// StateResource is one managed resource instance from the state file.
type StateResource struct {
	Address    string            `json:"address"`    // e.g. nimbus_subnet.public
	Module     string            `json:"module"`     // e.g. module.net ("" for root)
	Type       string            `json:"type"`       // nimbus_subnet
	Name       string            `json:"name"`       // public
	ID         string            `json:"id"`         // cloud id from attributes.id
	Attributes map[string]string `json:"attributes"` // scalar attributes as strings
}

// raw mirrors the documented v4 schema subset.
type rawState struct {
	Version          int    `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial           int    `json:"serial"`
	Resources        []struct {
		Module    string `json:"module"`
		Mode      string `json:"mode"`
		Type      string `json:"type"`
		Name      string `json:"name"`
		Instances []struct {
			IndexKey   any                        `json:"index_key"`
			Attributes map[string]json.RawMessage `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

// ParseState parses raw v4 state JSON. Only managed nimbus_* resources are kept.
func ParseState(data []byte) (*State, error) {
	var rs rawState
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, err
	}
	st := &State{Version: rs.Version, TerraformVersion: rs.TerraformVersion, Serial: rs.Serial}
	for _, res := range rs.Resources {
		if res.Mode != "" && res.Mode != "managed" {
			continue
		}
		if !strings.HasPrefix(res.Type, "nimbus_") {
			continue
		}
		for _, inst := range res.Instances {
			attrs := flattenScalars(inst.Attributes)
			sr := StateResource{
				Address:    buildAddress(res.Module, res.Type, res.Name, inst.IndexKey),
				Module:     res.Module,
				Type:       res.Type,
				Name:       res.Name,
				ID:         attrs["id"],
				Attributes: attrs,
			}
			st.Resources = append(st.Resources, sr)
		}
	}
	return st, nil
}

// LoadState reads and parses a state file. A missing file yields (nil, nil) so
// callers can treat "no state yet" as an empty topology.
func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	return ParseState(data)
}

// flattenScalars converts the raw attribute map into string values for the scalar
// attributes the simulator compares (skips nested objects/arrays).
func flattenScalars(raw map[string]json.RawMessage) map[string]string {
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		s := strings.TrimSpace(string(v))
		if s == "" || s == "null" {
			continue
		}
		switch s[0] {
		case '{', '[':
			continue // skip composite attributes
		case '"':
			var str string
			if err := json.Unmarshal(v, &str); err == nil {
				out[k] = str
			}
		default:
			out[k] = s // number/bool, keep raw token
		}
	}
	return out
}

// buildAddress reconstructs the Terraform resource address, including module
// prefix and index key suffix when present.
func buildAddress(module, typ, name string, indexKey any) string {
	addr := typ + "." + name
	if module != "" {
		addr = module + "." + addr
	}
	switch k := indexKey.(type) {
	case string:
		addr += fmt.Sprintf("[%q]", k)
	case float64:
		addr += fmt.Sprintf("[%d]", int(k))
	}
	return addr
}
