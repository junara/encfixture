package domain

// JobType identifies the kind of media a batch job generates.
type JobType string

// JobType constants enumerate the supported batch job kinds.
const (
	JobTypeImage JobType = "image"
	JobTypeVideo JobType = "video"
	JobTypeAudio JobType = "audio"
)

// Job is a single entry in a batch: a type tag and the config for that type.
// Exactly one of Image, Video, or Audio is populated, matching Type.
type Job struct {
	Type  JobType
	Image *ImageConfig
	Video *VideoConfig
	Audio *AudioConfig
}

// Batch is a collection of media generation jobs.
type Batch struct {
	Jobs []Job
}

// JobResult captures the outcome of running a single job.
type JobResult struct {
	Index  int
	Type   JobType
	Output string
	Err    error
}
