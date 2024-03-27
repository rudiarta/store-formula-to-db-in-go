package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// Product Type
var (
	PRODUCTA = `PRODUCTA`
	PRODUCTB = `PRODUCTB`
	PRODUCTC = `PRODUCTC`
	PRODUCTD = `PRODUCTD`
	PRODUCTE = `PRODUCTE`
	PRODUCTF = `PRODUCTF`

	UncommonProduct = map[string]bool{
		PRODUCTA: true,
		PRODUCTB: true,
		PRODUCTF: true,
		PRODUCTC: true,
	}

	// add on function to expression ===================================
	truncateFunction cel.EnvOption = cel.Function("truncate_float", cel.Overload(
		"truncate_float",
		[]*cel.Type{cel.DoubleType, cel.IntType},
		cel.DoubleType,
		cel.BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
			trimmedStringFloat := ""

			switch rhs.Value().(int64) {
			case 2:
				trimmedStringFloat = fmt.Sprintf("%.2f", lhs.Value().(float64))
			case 4:
				trimmedStringFloat = fmt.Sprintf("%.4f", lhs.Value().(float64))
			default:
				trimmedStringFloat = fmt.Sprintf("%.2f", lhs.Value().(float64))
			}

			trimmedFloat, _ := strconv.ParseFloat(trimmedStringFloat, 64)

			return types.Double(trimmedFloat)
		}),
	))

	mathCeil cel.EnvOption = cel.Function("math_ceil", cel.Overload(
		"math_ceil",
		[]*cel.Type{cel.DoubleType},
		cel.DoubleType,
		cel.UnaryBinding(func(lhs ref.Val) ref.Val {
			return types.Double(math.Ceil(lhs.Value().(float64)))
		}),
	))

	mathRound cel.EnvOption = cel.Function("math_round", cel.Overload(
		"math_round",
		[]*cel.Type{cel.DoubleType},
		cel.DoubleType,
		cel.UnaryBinding(func(lhs ref.Val) ref.Val {
			return types.Double(math.Round(lhs.Value().(float64)))
		}),
	))
	// end of add on function to expression ===================================
)

const (
	/* Rule:
	 *	Function: math_round(x) , math_ceil(x) , truncate_float(x,y)
	 *	Macros: UncommonProduct
	 */

	// ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight )
	// truncate_float( xxx , 2)
	// ( UncommonProduct.exists(t, t == ProductType) ? math_ceil( xxx ) : ( xxx < 1.31 ? 1.0 : math_round( xxx + 0.19 ) ) )
	DEFAULT = "( UncommonProduct.exists(t, t == ProductType) ? math_ceil( truncate_float( ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight ) , 2) ) : ( truncate_float( ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight ) , 2) < 1.31 ? 1.0 : math_round( truncate_float( ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight ) , 2) + 0.19 ) ) )"

	// ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight )
	// truncate_float( xxx , 2)
	// xxx < 1.31 ? 1.0 : math_round( xxx + 0.19 )
	RULE_0_31 = " truncate_float( ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight ) , 2 ) < 1.31 ? 1.0 : math_round( truncate_float( ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight ) , 2 ) + 0.19 ) "

	// ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight )
	// truncate_float( xxx , 2)
	// math_ceil( xxx )
	RULE_0_01 = " math_ceil( truncate_float( ( VolumeWeight > GrossWeight ? VolumeWeight : GrossWeight ) , 2) ) "
)

type ChargeableWeightData struct {
	GrossWeight  float64
	VolumeWeight float64
	ProductType  string
}

func getRule(data ChargeableWeightData, formula string) (chargeAbleWeight float64, volumeWeight float64, err error) {
	env, err := cel.NewEnv(
		cel.Variable("GrossWeight", cel.DoubleType),
		cel.Variable("VolumeWeight", cel.DoubleType),
		cel.Variable("ProductType", cel.StringType),
		cel.Variable("UncommonProduct", cel.MapType(cel.StringType, cel.BoolType)),
		truncateFunction,
		mathCeil,
		mathRound,
	)
	if err != nil {
		return 0, 0, err
	}

	ast, iss := env.Compile(formula)
	if iss.Err() != nil {
		return 0, 0, iss.Err()
	}

	ast, iss = env.Check(ast)
	if iss.Err() != nil {
		return 0, 0, iss.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		return 0, 0, err
	}

	out, _, err := prg.Eval(map[string]any{
		"GrossWeight":     data.GrossWeight,
		"VolumeWeight":    data.VolumeWeight,
		"ProductType":     data.ProductType,
		"UncommonProduct": UncommonProduct,
	})
	if err != nil {
		return 0, 0, err
	}

	if _, ok := out.Value().(float64); !ok {
		return 0, 0, errors.New("value is not a float64")
	}

	return out.Value().(float64), data.VolumeWeight, nil
}
