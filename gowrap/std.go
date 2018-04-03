package fproto_gowrap_validator_std

import (
	"fmt"
	"strconv"
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
		validatorType.Option.DepFile.ProtoFile.PackageName == "validate" {
		if validatorType.Name == "field" {
			return &Validator_Std{validatorType: validatorType}
		} else if validatorType.Name == "rfield" {
			return &Validator_Std_Repeated{validatorType: validatorType}
		}
	}
	return nil
}

func (tp *ValidatorPlugin_Std) ValidatorPrefixes() []string {
	return []string{"validate"}
}

//
// Validator_Std
//
type Validator_Std struct {
	validatorType *fdep.OptionType
}

func (tp *Validator_Std) FPValidator() {

}

func (t *Validator_Std) GenerateValidation(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	if option.ParenthesizedName != "validate.field" {
		return nil
	}

	tinfo := g.G().GetTypeInfo(tp)

	if tinfo.Converter().TCID() == fproto_gowrap.TCID_SCALAR {
		return t.generateValidation_scalar(g, vh, tp, tinfo, option, varSrc)
	}

	tv := vh.GetTypeValidator(t.validatorType, tinfo, tp)

	if tv != nil {
		return tv.GenerateValidation(g, vh, tp, option, varSrc)
	}

	return fmt.Errorf("Unknown type for validator: %s", tp.FullOriginalName())
}

func (t *Validator_Std) generateValidation_scalar(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	var opag []string
	for _, agn := range option.AggregatedSorted() {
		opag = append(opag, fmt.Sprintf("%s=%s", agn, option.AggregatedValues[agn].Source))
	}

	g.P("// ", option.Name, " -- ", option.ParenthesizedName, " ** ", option.NPName, " @@ ", option.Value.Source, " %% ", strings.Join(opag, ", "))

	switch *tp.ScalarType {
	case fproto.Fixed32Scalar, fproto.Fixed64Scalar, fproto.Int32Scalar, fproto.Int64Scalar,
		fproto.Sfixed32Scalar, fproto.Sfixed64Scalar, fproto.Sint32Scalar, fproto.Sint64Scalar,
		fproto.Uint32Scalar, fproto.Uint64Scalar:
		//
		// INTEGER
		//
		return t.generateValidation_scalar_int(g, vh, tp, tinfo, option, varSrc)
	case fproto.DoubleScalar, fproto.FloatScalar:
		//
		// FLOAT
		//
		return t.generateValidation_scalar_float(g, vh, tp, tinfo, option, varSrc)
	case fproto.StringScalar:
		//
		// STRING
		//
		return t.generateValidation_scalar_string(g, vh, tp, tinfo, option, varSrc)
	}

	return fmt.Errorf("Validation not supported for type %s", tp.FullOriginalName())
}

func (t *Validator_Std) generateValidation_scalar_int(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false
		if agn == "xrequired" {
			//
			// xrequired
			//
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if ", varSrc, " == 0 {")
				g.In()
				vh.GenerateValidationErrorAdd(g.G(), errors_alias+".New(\"Cannot be blank\")", agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if agn == "int_gt" {
			//
			// int_gt
			//
			supported = true
			g.P("if ", varSrc, " <= ", option.AggregatedValues[agn].Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must be greater than %s")`, errors_alias, option.AggregatedValues[agn].Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_MINMAX, "int_gt", option.AggregatedValues[agn].Source)
			g.Out()
			g.P("}")
		} else if agn == "int_lt" {
			//
			// int_lt
			//
			supported = true
			g.P("if ", varSrc, " >= ", option.AggregatedValues[agn].Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must be lower than %s")`, errors_alias, option.AggregatedValues[agn].Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_MINMAX, "int_lt", option.AggregatedValues[agn].Source)
			g.Out()
			g.P("}")
		} else if agn == "int_gte" {
			//
			// int_gte
			//
			supported = true
			g.P("if ", varSrc, " < ", option.AggregatedValues[agn].Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must be greater or equals to %s")`, errors_alias, option.AggregatedValues[agn].Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_MINMAX, "int_gte", option.AggregatedValues[agn].Source)
			g.Out()
			g.P("}")
		} else if agn == "int_lte" {
			//
			// int_lte
			//
			supported = true
			g.P("if ", varSrc, " > ", option.AggregatedValues[agn].Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must be lower or equals to %s")`, errors_alias, option.AggregatedValues[agn].Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_MINMAX, "int_lte", option.AggregatedValues[agn].Source)
			g.Out()
			g.P("}")
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	return nil
}

func (t *Validator_Std) generateValidation_scalar_float(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	// parse all float values
	type t_float_fields struct {
		float_epsilon *float64
		float_gt      *float64
		float_lt      *float64
		float_gte     *float64
		float_lte     *float64
	}
	float_fields := &t_float_fields{}
	for agn, agv := range option.AggregatedValues {
		if strings.HasPrefix(agn, "float_") {
			float_value, err := strconv.ParseFloat(agv.Source, 64)
			if err != nil {
				return fmt.Errorf("Invalid '%s' value '%s': %v", agn, agv.Source, err)
			}
			switch agn {
			case "float_gt":
				float_fields.float_gt = &float_value
			case "float_lt":
				float_fields.float_lt = &float_value
			case "float_epsilon":
				float_fields.float_epsilon = &float_value
			case "float_gte":
				float_fields.float_gte = &float_value
			case "float_lte":
				float_fields.float_lte = &float_value
			}
		}
	}

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
				error_msg := fmt.Sprintf(`%s.New("Cannot be blank")`)
				vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if strings.HasPrefix(agn, "float_") {
			supported = true
			// checked at bottom
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	if float_fields.float_gt != nil || float_fields.float_gte != nil || float_fields.float_lt != nil || float_fields.float_lte != nil {
		// float_gt
		if float_fields.float_gt != nil {
			ffvalue := *float_fields.float_gt
			if float_fields.float_epsilon != nil {
				ffvalue -= *float_fields.float_epsilon
			}

			g.P("if ", varSrc, " <= ", fmt.Sprintf("%f", ffvalue), " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must be greater than %f")`, errors_alias, *float_fields.float_gt)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, "float_gt", fproto_gowrap_validator.VEID_MINMAX, "float_gt", fmt.Sprintf("%f", *float_fields.float_gt))
			g.Out()
			g.P("}")
		}
		// float_gte
		if float_fields.float_gte != nil {
			ffvalue := *float_fields.float_gte
			if float_fields.float_epsilon != nil {
				ffvalue -= *float_fields.float_epsilon
			}

			g.P("if ", varSrc, " < ", fmt.Sprintf("%f", ffvalue), " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must be greater or equals to %f")`, errors_alias, *float_fields.float_gte)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, "float_gte", fproto_gowrap_validator.VEID_MINMAX, "float_gte", fmt.Sprintf("%f", *float_fields.float_gte))
			g.Out()
			g.P("}")
		}
	}

	return nil
}

func (t *Validator_Std) generateValidation_scalar_string(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false
		//
		// xrequired
		//
		if agn == "xrequired" {
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if ", varSrc, " == \"\" {")
				g.In()
				error_msg := fmt.Sprintf(`%s.New("Cannot be blank")`, errors_alias)
				vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if agn == "length_eq" {
			supported = true
			g.P("if len(", varSrc, ") != ", option.AggregatedValues[agn].Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Length must be %s")`, errors_alias, option.AggregatedValues[agn].Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_LENGTH)
			g.Out()
			g.P("}")
		} else if agn == "regex" {
			supported = true
			regex_alias := g.DeclDep("regexp", "regexp")

			g.P("if stdrematch, stdreerr := ", regex_alias, ".MatchString(`", option.AggregatedValues[agn].Source, "`, ", varSrc, "); stdreerr != nil || !stdrematch {")
			g.In()
			g.P("if stdreerr != nil {")
			g.In()
			vh.GenerateValidationErrorAdd(g.G(), "stdreerr", agn, fproto_gowrap_validator.VEID_INTERNAL_ERROR)
			g.Out()
			g.P("} else {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Text doesn't match the required pattern")`, errors_alias)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_PATTERN)
			g.Out()
			g.P("}")
			g.Out()
			g.P("}")
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	return nil
}

func (t *Validator_Std) generateValidation_scalar_repeated(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, repeatedType fproto_gowrap_validator.RepeatedType, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	return nil
}

//
// Validator_Std_Repeated
//
type Validator_Std_Repeated struct {
	validatorType *fdep.OptionType
}

func (tp *Validator_Std_Repeated) FPValidator() {

}

func (t *Validator_Std_Repeated) GenerateValidationRepeated(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, repeatedType fproto_gowrap_validator.RepeatedType, tp *fdep.DepType, option *fproto.OptionElement, varSrc string) error {
	if option.ParenthesizedName != "validate.rfield" {
		return nil
	}

	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		supported := false
		if agn == "xrequired" {
			//
			// xrequired
			//
			supported = true
			if option.AggregatedValues[agn].Source == "true" {
				g.P("if ", varSrc, " == nil || len(", varSrc, ") == 0 {")
				g.In()
				error_msg := fmt.Sprintf(`%s.New("Is required")`, errors_alias)
				vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if agn == "length_eq" {
			//
			// length_eq
			//
			supported = true
			g.P("if len(", varSrc, ") != ", option.AggregatedValues[agn].Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must have exactly %s items")`, errors_alias, option.AggregatedValues[agn].Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_LENGTH, "eq", option.AggregatedValues[agn].Source)
			g.Out()
			g.P("}")
			/*
				} else if agn == "int_gt" {
					//
					// int_gt
					//
					supported = true
					g.P("if ", varSrc, " <= ", option.AggregatedValues[agn].Source, " {")
					g.In()
					g.P(varError, " = ", errors_alias, `.New("Must be greater than `, option.AggregatedValues[agn].Source, `")`)
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), varError, agn, fproto_gowrap_validator.VEID_MINMAX, "int_gt", option.AggregatedValues[agn].Source)
				} else if agn == "int_lt" {
					//
					// int_lt
					//
					supported = true
					g.P("if ", varSrc, " >= ", option.AggregatedValues[agn].Source, " {")
					g.In()
					g.P(varError, " = ", errors_alias, `.New("Must be lower than `, option.AggregatedValues[agn].Source, `")`)
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), varError, agn, fproto_gowrap_validator.VEID_MINMAX, "int_lt", option.AggregatedValues[agn].Source)
				} else if agn == "int_gte" {
					//
					// int_gte
					//
					supported = true
					g.P("if ", varSrc, " < ", option.AggregatedValues[agn].Source, " {")
					g.In()
					g.P(varError, " = ", errors_alias, `.New("Must be greater or equals to `, option.AggregatedValues[agn].Source, `")`)
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), varError, agn, fproto_gowrap_validator.VEID_MINMAX, "int_gte", option.AggregatedValues[agn].Source)
				} else if agn == "int_lte" {
					//
					// int_lte
					//
					supported = true
					g.P("if ", varSrc, " > ", option.AggregatedValues[agn].Source, " {")
					g.In()
					g.P(varError, " = ", errors_alias, `.New("Must be lower or equals to `, option.AggregatedValues[agn].Source, `")`)
					g.Out()
					g.P("}")
					vh.GenerateValidationErrorCheck(g.G(), varError, agn, fproto_gowrap_validator.VEID_MINMAX, "int_lte", option.AggregatedValues[agn].Source)
			*/
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for repeated of type %s", agn, tp.FullOriginalName())
		}
	}

	return nil
}
