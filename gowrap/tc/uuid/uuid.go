package fproto_gowrap_validator_std_uuid

import (
	"fmt"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-std/gowrap/tc/uuid"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

//
// UUID
// Validates fproto_wrap.UUID
//

type TypeValidatorPlugin_UUID struct {
}

func (t *TypeValidatorPlugin_UUID) GetTypeValidator(validatorType *fdep.OptionType, typeinfo fproto_gowrap.TypeInfo, tp *fdep.DepType) fproto_gowrap_validator.TypeValidator {
	// validate.field
	if validatorType.Option != nil &&
		validatorType.Option.DepFile.FilePath == "github.com/RangelReale/fproto-wrap-validator-std/validate.proto" &&
		validatorType.Option.DepFile.ProtoFile.PackageName == "validate" &&
		validatorType.Name == "field" {
		if typeinfo.Converter().TCID() == fprotostd_gowrap_uuid.TCID_UUID {
			return &TypeValidator_UUID{}
		}
		if typeinfo.Converter().TCID() == fprotostd_gowrap_uuid.TCID_NULLUUID {
			return &TypeValidator_NullUUID{}
		}
	}

	return nil
}

//
// UUID
//
type TypeValidator_UUID struct {
}

func (v *TypeValidator_UUID) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	uuid_alias := g.DeclDep("github.com/RangelReale/go.uuid", "uuid")
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false

		//
		// xrequired
		//
		if agn == "xrequired" {
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if ", uuid_alias, ".Equal(", varSrc, ", uuid.Nil) {")
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

func (v *TypeValidator_UUID) GenerateValidationRepeated(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, repeatedType fproto_gowrap_validator.RepeatedType, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	return nil
}

//
// NullUUID
//
type TypeValidator_NullUUID struct {
}

func (v *TypeValidator_NullUUID) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	uuid_alias := g.DeclDep("github.com/RangelReale/go.uuid", "uuid")
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false

		//
		// xrequired
		//
		if agn == "xrequired" {
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if !", varSrc, ".Valid || ", uuid_alias, ".Equal(", varSrc, ".UUID, uuid.Nil) {")
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
