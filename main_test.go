package main

import (
	"testing"
	"v2ray.com/core/app/router"
)

func testParser(t *testing.T, parser parser, cases map[string]router.Domain) {
	t.Run("Cases", func(t *testing.T) {
		for rule, expect := range cases {
			t.Run(rule, func(t *testing.T) {
				ruleType, domain := parser(rule)
				if ruleType != expect.Type || domain != expect.Value {
					t.Errorf("%s: expect %s:%s, got %s:%s",
						rule, expect.Type.String(), expect.Value, ruleType.String(), domain)
				}
			})
		}
	})
}

func TestV2rayRuleParser(t *testing.T) {
	cases := map[string]router.Domain{
		"a.com":        {Type: router.Domain_Domain, Value: "a.com"},
		"a.com:a.com":  {Type: router.Domain_Domain, Value: "a.com:a.com"},
		"domain:a.com": {Type: router.Domain_Domain, Value: "a.com"},
		"plain:a.com":  {Type: router.Domain_Plain, Value: "a.com"},
		"regex:.*":     {Type: router.Domain_Regex, Value: ".*"},
		"full:a.com":   {Type: router.Domain_Full, Value: "a.com"},
		":a.com":       {Type: router.Domain_Domain, Value: ":a.com"},
		"#comment":     {Type: Domain_Comment, Value: ""},
		"doMAin:a.com": {Type: router.Domain_Domain, Value: "a.com"},
		"pLAin:a.com":  {Type: router.Domain_Plain, Value: "a.com"},
		"reGex:.*":     {Type: router.Domain_Regex, Value: ".*"},
		"fuLl:a.com":   {Type: router.Domain_Full, Value: "a.com"},
	}
	testParser(t, v2rayRuleParser, cases)
}

func TestAutoProxyRuleParser(t *testing.T) {
	cases := map[string]router.Domain{
		"a.com":       {Type: router.Domain_Plain, Value: "a.com"},
		"||a.com":     {Type: router.Domain_Domain, Value: "a.com"},
		"|a.com":      {Type: router.Domain_Plain, Value: "a.com"},
		"a.com|":      {Type: router.Domain_Plain, Value: "a.com"},
		"|a.com|":     {Type: router.Domain_Full, Value: "a.com"},
		"!comment":    {Type: Domain_Comment, Value: ""},
		"@@b.com":     {Type: Domain_Comment, Value: ""},
		"[autoproxy]": {Type: Domain_Comment, Value: ""},
		"/.*/":        {Type: router.Domain_Regex, Value: ".*"},
		"a*.com":      {Type: router.Domain_Regex, Value: "a.*\\.com"},
		"||a*.com":    {Type: router.Domain_Regex, Value: "a.*\\.com"},
		"|a*.com":     {Type: router.Domain_Regex, Value: "a.*\\.com"},
		"a*.com|":     {Type: router.Domain_Regex, Value: "a.*\\.com"},
		"|a*.com|":    {Type: router.Domain_Regex, Value: "^a.*\\.com$"},
		"*.a*.com":    {Type: router.Domain_Regex, Value: ".*\\.a.*\\.com"},
	}
	testParser(t, autoProxyRuleParser, cases)
}
