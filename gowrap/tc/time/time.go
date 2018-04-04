package fproto_gowrap_validator_std_time

import (
	"fmt"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-std/gowrap/tc/time"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

//
// UUID
// Validates google.protobuf.Timestamp as time.Time
//

type TypeValidatorPlugin_Time struct {
}

func (t *TypeValidatorPlugin_Time) GetTypeValidator(validatorType *fdep.OptionType, typeinfo fproto_gowrap.TypeInfo, tp *fdep.DepType) fproto_gowrap_validator.TypeValidator {
	// validate.field
	if validatorType.Option != nil &&
		validatorType.Option.DepFile.FilePath == "github.com/RangelReale/fproto-wrap-validator-std/validate.proto" &&
		validatorType.Option.DepFile.ProtoFile.PackageName == "validate" &&
		validatorType.Name == "field" {
		if typeinfo.Converter().TCID() == fprotostd_gowrap_time.TCID_TIME {
			return &TypeValidator_Time{}
		}
		if typeinfo.Converter().TCID() == fprotostd_gowrap_time.TCID_NULLTIME {
			return &TypeValidator_NullTime{}
		}
	}

	return nil
}

//
// Time
//
type TypeValidator_Time struct {
}

func (v *TypeValidator_Time) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false

		//
		// required
		//
		if agn == "required" {
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if ", varSrc, ".IsZero() {")
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

//
// NullTime
//
type TypeValidator_NullTime struct {
}

func (v *TypeValidator_NullTime) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false

		//
		// required
		//
		if agn == "required" {
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if !", varSrc, ".Valid || ", varSrc, ".IsZero() {")
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
