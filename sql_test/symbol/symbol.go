//go:generate go install github.com/traefik/yaegi/cmd/yaegi@v0.16.1
//go:generate yaegi extract github.com/go-sql-driver/mysql
package symbol

import "reflect"

var Symbols = map[string]map[string]reflect.Value{}
