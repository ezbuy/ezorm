package people

import (
	"bytes"
	"fmt"
	"github.com/ezbuy/ezorm/db"
)

var (
	_ db.M
	_ = fmt.Println
	_ bytes.Buffer
)
