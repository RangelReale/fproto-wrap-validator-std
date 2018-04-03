package fproto_gowrap_validator_std_duration

import (
	"fmt"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-std/gowrap/tc/duration"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

//
// Duration
// Validates google.protobuf.Duration as time.Duration
//

type TypeValidatorPlugin_Duration struct {
}

func (t *TypeValidatorPlugin_Duration) GetTypeValidator(validatorType *fdep.OptionType, typeinfo fproto_gowrap.TypeInfo, tp *fdep.DepType) fproto_gowrap_validator.TypeValidator {
	// validate.field
	if validatorType.Option != nil &&
		validatorType.Option.DepFile.FilePath == "github.com/RangelReale/fproto-wrap-validator-std/validate.proto" &&
		validatorType.Option.DepFile.ProtoFile.PackageName == "validate" &&
		validatorType.Name == "field" {
		if typeinfo.Converter().TCID() == fprotostd_gowrap_duration.TCID_DURATION {
			return &TypeValidator_Duration{}
		}
	}

	return nil
}

//
// Time
//
type TypeValidator_Duration struct {
}

func (v *TypeValidator_Duration) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false

		//
		// xrequired
		//
		if agn == "xrequired" {
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if ", varSrc, " == 0 {")
				g.In()
				vh.GenerateValidationErrorAdd(g.G(), errors_alias+".New(\"Cannot be blank\")", agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	return nil
}
