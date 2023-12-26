package resource

type Deployment struct {
	Status          DeploymentStatus   `json:"status"`
	Strategy        string             `json:"strategy"`
	Droplet         Relationship       `json:"droplet"`
	PreviousDroplet Relationship       `json:"previous_droplet"`
	NewProcesses    []ProcessReference `json:"new_processes"`
	Revision        DeploymentRevision `json:"revision"`
	Metadata        *Metadata          `json:"metadata"`
	Relationships   AppRelationship    `json:"relationships"`
	Resource        `json:",inline"`
}

type DeploymentCreate struct {
	Relationships AppRelationship     `json:"relationships"`
	Droplet       *Relationship       `json:"droplet,omitempty"`
	Revision      *DeploymentRevision `json:"revision,omitempty"`
	Strategy      string              `json:"strategy,omitempty"`
	Metadata      *Metadata           `json:"metadata,omitempty"`
}

type DeploymentUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type DeploymentList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*Deployment `json:"resources"`
}

type DeploymentRevision struct {
	GUID    string `json:"guid"`
	Version *int   `json:"version,omitempty"`
}

type ProcessReference struct {
	GUID string `json:"guid"`
	Type string `json:"type"`
}

type DeploymentStatus struct {
	Value   string            `json:"value"`
	Reason  string            `json:"reason"`
	Details map[string]string `json:"details"`
}

func NewDeploymentCreate(appGUID string) *DeploymentCreate {
	return &DeploymentCreate{
		Relationships: AppRelationship{
			App: ToOneRelationship{
				Data: &Relationship{
					GUID: appGUID,
				},
			},
		},
	}
}
