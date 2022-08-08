package definer

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/zhiting-tech/smartassistant/pkg/thingmodel"
)

func NewAttribute(attribute thingmodel.Attribute) *Attribute {
	return &Attribute{
		meta: &attribute,
	}
}

type Attribute struct {
	meta       *thingmodel.Attribute
	iAttribute thingmodel.IAttribute
}

func (a *Attribute) SetVal(val interface{}) {
	if err := checkAttrVal(a, val); err != nil {
		logrus.Debugf("%s set val err: %s", a.Type(), err)
		return
	}
	a.meta.Val = val
}

func (a *Attribute) GetVal() (val interface{}) {
	return a.meta.Val
}

func (a *Attribute) Set(val interface{}) error {
	if a.iAttribute == nil {
		return NotEnableErr
	}

	if err := checkAttrVal(a, val); err != nil {
		return err
	}
	return a.iAttribute.Set(val)
}

func checkAttrVal(a *Attribute, val interface{}) error {

	switch a.meta.ValType {
	case thingmodel.String:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("invalid val type of %s", a.Type())
		}
	case thingmodel.Int, thingmodel.Int32, thingmodel.Int64, thingmodel.Float32, thingmodel.Float64:
		switch val.(type) {
		case float64, int, int32, int64, float32:
			return nil
		default:
			return fmt.Errorf("invalid val type of %s", a.Type())
		}
	case thingmodel.Bool:
		switch val.(type) {
		case float64, int, bool:
			return nil
		default:
			return fmt.Errorf("invalid val type of %s", a.Type())
		}
	case thingmodel.JSON:
		if _, ok := val.(string); !ok {
			return fmt.Errorf("invalid val type of %s", a.Type())
		}
	}
	return nil
}

// Enable 启用属性并通过实现接口设置方法
func (a *Attribute) Enable(attr thingmodel.IAttribute) *Attribute {
	a.iAttribute = attr
	return a
}

// SetRange 设置最小值最大值
func (a *Attribute) SetRange(min, max interface{}) *Attribute {
	a.meta.Min = min
	a.meta.Max = max
	return a
}

// Type 属性类型
func (a *Attribute) Type() thingmodel.Attribute {
	return *a.meta
}
