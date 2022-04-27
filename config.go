package importas

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type Config struct {
	RequiredAlias        map[string]string
	Rules                []*Rule
	DisallowUnaliased    bool
	DisallowExtraAliases bool
}

func (c *Config) CompileRegexp() error {
	rules := make([]*Rule, 0, len(c.RequiredAlias))
	for path, alias := range c.RequiredAlias {
		reg, err := regexp.Compile(fmt.Sprintf("^%s$", path))
		if err != nil {
			return err
		}

		rules = append(rules, &Rule{
			Regexp: reg,
			Alias:  alias,
		})
	}

	c.Rules = rules
	return nil
}

func (c *Config) findRule(path string) *Rule {
	for _, rule := range c.Rules {
		if rule.Regexp.MatchString(path) {
			return rule
		}
	}

	return nil
}

func (c *Config) AliasFor(path string) (string, bool) {
	rule := c.findRule(path)
	if rule == nil {
		return "", false
	}

	alias, err := rule.aliasFor(path)
	if err != nil {
		return "", false
	}

	return alias, true
}

type Rule struct {
	Alias  string
	Regexp *regexp.Regexp
}

func (r *Rule) aliasFor(path string) (string, error) {
	groups := r.Regexp.FindAllString(path, -1)
	if len(groups) > 0 {

		funcMap := template.FuncMap{
			"Title": func(s string) string {
				return strings.ToTitle(strings.ToLower(s))
			},
		}
		tmpl, err := template.New("test").Funcs(funcMap).Parse(r.Alias)
		if err != nil {
			return "", err
		}

		alias := &strings.Builder{}
		err = tmpl.Execute(alias, groups)

		return alias.String(), err
	}

	return "", errors.New("mismatch rule")
}
