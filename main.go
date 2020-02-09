package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/app/router"
)

type parser func(string) (router.Domain_Type, string)

const Domain_Comment = -1

func v2rayRuleParser(rule string) (router.Domain_Type, string) {
	const defaultType = router.Domain_Domain
	prefixes := map[string]router.Domain_Type{
		"domain": router.Domain_Domain,
		"plain":  router.Domain_Plain,
		"regex":  router.Domain_Regex,
		"regexp": router.Domain_Regex,
		"full":   router.Domain_Full,
	}
	if strings.HasPrefix(rule, "#") {
		return Domain_Comment, ""
	}
	if s := strings.Split(rule, ":"); len(s) > 1 {
		if val, found := prefixes[strings.ToLower(s[0])]; found {
			return val, strings.Join(s[1:], ":")
		} else {
			return defaultType, strings.Join(s, ":")
		}
	} else {
		return defaultType, s[0]
	}
}

func autoProxyRuleParser(rule string) (ruleType router.Domain_Type, domain string) {
	switch {
	case strings.HasPrefix(rule, "[") || strings.HasPrefix(rule, "!") || strings.HasPrefix(rule, "@@"):
		ruleType, domain = Domain_Comment, ""
	case strings.HasPrefix(rule, "||"):
		ruleType, domain = router.Domain_Domain, strings.TrimPrefix(rule, "||")
	case strings.HasPrefix(rule, "|") && strings.HasSuffix(rule, "|"):
		ruleType, domain = router.Domain_Full, strings.Trim(rule, "|")
	case strings.HasPrefix(rule, "|"):
		fmt.Printf("Unsupported rule (start anchor): %s. Regarded as plaintext rule.\n", rule)
		ruleType, domain = router.Domain_Plain, strings.TrimPrefix(rule, "|")
	case strings.HasSuffix(rule, "|"):
		fmt.Printf("Unsupported rule (end anchor): %s. Regarded as plaintext rule.\n", rule)
		ruleType, domain = router.Domain_Plain, strings.TrimSuffix(rule, "|")
	default:
		ruleType, domain = router.Domain_Plain, rule
	}
	if strings.HasSuffix(rule, "/") && strings.HasPrefix(rule, "/") {
		ruleType, domain = router.Domain_Regex, strings.Trim(rule, "/")
	} else if strings.Contains(domain, "*") {
		regex := strings.ReplaceAll(regexp.QuoteMeta(domain), "\\*", ".*")
		if ruleType == router.Domain_Full {
			regex = "^" + regex + "$"
		}
		ruleType, domain = router.Domain_Regex, regex
	}
	return
}

func getSitesList(fileName string, parser parser) (list []*router.Domain) {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	rules := strings.Split(string(d), "\n")

	for _, rule := range rules {
		rule := strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		t, val := parser(rule)
		if t == Domain_Comment {
			continue
		}
		list = append(list, &router.Domain{
			Type:  t,
			Value: val,
		})
	}
	return
}

func main() {
	parsers := map[string]parser{
		"v2ray":     v2rayRuleParser,
		"autoproxy": autoProxyRuleParser,
	}

	var (
		sites  = flag.String("sites", "sites", "Folder storing site files.")
		output = flag.String("output", "geosite.dat", "Path of the output .dat file.")
		format = flag.String("format", "v2ray", "Format of the site files.")
	)
	flag.Parse()

	var parser parser
	var ok bool
	if parser, ok = parsers[strings.ToLower(*format)]; !ok {
		panic(fmt.Sprintf("Unsupported format %s.", *format))
	}

	siteList := new(router.GeoSiteList)

	files, err := ioutil.ReadDir(*sites)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filename := file.Name()
		siteList.Entry = append(siteList.Entry, &router.GeoSite{
			CountryCode: strings.ToUpper(strings.TrimSuffix(filename, filepath.Ext(filename))),
			Domain:      getSitesList(filepath.Join(*sites, filename), parser),
		})
	}

	siteListBytes, err := proto.Marshal(siteList)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(*output, siteListBytes, 0666); err != nil {
		panic(err)
	}

	fmt.Printf("File generated successfully: %s.\n", *output)
}
