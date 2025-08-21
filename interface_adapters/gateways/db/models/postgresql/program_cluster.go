package postgresql

import (
	"time"

	"github.com/jackc/pgtype"
)

// ProgramCluster represents a group of related programs with common keywords
type ProgramCluster struct {
	ID         *int64           `json:"id" gorm:"primaryKey"`
	ProgramIDs pgtype.Int4Array `json:"program_ids" gorm:"type:integer[]"`
	Keywords   pgtype.TextArray `json:"keywords" gorm:"type:text[]"`
	Metadata   pgtype.JSONB     `json:"metadata" gorm:"type:jsonb"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

type Metadata struct {
	Source     string  `json:"source"`
	Confidence float64 `json:"confidence"`
}

// TableName specifies the database table name
func (ProgramCluster) TableName() string {
	return "programs.program_clusters"
}

// GetProgramIDs convierte pgtype.Int4Array a []int
func (pc *ProgramCluster) GetProgramIDs() []int {
	if pc.ProgramIDs.Status != pgtype.Present {
		return []int{}
	}

	result := make([]int, len(pc.ProgramIDs.Elements))
	for i, elem := range pc.ProgramIDs.Elements {
		if elem.Status == pgtype.Present {
			result[i] = int(elem.Int)
		}
	}
	return result
}

// GetKeywords convierte pgtype.TextArray a []string
func (pc *ProgramCluster) GetKeywords() []string {
	if pc.Keywords.Status != pgtype.Present {
		return []string{}
	}

	result := make([]string, len(pc.Keywords.Elements))
	for i, elem := range pc.Keywords.Elements {
		if elem.Status == pgtype.Present {
			result[i] = elem.String
		}
	}
	return result
}

// GetMetadata convierte pgtype.JSONB a Metadata
func (pc *ProgramCluster) GetMetadata() Metadata {
	if pc.Metadata.Status != pgtype.Present {
		return Metadata{}
	}

	var metadata Metadata
	if err := pc.Metadata.AssignTo(&metadata); err != nil {
		return Metadata{}
	}
	return metadata
}
