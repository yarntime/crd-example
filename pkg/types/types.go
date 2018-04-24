package types

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type ExamplePhase string

const (
	Creating ExamplePhase = "creating"
	Finished ExamplePhase = "finished"
)

type Example struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata,omitempty"`
	Spec               ExampleSpec   `json:"spec"`
	Status             ExampleStatus `json:"status"`
}

type ExampleSpec struct {
	ResourceType string        `json:"type"`
	Period       time.Duration `json:"period"`
}

type ExampleStatus struct {
	Phase     ExamplePhase  `json:"phase,omitempty"`
	StartTime *meta_v1.Time `json:"startTime,omitempty"`
}

type ExampleList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []Example `json:"items"`
}
