package opt

import (
	"sexp"
	"sexpconv"
)

// ReduceStrength replaces operations with their less expensive
// equivalents.
func ReduceStrength(form sexp.Form) sexp.Form {
	switch form := form.(type) {
	case *sexp.NumAdd:
		return weakenAdd(form)
	case *sexp.NumSub:
		return weakenSub(form)

	case *sexp.Bind:
		form.Init = ReduceStrength(form.Init)
	case *sexp.Rebind:
		form.Expr = ReduceStrength(form.Expr)
	case *sexp.Block:
		form.Forms = reduceStrength(form.Forms)
	case *sexp.FormList:
		form.Forms = reduceStrength(form.Forms)

	case *sexp.ArrayLit:
		return weakenArrayLit(form)
	case *sexp.SparseArrayLit:
		return weakenSparseArrayLit(form)
	}

	return form
}

func reduceStrength(forms []sexp.Form) []sexp.Form {
	for i, form := range forms {
		forms[i] = ReduceStrength(form)
	}
	return forms
}

func weakenAdd(form *sexp.NumAdd) sexp.Form {
	weaken := func(a, b int) sexp.Form {
		if numEq(form.Args[a], 1) {
			return addX(form.Args[b], 1, form.Typ)
		}
		if numEq(form.Args[a], 2) {
			return addX(form.Args[b], 2, form.Typ)
		}
		// Addition of negative number = substraction.
		if numEq(form.Args[a], -1) {
			return subX(form.Args[b], 1, form.Typ)
		}
		if numEq(form.Args[a], -2) {
			return subX(form.Args[b], 2, form.Typ)
		}
		return nil
	}

	if form := weaken(0, 1); form != nil {
		return form
	}
	// Because "+" is commutative, we can try to apply
	// same patterns against other argument.
	if form := weaken(1, 0); form != nil {
		return form
	}
	return form
}

func weakenSub(form *sexp.NumSub) sexp.Form {
	if numEq(form.Args[1], 1) {
		return subX(form.Args[0], 1, form.Typ)
	}
	if numEq(form.Args[1], 2) {
		return subX(form.Args[0], 2, form.Typ)
	}
	// Substraction of negative number = addition.
	if numEq(form.Args[1], -1) {
		return addX(form.Args[0], 1, form.Typ)
	}
	if numEq(form.Args[1], -2) {
		return addX(form.Args[0], 2, form.Typ)
	}
	return form
}

func weakenArrayLit(form *sexp.ArrayLit) sexp.Form {
	// #TODO: recognize array where all elements are the same.
	//        Replace with "make-vector" call.
	return form
}

func weakenSparseArrayLit(form *sexp.SparseArrayLit) sexp.Form {
	// Sparse arrays worth it only when zero values are prevalent.
	zeroVals := form.Typ.Len() - int64(len(form.Vals))
	if zeroVals*4 < form.Typ.Len()*3 {
		// Count of zero values < 75%.
		// Convert to ArrayLit.
		zv := sexpconv.ZeroValue(form.Typ.Elem())
		vals := make([]sexp.Form, int(form.Typ.Len()))
		for _, val := range form.Vals {
			vals[val.Index] = val.Expr
		}
		for i := range vals {
			if vals[i] == nil {
				vals[i] = zv
			}
		}
		return &sexp.ArrayLit{Vals: vals}
	}

	return form
}
