package plan

import (
	"encoding/json"
	"fmt"
	"os"
)

// Plan represents the terraform plan JSON structure
type Plan struct {
	FormatVersion    string           `json:"format_version"`
	TerraformVersion string           `json:"terraform_version"`
	PlannedValues    PlannedValues    `json:"planned_values"`
	ResourceChanges  []ResourceChange `json:"resource_changes"`
	PriorState       *State           `json:"prior_state,omitempty"`
}

type PlannedValues struct {
	RootModule Module `json:"root_module"`
}

type Module struct {
	Resources    []Resource `json:"resources,omitempty"`
	ChildModules []Module   `json:"child_modules,omitempty"`
}

type Resource struct {
	Address      string                 `json:"address"`
	Mode         string                 `json:"mode"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	ProviderName string                 `json:"provider_name"`
	Values       map[string]interface{} `json:"values"`
}

type ResourceChange struct {
	Address      string                 `json:"address"`
	Mode         string                 `json:"mode"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	ProviderName string                 `json:"provider_name"`
	Change       Change                 `json:"change"`
}

type Change struct {
	Actions []string               `json:"actions"`
	Before  map[string]interface{} `json:"before"`
	After   map[string]interface{} `json:"after"`
}

type State struct {
	Values StateValues `json:"values"`
}

type StateValues struct {
	RootModule Module `json:"root_module"`
}

// ParsePlanFile reads and parses a terraform plan JSON file
func ParsePlanFile(path string) (*Plan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	return ParsePlanJSON(data)
}

// ParsePlanJSON parses terraform plan JSON data
func ParsePlanJSON(data []byte) (*Plan, error) {
	var plan Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan JSON: %w", err)
	}

	return &plan, nil
}

// GetResourceChanges returns all resource changes from the plan
func (p *Plan) GetResourceChanges() []ResourceChange {
	return p.ResourceChanges
}

// GetCreatedResources returns resources that will be created
func (p *Plan) GetCreatedResources() []ResourceChange {
	var created []ResourceChange
	for _, rc := range p.ResourceChanges {
		for _, action := range rc.Change.Actions {
			if action == "create" {
				created = append(created, rc)
				break
			}
		}
	}
	return created
}

// GetDestroyedResources returns resources that will be destroyed
func (p *Plan) GetDestroyedResources() []ResourceChange {
	var destroyed []ResourceChange
	for _, rc := range p.ResourceChanges {
		for _, action := range rc.Change.Actions {
			if action == "delete" {
				destroyed = append(destroyed, rc)
				break
			}
		}
	}
	return destroyed
}

// GetUpdatedResources returns resources that will be updated in-place
func (p *Plan) GetUpdatedResources() []ResourceChange {
	var updated []ResourceChange
	for _, rc := range p.ResourceChanges {
		for _, action := range rc.Change.Actions {
			if action == "update" {
				updated = append(updated, rc)
				break
			}
		}
	}
	return updated
}

// GetReplacedResources returns resources that will be replaced (destroy + create)
func (p *Plan) GetReplacedResources() []ResourceChange {
	var replaced []ResourceChange
	for _, rc := range p.ResourceChanges {
		hasCreate := false
		hasDelete := false
		for _, action := range rc.Change.Actions {
			if action == "create" {
				hasCreate = true
			}
			if action == "delete" {
				hasDelete = true
			}
		}
		if hasCreate && hasDelete {
			replaced = append(replaced, rc)
		}
	}
	return replaced
}
