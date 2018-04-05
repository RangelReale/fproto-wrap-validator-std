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

func (r *rangeValidation) setGt(v string) {
	r.gt = &v
}

func (r *rangeValidation) setGte(v string) {
	r.gte = &v
}

func (r *rangeValidation) setLt(v string) {
	r.lt = &v
}

func (r *rangeValidation) setLte(v string) {
	r.lte = &v
}

func (r *rangeValidation) setEpsilon(v string) {
	r.epsilon = &v
}

func (r *rangeValidation) setEq(v string) {
	r.eq = &v
}

func generateRangeValidation(ranges *rangeValidation, g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType,
	option *fproto.OptionElement, varSrc string, errorId fproto_gowrap_validator.ValidationErrorId) error {
	errors_alias := g.DeclDep("errors", "errors")

	validationItemPrefix := ""
	validationDescription := "Value"

	if errorId == fproto_gowrap_validator.VEID_LENGTH {
		validationItemPrefix = "length_"
		validationDescription = "Length"
	}

	gtadd := ""
	ltadd := ""
	if ranges.epsilon != nil {
		gtadd = " - " + *ranges.epsilon
		ltadd = " + " + *ranges.epsilon
	}

	if ranges.gt != nil {
		g.P("if ", varSrc, " <= ", *ranges.gt, gtadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("%s must be greater than %s")`, errors_alias, validationDescription, *ranges.gt)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, validationItemPrefix+"gt", errorId, "gt", *ranges.gt)
		g.Out()
		g.P("}")
	}

	if ranges.gte != nil {
		g.P("if ", varSrc, " < ", *ranges.gte, gtadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("%s must be greater or equals to %s")`, errors_alias, validationDescription, *ranges.gte)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, validationItemPrefix+"gte", errorId, "gte", *ranges.gte)
		g.Out()
		g.P("}")
	}

	if ranges.lt != nil {
		g.P("if ", varSrc, " >= ", *ranges.lt, ltadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("%s must be lower than %s")`, errors_alias, validationDescription, *ranges.lt)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, validationItemPrefix+"lt", errorId, "lt", *ranges.lt)
		g.Out()
		g.P("}")
	}

	if ranges.lte != nil {
		g.P("if ", varSrc, " > ", *ranges.lte, ltadd, " {")
		g.In()
		error_msg := fmt.Sprintf(`%s.New("%s must be lower or equals to %s")`, errors_alias, validationDescription, *ranges.lte)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, validationItemPrefix+"lte", errorId, "lte", *ranges.lte)
		g.Out()
		g.P("}")
	}

	if ranges.eq != nil {
		if ranges.epsilon != nil {
			g.P("if ", varSrc, " >= ", *ranges.eq, "-", *ranges.epsilon, " && ", varSrc, " <= ", *ranges.eq, "+", *ranges.epsilon, " {")
		} else {
			g.P("if ", varSrc, " != ", *ranges.eq, " {")
		}

		g.In()
		error_msg := fmt.Sprintf(`%s.New("%s must be exactly equals to %s")`, errors_alias, validationDescription, *ranges.eq)
		vh.GenerateValidationErrorAdd(g.G(), error_msg, validationItemPrefix+"eq", errorId, "eq", *ranges.eq)
		g.Out()
		g.P("}")
	}

	return nil
}
