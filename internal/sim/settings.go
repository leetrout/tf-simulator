package sim

// SettingRow mirrors the frontend's Settings rows (kind ∈ pill|color|toggle).
type SettingRow struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
	Value       any    `json:"value"`
	Mono        bool   `json:"mono,omitempty"`
}

// Settings is the /api/sim/settings payload, split into environment + simulator.
type Settings struct {
	Environment []SettingRow `json:"environment"`
	Simulator   []SettingRow `json:"simulator"`
}

// Settings returns the (currently read-only) environment and simulator settings.
func (s *Simulator) Settings() Settings {
	_, _, _, _, meta := s.compute()
	tfVersion := meta.TerraformVersion
	if tfVersion == "" {
		tfVersion = "— no state —"
	}
	return Settings{
		Environment: []SettingRow{
			{Label: "Cloud provider", Description: "fictional provider for the sim", Kind: "pill", Value: meta.Provider},
			{Label: "Region", Description: "default region for new resources", Kind: "pill", Value: meta.Region, Mono: true},
			{Label: "State backend", Description: "where terraform.tfstate lives", Kind: "pill", Value: meta.Backend, Mono: true},
			{Label: "Working directory", Description: "watched for terraform.tfstate", Kind: "pill", Value: meta.WorkDir, Mono: true},
			{Label: "Terraform version", Description: "reported in the status bar", Kind: "pill", Value: tfVersion, Mono: true},
		},
		Simulator: []SettingRow{
			{Label: "Default scenario", Description: "loaded on first launch", Kind: "pill", Value: "In Sync"},
			{Label: "Accent color", Description: "UI highlight color", Kind: "color", Value: "#80c5f6"},
			{Label: "Dependency edges", Description: "dashed links between resources", Kind: "toggle", Value: true},
			{Label: "Inline drift comments", Description: "annotations in the JSON explorer", Kind: "toggle", Value: true},
		},
	}
}
