package test

import (
	"k3gin/app/util/regexpr"
	"testing"
)

func TestMatch(t *testing.T) {
	t.Log(regexpr.Match(`(.*)(index\.html|doc\.json)[?|.]*`, "index.html/cash.json/test"))
}

func TestMatchEmail(t *testing.T) {
	t.Log(regexpr.MatchEmail("yelei-name-@4k@3k.com"))
}

func TestMatchPhone(t *testing.T) {
	t.Log(regexpr.MatchPhone("13631375979"))
}

func TestMatchDomain(t *testing.T) {
	t.Log(regexpr.MatchDomain("https://www.cnblogs.com/speeding/p/5097790.html/https://www.cnblogs.com/"))
}

func TestMatchIP(t *testing.T) {
	t.Log(regexpr.MatchIP("192.2.0.112.12", regexpr.IPV4))
	t.Log(regexpr.MatchIP("2001:0db8:863:08d3:1319:8a2e:0370:7344:1xxx", regexpr.IPV6))
}
func TestMatchUserName(t *testing.T) {
	t.Log(regexpr.MatchUserName("sss1212sdf"))
}
