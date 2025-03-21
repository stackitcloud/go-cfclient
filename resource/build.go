package resource

type BuildState string

func (b BuildState) String() string {
	return string(b)
}

// The 3 lifecycle states
const (
	BuildStateStaging BuildState = "STAGING"
	BuildStateStaged  BuildState = "STAGED"
	BuildStateFailed  BuildState = "FAILED"
)

type Build struct {
	State BuildState `json:"state"`
	Error *string    `json:"error"`

	StagingMemoryInMB                 int `json:"staging_memory_in_mb"`
	StagingDiskInMB                   int `json:"staging_disk_in_mb"`
	StagingLogRateLimitBytesPerSecond int `json:"staging_log_rate_limit_bytes_per_second"`

	Lifecycle     Lifecycle       `json:"lifecycle"`
	Package       Relationship    `json:"package"`
	Droplet       *Relationship   `json:"droplet"`
	CreatedBy     CreatedBy       `json:"created_by"`
	Relationships AppRelationship `json:"relationships"`
	Metadata      *Metadata       `json:"metadata"`
	Resource      `json:",inline"`
}

type BuildCreate struct {
	Package                           Relationship `json:"package"`
	Lifecycle                         *Lifecycle   `json:"lifecycle,omitempty"`
	StagingMemoryInMB                 int          `json:"staging_memory_in_mb,omitempty"`
	StagingDiskInMB                   int          `json:"staging_disk_in_mb,omitempty"`
	StagingLogRateLimitBytesPerSecond int          `json:"staging_log_rate_limit_bytes_per_second,omitempty"`
	Metadata                          *Metadata    `json:"metadata,omitempty"`
}

type BuildUpdate struct {
	Metadata  *Metadata  `json:"metadata,omitempty"`
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
	State     string     `json:"state,omitempty"`
}

type BuildList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Build   `json:"resources"`
}

type CreatedBy struct {
	GUID  string `json:"guid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewBuildCreate(packageGUID string) *BuildCreate {
	return &BuildCreate{
		Package: Relationship{
			GUID: packageGUID,
		},
	}
}

func NewBuildUpdate() *BuildUpdate {
	return &BuildUpdate{}
}
