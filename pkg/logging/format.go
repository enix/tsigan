package logging

import (
	"github.com/enix/tsigan/pkg/utils"
)

type Format string

const (
	SimpleFormat     Format = "simple"
	StructuredFormat Format = "structured"
	JSONFormat       Format = "json"
	DeveloperFormat  Format = "developer"
)

type FormatFlag struct {
	*utils.Enum
	defaultValue Format
}

func NewServerFormatFlag(defaultValue Format) *FormatFlag {
	return &FormatFlag{
		Enum: utils.NewEnum(
			string(defaultValue),
			string(StructuredFormat),
			string(JSONFormat),
			string(DeveloperFormat),
		),
		defaultValue: defaultValue,
	}
}
