package fproto_gowrap_validator_std_duration

import (
	"fmt"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-std/gowrap/tc/duration"
	"github.com/RangelReale/fproto-wrap-validator-std/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

//
// Duration
// Validates google.protobuf.Duration as time.Duration
//

type DefaultTypeValidatorPlugin_Duration struct {
}

func (t *DefaultTypeValidatorPlugin_Duration) GetDefaultTypeValidator(typeinfo fproto_gowrap.TypeInfo, tp *fdep.DepType) fproto_gowrap_validator_std.DefaultTypeValidator {
	if typeinfo.Converter().TCID() == fprotostd_gowrap_duration.TCID_DURATION {
		return &DefaultTypeValidator_Duration{}
	}

	return nil
}

//
// Time
//
type DefaultTypeValidator_Duration struct {
}

func (v *DefaultTypeValidator_Duration) GenerateValidation(g *fproto_gowrap.GeneratorFile, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, varError string) error {
	errors_alias := g.DeclDep("errors", "errors")

	for agn, agv := range option.AggregatedValues {
		supported := false

		//
		// xrequired
		//
		if agn == "xrequired" {
			supported = true
			if agv.Source == "true" {
				g.P("if ", varSrc, " == 0 {")
				g.In()
				g.P("err = ", errors_alias, ".New(\"Cannot be blank\")")
				g.Out()
				g.P("}")
				g.GenerateSimpleErrorCheck()
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	return nil
}
