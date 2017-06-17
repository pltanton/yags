package suspend

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	startMonitorIfNeed()
	fmt.Println(instance)
	startMonitorIfNeed()
	fmt.Println(instance)
}
