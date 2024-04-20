package dbpostgres

import "fmt"

func genParam(i *int) string {
	param := fmt.Sprintf("$%d", *i)
	*i++
	return param
}
