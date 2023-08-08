package inter

import "fmt"

// 要把一个Glassdoor 对象的指针赋值给一个Door类型的变量
var _ Door = &GlassDoor{}
var _ Door = &WoodDoor{}

type Door interface {
	Unlock()
	Open()
	Close()
	Lock()
}

type GlassDoor struct {
}

func (g *GlassDoor) Open() {
	fmt.Println("GlassDor Open")
}

func (g *GlassDoor) Close() {

	fmt.Println("GlassDor close")
}
func (g *GlassDoor) Unlock() {

	fmt.Println("GlassDor  Unlock")
}

func (g *GlassDoor) Lock() {

	fmt.Println("GlassDor  lock")
}

type WoodDoor struct {
}

func (w *WoodDoor) Unlock() {
	fmt.Println("WoodDoor  lock")
}

func (w *WoodDoor) Lock() {
	fmt.Println("WoodDoor  Unlock")
}

func (w *WoodDoor) Open() {
	fmt.Println("WoodDoor Open")
}

func (w *WoodDoor) Close() {
	fmt.Println("WoodDoor Close")
}
