package mongokit

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/256dpi/lungo/bsonkit"
)

type Operator func(*Context, bsonkit.Doc, string, interface{}) error

type Context struct {
	Upsert     bool
	TopLevel   map[string]Operator
	Expression map[string]Operator
}

func Process(ctx *Context, doc bsonkit.Doc, query bson.D, root bool) error {
	// process all expressions (implicit and)
	for _, exp := range query {
		err := ProcessExpression(ctx, doc, "", exp, root)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessExpression(ctx *Context, doc bsonkit.Doc, prefix string, pair bson.E, root bool) error {
	// force top level operators if there are no expression operators
	force := len(ctx.Expression) == 0

	// check for top level operators which may appear together with field
	// expressions in the document
	if force || (len(pair.Key) > 0 && pair.Key[0] == '$') {
		// lookup top level operator
		var operator Operator
		if root {
			operator = ctx.TopLevel[pair.Key]
			if operator == nil {
				return fmt.Errorf("unknown top level operator %q", pair.Key)
			}
		} else {
			operator = ctx.Expression[pair.Key]
			if operator == nil {
				return fmt.Errorf("unknown expression operator %q", pair.Key)
			}
		}

		// call operator
		return operator(ctx, doc, prefix, pair.Value)
	}

	// get path
	path := pair.Key
	if prefix != "" {
		path = prefix + "." + path
	}

	// check for field expressions with a document which may contain either
	// only expression operators or only simple conditions
	if exps, ok := pair.Value.(bson.D); ok {
		// process all expressions (implicit and)
		for i, exp := range exps {
			// break and leave document as a simple condition if the
			// first key does not look like an operator
			if i == 0 && (len(exp.Key) == 0 || exp.Key[0] != '$') {
				break
			}

			// check operator validity
			if len(exp.Key) == 0 || exp.Key[0] != '$' {
				return fmt.Errorf("expected operator, got %q", exp.Key)
			}

			// lookup operator
			operator := ctx.Expression[exp.Key]
			if operator == nil {
				return fmt.Errorf("unknown expression operator %q", exp.Key)
			}

			// call operator
			err := operator(ctx, doc, path, exp.Value)
			if err != nil {
				return err
			}

			// return success if last one
			if i == len(exps)-1 {
				return nil
			}
		}
	}

	// handle pair as a simple condition

	// get the default operator
	operator := ctx.Expression[""]
	if operator == nil {
		return fmt.Errorf("missing default operator")
	}

	// call operator
	err := operator(ctx, doc, path, pair.Value)
	if err != nil {
		return err
	}

	return nil
}