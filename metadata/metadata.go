package metadata

import (
	"google.golang.org/grpc/attributes"
)

// metadata common key
const (
	MetaWeight  = "weight"
	MetaCluster = "cluster"
	MetaZone    = "zone"
	MetaColor   = "color"
)

type mdKey struct{}

// MD is context metadata for balancer and resolver
type MD struct {
	Weight int64
}

var mdAttributes = MD{}

// NewAttributes .
func NewAttributes(attr *attributes.Attributes, md MD) *attributes.Attributes {
	if attr == nil {
		attr = attributes.New()
	}
	return attr.WithValues(mdKey{}, md)
}

// FromAttributes .
func FromAttributes(attr *attributes.Attributes) (md MD, ok bool) {
	if attr != nil {
		md, ok = attr.Value(mdKey{}).(MD)
	}
	return
}
