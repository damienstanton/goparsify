package calc

import (
	"fmt"

	"github.com/damienstanton/goparsify"
)

var (
	value goparsify.Parser

	sumOp  = goparsify.Chars("+-", 1, 1)
	prodOp = goparsify.Chars("/*", 1, 1)

	groupExpr = goparsify.Seq("(", sum, ")").Map(func(n *goparsify.Result) {
		n.Result = n.Child[1].Result
	})

	number = goparsify.NumberLit().Map(func(n *goparsify.Result) {
		switch i := n.Result.(type) {
		case int64:
			n.Result = float64(i)
		case float64:
			n.Result = i
		default:
			panic(fmt.Errorf("unknown value %#v", i))
		}
	})

	sum = goparsify.Seq(prod, goparsify.Some(goparsify.Seq(sumOp, prod))).Map(func(n *goparsify.Result) {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "+":
				i += op.Child[1].Result.(float64)
			case "-":
				i -= op.Child[1].Result.(float64)
			}
		}

		n.Result = i
	})

	prod = goparsify.Seq(&value, goparsify.Some(goparsify.Seq(prodOp, &value))).Map(func(n *goparsify.Result) {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "/":
				i /= op.Child[1].Result.(float64)
			case "*":
				i *= op.Child[1].Result.(float64)
			}
		}

		n.Result = i
	})

	y = goparsify.Maybe(sum)
)

func init() {
	value = goparsify.Any(number, groupExpr)
}

func calc(input string) (float64, error) {
	result, err := goparsify.Run(y, input)
	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}
