package inter

import "fmt"

type Asset interface {
}

type Assets struct {
	assets []Asset
}

// DoStartWork 打开全部门
func (a *Assets) DoStartWork() {
	fmt.Println("公司准备开门")
	for _, item := range a.assets {
		if d, ok := item.(Door); ok {
			d.Unlock()
			d.Open()
		}
	}
}

// DoStopWork 关闭全部门
func (a *Assets) DoStopWork() {
	fmt.Println("公司准备关门")
	for _, item := range a.assets {
		if d, ok := item.(Door); ok {

			d.Close()
			d.Lock()
		}

	}

}
