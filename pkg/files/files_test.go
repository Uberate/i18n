package files

import (
	"fmt"
	"github.com/uberate/i18n/pkg/provider"
	"testing"
)

var BaseI18nValue *provider.I18n

func Init() {
	BaseI18nValue = provider.NewI18n(provider.ISO6391)
	p := BaseI18nValue.Pusher("system", "text", "error")
	p(provider.EnglishLn, "error occur")
	p(provider.ChineseLn, "错误")
	p = BaseI18nValue.Pusher("system", "error", "unknown")
	p(provider.EnglishLn, "Unknown error")
	p(provider.ChineseLn, "未知错误")
	ip := BaseI18nValue.IPusher("user", "text", "test")
	ip(provider.EnglishLn, "test")(provider.ChineseLn, "测试")
}

func TestToJson(t *testing.T) {
	Init()
	value, err := ToJSON(BaseI18nValue)
	if err != nil {
		t.Error(err)
	}
	res, err := FromJson(value)
	if err != nil {
		t.Error(err)
	}
	if !res.IsMessageEquals(BaseI18nValue) {
		t.Error("Res should equals BaseI18nValue. But not.")
		return
	}
	res.PushMessage(provider.EnglishLn, "test", "test")
	if res.IsMessageEquals(BaseI18nValue) {
		t.Error("Res should different from BaseI18nValue. But not.")
		fmt.Println(ToJSON(res))
		fmt.Println(ToJSON(BaseI18nValue))
	}
}
