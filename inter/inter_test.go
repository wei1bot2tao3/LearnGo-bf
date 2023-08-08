package inter

import (
	"fmt"
	"testing"
)

// Test
func Test(t *testing.T) {

	ttgs := &Assets{
		assets: []Asset{
			&GlassDoor{},
			&WoodDoor{},
		},
	}

	ttgs.DoStartWork()
	fmt.Println("work 8h")
	ttgs.DoStopWork()
	fmt.Println("下班")

}
