package provider

import (
	"fmt"
	"github.com/uberate/i18n/pkg"
	"testing"
)

var BaseI18nValue *I18n

func Init() {
	BaseI18nValue = NewI18n(ISO6391)
	p := BaseI18nValue.Pusher("system", "text", "error")
	p(EnglishLn, "error occur")
	p(ChineseLn, "错误")
	p = BaseI18nValue.Pusher("system", "error", "unknown")
	p(EnglishLn, "Unknown error")
	p(ChineseLn, "未知错误")
	ip := BaseI18nValue.IPusher("user", "text", "test")
	ip(EnglishLn, "test")(ChineseLn, "测试")
}

func TestToJson(t *testing.T) {
	Init()
	value, err := pkg.ToJSON(BaseI18nValue)
	if err != nil {
		t.Error(err)
	}
	res, err := pkg.FromJson(value)
	if err != nil {
		t.Error(err)
	}
	if !res.IsMessageEquals(BaseI18nValue) {
		t.Error("Res should equals BaseI18nValue. But not.")
		return
	}
	res.PushMessage(EnglishLn, "test", "test")
	if res.IsMessageEquals(BaseI18nValue) {
		t.Error("Res should different from BaseI18nValue. But not.")
		fmt.Println(pkg.ToJSON(res))
		fmt.Println(pkg.ToJSON(BaseI18nValue))
	}
}

func TestToStrings(t *testing.T) {
	Init()
	BaseI18nValue.WalkRecord(func(languageValue, messageValue string, flags ...string) {
		if value, ok := BaseI18nValue.MessageByString(languageValue, flags...); !ok || messageValue != value {
			t.Errorf("Get: [%s], want: [%s], [%v]", messageValue, value, ok)
		}
	})
}