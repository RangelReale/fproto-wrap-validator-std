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

//
// range validation
//
type rangeValidation struct {
	gt      *string
	gte     *string
	lt      *string
	lte     *string
	epsilon *string
	eq      *string
}

func (r *rangeValidation) isEmpty() bool {
	return r.gt == nil && r.gte == nil && r.lt == nil && r.lte == nil && r.eq == nil
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

//
// in validation
//
type inValidation struct {
	list      []string
	epsilon   *string
	is_not    bool
	is_string bool
}

func (r *inValidation) isEmpty() bool {
	return len(r.list) == 0
}

func generateInValidation(in *inValidation, g *fproto_gowrap.GeneratorFile, vh fproto_gowrap_validator.ValidatorHelper, tp *fdep.DepType,
	option *fproto.OptionElement, varSrc string) error {

	// TODO: support epsilon
	errors_alias := g.DeclDep("errors", "errors")

	validationItem := "in"
	compare := "!="
	ccond := " && "
	if in.is_not {
		validationItem = "not_in"
		compare = "=="
		ccond = " || "
	}

	compare_list := make([]string, 0)
	for _, l := range in.list {
		cvalue := l
		if in.is_string {
			cvalue = strconv.Quote(l)
		}
		citem := fmt.Sprintf("%s %s %s", varSrc, compare, cvalue)
		compare_list = append(compare_list, citem)
	}

	g.P("if ", strings.Join(compare_list, ccond), " {")
	g.In()
	var error_msg string
	if in.is_not {
		error_msg = fmt.Sprintf(`%s.New("Must not be one of the invalid values")`, errors_alias)
	} else {
		error_msg = fmt.Sprintf(`%s.New("Must be one of the valid values")`, errors_alias)
	}
	vh.GenerateValidationErrorAdd(g.G(), error_msg, validationItem, fproto_gowrap_validator.VEID_INVALID_VALUE)
	g.Out()
	g.P("}")

	return nil
}
