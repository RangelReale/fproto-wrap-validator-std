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
	// validator.field
	if validatorType.Option != nil &&
		validatorType.Option.DepFile.FilePath == "github.com/RangelReale/fproto-wrap-validator-std/validator.proto" &&
		validatorType.Option.DepFile.ProtoFile.PackageName == "validator" {
		if validatorType.Name == "field" {
			return &Validator_Std{validatorType: validatorType}
		} else if validatorType.Name == "rfield" {
			return &Validator_Std_Repeated{validatorType: validatorType}
		}
	}
	return nil
}

func (tp *ValidatorPlugin_Std) ValidatorPrefixes() []string {
	return []string{"validator"}
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
	if option.ParenthesizedName != "validator.field" {
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

	int_fields := &rangeValidation{}

	for _, agn := range option.AggregatedSorted() {
		agv := option.AggregatedValues[agn]

		supported := false
		if agn == "required" {
			//
			// required
			//
			supported = true
			if agv.Source == "true" {
				g.P("if ", varSrc, " == 0 {")
				g.In()
				vh.GenerateValidationErrorAdd(g.G(), errors_alias+`.New("Cannot be blank")`, agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if strings.HasPrefix(agn, "int_") {
			supported = true

			// checked at bottom
			_, err := strconv.ParseInt(agv.Source, 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid '%s' value '%s': %v", agn, agv.Source, err)
			}
			switch agn {
			case "int_gt":
				int_fields.setGt(agv.Source)
			case "int_lt":
				int_fields.setLt(agv.Source)
			case "int_gte":
				int_fields.setGte(agv.Source)
			case "int_lte":
				int_fields.setLte(agv.Source)
			default:
				supported = false
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	// Range: int
	err := generateRangeValidation(int_fields, g, vh, tp, option, varSrc, fproto_gowrap_validator.VEID_MINMAX)
	if err != nil {
		return err
	}

	return nil
}

func (t *Validator_Std) generateValidation_scalar_float(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {

	float_fields := &rangeValidation{}

	for _, agn := range option.AggregatedSorted() {
		agv := option.AggregatedValues[agn]

		supported := false
		//
		// required
		//
		if agn == "required" {
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
			_, err := strconv.ParseFloat(agv.Source, 64)
			if err != nil {
				return fmt.Errorf("Invalid '%s' value '%s': %v", agn, agv.Source, err)
			}
			switch agn {
			case "float_gt":
				float_fields.setGt(agv.Source)
			case "float_lt":
				float_fields.setLt(agv.Source)
			case "float_epsilon":
				float_fields.setEpsilon(agv.Source)
			case "float_gte":
				float_fields.setGte(agv.Source)
			case "float_lte":
				float_fields.setLte(agv.Source)
			default:
				supported = false
			}

		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	// Range: float
	err := generateRangeValidation(float_fields, g, vh, tp, option, varSrc, fproto_gowrap_validator.VEID_MINMAX)
	if err != nil {
		return err
	}

	return nil
}

func (t *Validator_Std) generateValidation_scalar_string(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	length_fields := &rangeValidation{}

	for _, agn := range option.AggregatedSorted() {
		agv := option.AggregatedValues[agn]

		supported := false
		//
		// required
		//
		if agn == "required" {
			supported = true
			if agv.Source == "true" {
				g.P("if ", varSrc, " == \"\" {")
				g.In()
				error_msg := fmt.Sprintf(`%s.New("Cannot be blank")`, errors_alias)
				vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if agn == "length_eq" {
			supported = true
			g.P("if len(", varSrc, ") != ", agv.Source, " {")
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Length must be %s")`, errors_alias, agv.Source)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_LENGTH)
			g.Out()
			g.P("}")
		} else if agn == "regex" {
			supported = true
			regex_alias := g.DeclDep("regexp", "regexp")

			g.P("if stdrematch, stdreerr := ", regex_alias, ".MatchString(`", agv.Source, "`, ", varSrc, "); stdreerr != nil || !stdrematch {")
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
		} else if strings.HasPrefix(agn, "length") {
			supported = true

			// checked at bottom
			_, err := strconv.ParseFloat(agv.Source, 64)
			if err != nil {
				return fmt.Errorf("Invalid '%s' value '%s': %v", agn, agv.Source, err)
			}
			switch agn {
			case "length_gt":
				length_fields.setGt(agv.Source)
			case "length_lt":
				length_fields.setLt(agv.Source)
			case "length_gte":
				length_fields.setGte(agv.Source)
			case "length_lte":
				length_fields.setLte(agv.Source)
			case "length_eq":
				length_fields.setEq(agv.Source)
			default:
				supported = false
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.FullOriginalName())
		}
	}

	// Range: length
	err := generateRangeValidation(length_fields, g, vh, tp, option, fmt.Sprintf("len(%s)", varSrc), fproto_gowrap_validator.VEID_LENGTH)
	if err != nil {
		return err
	}

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
	if option.ParenthesizedName != "validator.rfield" {
		return nil
	}

	errors_alias := g.DeclDep("errors", "errors")

	length_fields := &rangeValidation{}

	for _, agn := range option.AggregatedSorted() {
		agv := option.AggregatedValues[agn]

		supported := false
		if agn == "required" {
			//
			// required
			//
			supported = true
			if agv.Source == "true" {
				g.P("if ", varSrc, " == nil || len(", varSrc, ") == 0 {")
				g.In()
				error_msg := fmt.Sprintf(`%s.New("Is required")`, errors_alias)
				vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
				g.Out()
				g.P("}")
			}
		} else if strings.HasPrefix(agn, "length") {
			supported = true

			// checked at bottom
			_, err := strconv.ParseFloat(agv.Source, 64)
			if err != nil {
				return fmt.Errorf("Invalid '%s' value '%s': %v", agn, agv.Source, err)
			}
			switch agn {
			case "length_gt":
				length_fields.setGt(agv.Source)
			case "length_lt":
				length_fields.setLt(agv.Source)
			case "length_gte":
				length_fields.setGte(agv.Source)
			case "length_lte":
				length_fields.setLte(agv.Source)
			case "length_eq":
				length_fields.setEq(agv.Source)
			default:
				supported = false
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for repeated of type %s", agn, tp.FullOriginalName())
		}
	}

	// Range: length
	err := generateRangeValidation(length_fields, g, vh, tp, option, fmt.Sprintf("len(%s)", varSrc), fproto_gowrap_validator.VEID_LENGTH)
	if err != nil {
		return err
	}

	return nil
}
