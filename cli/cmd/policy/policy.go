// Package policy peforms the genreic session key management
package policy

import (
	"encoding/json"
	"pie-rum/cmd/keys"

	"github.com/spf13/cobra"
)

type InitPolicy struct {
	CMDORG *IOrganization `json:"cmdOrg"`
}

type IOrganization struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	ShortHand string `json:"shortHand"`
	Command   string `json:"command"`
}

func (p *InitPolicy) Pack() []byte {
	e, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return e
}

func PackInitPolicy(cmd *cobra.Command) {
	var p InitPolicy

	cmd.Flags().Set(keys.Init_Flag, string(p.Pack()))
}

func NewPolicyPack() *InitPolicy {
	return &InitPolicy{
		&IOrganization{},
	}
}
func UnPackInitPolicy(cmd *cobra.Command) (*InitPolicy, error) {
	f, err := cmd.Flags().GetString(keys.Init_Flag)
	if err != nil {
		return nil, err
	}

	// 1. Allocate the struct
	src := NewPolicyPack()

	// 2. Unmarshal into the allocated pointer
	if err := json.Unmarshal([]byte(f), src); err != nil {
		return nil, err
	}
	return src, nil
}
