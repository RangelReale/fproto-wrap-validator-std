package fproto_gowrap_validator_std_pbwrappers

import (
	"fmt"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-std/gowrap/tc/pbwrappers"
	"github.com/RangelReale/fproto-wrap-validator-std/gowrap/scalar"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

//
// UUID
// Validates google.protobuf.XXXValue from google/protobuf/wrappers.proto as the scalar value
//

type TypeValidatorPlugin_PBWrappers struct {
}

func (t *TypeValidatorPlugin_PBWrappers) GetTypeValidator(validatorType *fdep.OptionType, typeinfo fproto_gowrap.TypeInfo, tp *fdep.DepType) fproto_gowrap_validator.TypeValidator {
	// validator.field
	if validatorType.Option != nil &&
		validatorType.Option.DepFile.FilePath == "github.com/RangelReale/fproto-wrap-validator-std/validator.proto" &&
		validatorType.Option.DepFile.ProtoFile.PackageName == "validator" &&
		validatorType.Name == "field" {
		if typeinfo.Converter().TCID() == fprotostd_gowrap_pbwrappers.TCID_PBWRAPPERS {
			return &TypeValidator_PBWrappers{}
		}
	}

	return nil
}

//
// NullTime
//
type TypeValidator_PBWrappers struct {
}

func (v *TypeValidator_PBWrappers) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	// Check only required
	for _, agn := range option.AggregatedSorted() {
		//
		// required
		//
		if agn == "required" {
			if option.AggregatedValues[agn].Source == "true" {
				errors_alias := g.DeclDep("errors", "errors")

				g.P("if !", varSrc, ".Valid {")
				g.In()
				vh.GenerateValidationErrorAdd(g.G(), errors_alias+".New(\"Cannot be blank\")", agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		}
	}

	// only check if valid
	g.P("if ", varSrc, ".Valid {")
	g.In()

	switch tp.Name {
	case "Int64Value", "Int32Value", "UInt64Value", "UInt32Value":
		err := fproto_gowrap_validator_std_scalar.GenerateValidation_scalar_int(g, vh, tp, option, varSrc+".WValue", true)
		if err != nil {
			return err
		}
	case "DoubleValue", "FloatValue":
		err := fproto_gowrap_validator_std_scalar.GenerateValidation_scalar_float(g, vh, tp, option, varSrc+".WValue", true)
		if err != nil {
			return err
		}
	case "StringValue":
		err := fproto_gowrap_validator_std_scalar.GenerateValidation_scalar_string(g, vh, tp, option, varSrc+".WValue", true)
		if err != nil {
			return err
		}
	case "ByteValue":
		err := fproto_gowrap_validator_std_scalar.GenerateValidation_scalar_byte(g, vh, tp, option, varSrc+".WValue", true)
		if err != nil {
			return err
		}
	case "BoolValue":
		err := fproto_gowrap_validator_std_scalar.GenerateValidation_scalar_bool(g, vh, tp, option, varSrc+".WValue", true)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown protobuf wrapper type: %s", tp.Name)
	}

	g.Out()
	g.P("}")

	return nil
}
