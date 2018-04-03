package fproto_gowrap_validator_std

import (
	"fmt"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-wrap-validator/gowrap"
	"github.com/RangelReale/fproto-wrap/gowrap"
)

type rangeValidation struct {
	gt      *string
	gte     *string
	lt      *string
	lte     *string
	epsilon *string
	eq      *string
}

func generateRangeValidation(ranges *rangeValidation, g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType, tinfo fproto_gowrap.TypeInfo, option *fproto.OptionElement, varSrc string) error {
	errors_alias := g.DeclDep("errors", "errors")

	gtadd := ""
	ltadd := ""
	if ranges.epsilon != nil {
		gtadd = " - " + *ranges.epsilon
		ltadd = " + " + *ranges.epsilon
	}

	if ranges.gt != nil {
		g.P("if ", varSrc, " <= ", *ranges.gt, gtadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("Must be greater than %s")`, errors_alias, *ranges.gt)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, "gt", fproto_gowrap_validator.VEID_MINMAX, "gt", *ranges.gt)
		g.Out()
		g.P("}")
	}

	if ranges.gte != nil {
		g.P("if ", varSrc, " < ", *ranges.gte, gtadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("Must be greater or equals to %s")`, errors_alias, *ranges.gte)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, "gte", fproto_gowrap_validator.VEID_MINMAX, "gte", *ranges.gte)
		g.Out()
		g.P("}")
	}

	if ranges.lt != nil {
		g.P("if ", varSrc, " > ", *ranges.lt, ltadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("Must be lower than %s")`, errors_alias, *ranges.lt)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, "lt", fproto_gowrap_validator.VEID_MINMAX, "lt", *ranges.lt)
		g.Out()
		g.P("}")
	}

	if ranges.lte != nil {
		g.P("if ", varSrc, " <= ", *ranges.lte, ltadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("Must be lower or equals to %s")`, errors_alias, *ranges.lte)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, "lte", fproto_gowrap_validator.VEID_MINMAX, "lte", *ranges.lte)
		g.Out()
		g.P("}")
	}

	if ranges.eq != nil {
		g.P("if ", varSrc, " != ", *ranges.eq, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("Must be exactly equals to %s")`, errors_alias, *ranges.eq)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, "eq", fproto_gowrap_validator.VEID_MINMAX, "eq", *ranges.eq)
		g.Out()
		g.P("}")
	}

	return nil
}
