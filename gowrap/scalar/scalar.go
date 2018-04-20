package fproto_gowrap_validator_std_scalar

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

func GenerateValidation_scalar_int(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, ignoreRequired bool) error {
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
			if !ignoreRequired {
				if agv.Source == "true" {
					g.P("if ", varSrc, " == 0 {")
					g.In()
					vh.GenerateValidationErrorAdd(g.G(), errors_alias+`.New("Cannot be blank")`, agn, fproto_gowrap_validator.VEID_REQUIRED)
					g.Out()
					g.P("}")
				}
			}
		} else if agn == "int_enum_check" {
			if tienum, tienumok := tp.Item.(*fproto.EnumElement); tienumok {
				supported = true

				var check_list []string
				for _, tv := range tienum.EnumConstants {
					check_list = append(check_list, fmt.Sprintf("int(%s) != %d", varSrc, tv.Tag))
				}

				if len(check_list) > 0 {
					g.P("if ", strings.Join(check_list, " && "), " {")
					g.In()
					vh.GenerateValidationErrorAdd(g.G(), errors_alias+`.New("Must be one of the declared values.")`, agn, fproto_gowrap_validator.VEID_INVALID_VALUE)
					g.Out()
					g.P("}")
				} else {
					vh.GenerateValidationErrorAdd(g.G(), errors_alias+`.New("Enum field type has no declared values.")`, agn, fproto_gowrap_validator.VEID_INTERNAL_ERROR)
				}
			}

		} else if agn == "int_in" || agn == "int_not_in" {
			supported = true

			inVal := &inValidation{
				is_not:    agn == "int_not_in",
				is_string: false,
				list:      make([]string, 0),
			}

			// build a list of values
			if len(agv.Array) == 0 {
				inVal.list = append(inVal.list, agv.Source)
			} else {
				for _, ai := range agv.Array {
					if len(ai.Array) > 0 {
						return fmt.Errorf("Nested arrays are not supported for %s", agn)
					}

					inVal.list = append(inVal.list, ai.Source)
				}
			}

			if inVal.isEmpty() {
				return fmt.Errorf("At least one value is required for %s", agn)
			}

			if err := generateInValidation(inVal, g, vh, tp, option, varSrc); err != nil {
				return err
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
			case "int_eq":
				int_fields.setEq(agv.Source)
			default:
				supported = false
			}
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.TypeDescription())
		}
	}

	// Range: int
	err := generateRangeValidation(int_fields, g, vh, tp, option, varSrc, fproto_gowrap_validator.VEID_MINMAX)
	if err != nil {
		return err
	}

	return nil
}

func GenerateValidation_scalar_float(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, ignoreRequired bool) error {
	errors_alias := g.DeclDep("errors", "errors")

	float_fields := &rangeValidation{}

	for _, agn := range option.AggregatedSorted() {
		agv := option.AggregatedValues[agn]

		supported := false
		//
		// required
		//
		if agn == "required" {
			supported = true
			if !ignoreRequired {
				if option.AggregatedValues[agn].Source == "true" {
					g.P("if ", varSrc, " == 0 {")
					g.In()
					error_msg := fmt.Sprintf(`%s.New("Cannot be blank")`, errors_alias)
					vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
					g.Out()
					g.P("}")
				}
			}
		} else if agn == "float_in" || agn == "float_not_in" {
			supported = true

			inVal := &inValidation{
				is_not:    agn == "float_not_in",
				is_string: false,
				list:      make([]string, 0),
			}

			if aep, aisep := option.AggregatedValues["float_epsilon"]; aisep {
				inVal.epsilon = &aep.Source
			}

			// build a list of values
			if len(agv.Array) == 0 {
				inVal.list = append(inVal.list, agv.Source)
			} else {
				for _, ai := range agv.Array {
					if len(ai.Array) > 0 {
						return fmt.Errorf("Nested arrays are not supported for %s", agn)
					}

					inVal.list = append(inVal.list, ai.Source)
				}
			}

			if inVal.isEmpty() {
				return fmt.Errorf("At least one value is required for %s", agn)
			}

			if err := generateInValidation(inVal, g, vh, tp, option, varSrc); err != nil {
				return err
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
			case "float_eq":
				float_fields.setEq(agv.Source)
			default:
				supported = false
			}

		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.TypeDescription())
		}
	}

	// Range: float
	err := generateRangeValidation(float_fields, g, vh, tp, option, varSrc, fproto_gowrap_validator.VEID_MINMAX)
	if err != nil {
		return err
	}

	return nil
}

func GenerateValidation_scalar_string(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, ignoreRequired bool) error {
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
			if !ignoreRequired {
				if agv.Source == "true" {
					g.P("if ", varSrc, ` == "" {`)
					g.In()
					error_msg := fmt.Sprintf(`%s.New("Cannot be blank")`, errors_alias)
					vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
					g.Out()
					g.P("}")
				}
			}
		} else if agn == "string_eq" {
			g.P("if ", varSrc, ` != `, strconv.Quote(agv.Source), ` {`)
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must have the value'`, agv.Source, `'")`, errors_alias)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
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
		} else if agn == "string_in" || agn == "string_not_in" {
			supported = true

			inVal := &inValidation{
				is_not:    agn == "string_not_in",
				is_string: true,
				list:      make([]string, 0),
			}

			// build a list of values
			if len(agv.Array) == 0 {
				inVal.list = append(inVal.list, agv.Source)
			} else {
				for _, ai := range agv.Array {
					if len(ai.Array) > 0 {
						return fmt.Errorf("Nested arrays are not supported for %s", agn)
					}

					inVal.list = append(inVal.list, ai.Source)
				}
			}

			if inVal.isEmpty() {
				return fmt.Errorf("At least one value is required for %s", agn)
			}

			if err := generateInValidation(inVal, g, vh, tp, option, varSrc); err != nil {
				return err
			}

		} else if strings.HasPrefix(agn, "length_") {
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
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.TypeDescription())
		}
	}

	// Range: length
	if !length_fields.isEmpty() {
		// check if blank first
		g.P("if ", varSrc, ` != "" {`)
		g.In()

		err := generateRangeValidation(length_fields, g, vh, tp, option, fmt.Sprintf("len(%s)", varSrc), fproto_gowrap_validator.VEID_LENGTH)
		if err != nil {
			return err
		}

		g.Out()
		g.P("}")
	}

	return nil
}

func GenerateValidation_scalar_byte(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, ignoreRequired bool) error {
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
			if !ignoreRequired {
				if agv.Source == "true" {
					g.P("if ", varSrc, " == nil || len(", varSrc, ") == 0 {")
					g.In()
					error_msg := fmt.Sprintf(`%s.New("Cannot be blank")`, errors_alias)
					vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
					g.Out()
					g.P("}")
				}
			}
		} else if strings.HasPrefix(agn, "length_") {
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
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.TypeDescription())
		}
	}

	// Range: length
	err := generateRangeValidation(length_fields, g, vh, tp, option, fmt.Sprintf("len(%s)", varSrc), fproto_gowrap_validator.VEID_LENGTH)
	if err != nil {
		return err
	}

	return nil
}

func GenerateValidation_scalar_bool(g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, option *fproto.OptionElement, varSrc string, ignoreRequired bool) error {
	errors_alias := g.DeclDep("errors", "errors")

	for _, agn := range option.AggregatedSorted() {
		agv := option.AggregatedValues[agn]

		supported := false
		//
		// required
		//
		if agn == "bool_eq" {
			g.P("if ", varSrc, ` != `, agv.Source, ` {`)
			g.In()
			error_msg := fmt.Sprintf(`%s.New("Must have the value'`, agv.Source, `'")`, errors_alias)
			vh.GenerateValidationErrorAdd(g.G(), error_msg, agn, fproto_gowrap_validator.VEID_REQUIRED)
			g.Out()
			g.P("}")
		}

		if !supported {
			return fmt.Errorf("Validation %s not supported for type %s", agn, tp.TypeDescription())
		}
	}

	return nil
}
