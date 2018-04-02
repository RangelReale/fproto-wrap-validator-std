package fproto_gowrap_validator_std

import (
	"fmt"
	"strings"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

type ValidatorPlugin_Std struct {
}

func (tp *ValidatorPlugin_Std) GetValidator(validatorType *fdep.OptionType) fproto_gowrap_validator.Validator {
	// validate.field
	if validatorType.Option != nil &&
		validatorType.Option.DepFile.FilePath == "github.com/RangelReale/fproto-wrap-validator-std/validate.proto" &&
		validatorType.Option.DepFile.ProtoFile.PackageName == "validate" &&
		validatorType.Name == "field" {
		return &Validator_Std{validatorType: validatorType}
	}
	return nil
}

func (tp *ValidatorPlugin_Std) ValidatorPrefixes() []string {
	return []string{"validate"}
}

type Validator_Std struct {
	validatorType *fdep.OptionType
}

func (t *Validator_Std) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, varError string) error {
	tinfo := g.G().GetTypeInfo(tp)

	if tinfo.Converter().TCID() == fproto_gowrap.TCID_SCALAR {
		return t.generateValidation_scalar(g, vh, tp, tinfo, option, varSrc, varError)
	}

	tv := vh.GetTypeValidator(t.validatorType, tinfo, tp)

	if tv != nil {
		return tv.GenerateValidation(g, vh, tp, option, varSrc, varError)
	}

	return fmt.Errorf("Unknown type for validator: %s", tp.FullOriginalName())
}

func (t *Validator_Std) generateValidation_scalar(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string, varError string) error {
	errors_alias := g.DeclDep("errors", "errors")

	var opag []string
	for agn, agv := range option.AggregatedValues {
		opag = append(opag, fmt.Sprintf("%s=%s", agn, agv.Source))
	}

	g.P("// ", option.Name, " -- ", option.ParenthesizedName, " ** ", option.NPName, " @@ ", option.Value.Source, " %% ", strings.Join(opag, ", "))

	for agn, agv := range option.AggregatedValues {
		supported := false

		switch *tp.ScalarType {
		//
		// INTEGER
		//
		case fproto.Fixed32Scalar, fproto.Fixed64Scalar, fproto.Int32Scalar, fproto.Int64Scalar,
			fproto.Sfixed32Scalar, fproto.Sfixed64Scalar, fproto.Sint32Scalar, fproto.Sint64Scalar,
			fproto.Uint32Scalar, fproto.Uint64Scalar:
			//
			// xrequired
			//
			if agn == "xrequired" {
				supported = true
				if agv.Source == "true" {
					g.P("if ", varSrc, " == 0 {")
					g.In()
					g.P(varError, " = ", errors_alias, ".New(\"Cannot be blank\")")
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), agn, fproto_gowrap_validator.VEID_REQUIRED)
				}
			}
			//
			// FLOAT
			//
		case fproto.DoubleScalar, fproto.FloatScalar:
			//
			// xrequired
			//
			if agn == "xrequired" {
				supported = true
				if agv.Source == "true" {
					g.P("if ", varSrc, " == 0 {")
					g.In()
					g.P(varError, " = ", errors_alias, ".New(\"Cannot be blank\")")
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), agn, fproto_gowrap_validator.VEID_REQUIRED)
				}
			}
			//
			// STRING
			//
		case fproto.StringScalar:
			//
			// xrequired
			//
			if agn == "xrequired" {
				supported = true
				if agv.Source == "true" {
					g.P("if ", varSrc, " == \"\" {")
					g.In()
					g.P(varError, " = ", errors_alias, ".New(\"Cannot be blank\")")
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), agn, fproto_gowrap_validator.VEID_REQUIRED)
				}
			} else if agn == "length_eq" {
				supported = true
				g.P("if len(", varSrc, ") != ", agv.Source, " {")
				g.In()
				g.P(varError, " = ", errors_alias, ".New(\"Length must be ", agv.Source, "\")")
				g.Out()
				g.P("}")
				vh.GenerateValidationErrorCheck(g.G(), agn, fproto_gowrap_validator.VEID_LENGTH)
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	return nil
}