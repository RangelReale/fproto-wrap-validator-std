package fproto_gowrap_validator_std

import (
	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

type DefaultTypeValidatorPlugin interface {
	// Returns a default type validator for the type
	GetDefaultTypeValidator(typeinfo fproto_gowrap.TypeInfo, tp *fdep.DepType) DefaultTypeValidator
}

type DefaultTypeValidator interface {
	GenerateValidation(g *fproto_gowrap.GeneratorFile, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, varError string) error
}
