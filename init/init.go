package init

import (
	"fmt"
	"log"
)

func init() {
	fmt.Println("init is call")
	log.SetFlags(log.Ltime | log.Lshortfile)
}
