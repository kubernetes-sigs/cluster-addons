package convert

import (
	"sort"
	"strings"

	v1 "k8s.io/api/rbac/v1"
)

// ruleSet is a wrapper around a set of PolicyRules, and allows for more fluent construction
type ruleSet struct {
	rules []v1.PolicyRule
}

// Add appends rule(s) to the rules in the set.
func (r *ruleSet) Add(rules ...v1.PolicyRule) {
	r.rules = append(r.rules, rules...)
}

// normalize attempts to simplify and combine redundant rules.
func (r *ruleSet) normalize() {
	r.rules = normalizeRules(r.rules)

	r.rules = foldResources(r.rules)

	r.rules = foldVerbs(r.rules)

	sort.Slice(r.rules, func(i, j int) bool { return ruleLT(&r.rules[i], &r.rules[j]) })
}

// ruleLT is a comparison function for PolicyRule, so we can sort into a consistent order.
func ruleLT(l, r *v1.PolicyRule) bool {
	lGroup := firstOrEmpty(l.APIGroups)
	rGroup := firstOrEmpty(r.APIGroups)
	if lGroup != rGroup {
		return lGroup < rGroup
	}
	lResource := firstOrEmpty(l.Resources)
	rResource := firstOrEmpty(r.Resources)
	if lResource != rResource {
		return lResource < rResource
	}
	lVerb := firstOrEmpty(l.Verbs)
	rVerb := firstOrEmpty(r.Verbs)
	if lVerb != rVerb {
		return lVerb < rVerb
	}
	return false
}

// firstOrEmpty is a utility function for our comparison function.
func firstOrEmpty(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

// foldResources combines rules that are the same other than in the resource field, so we can join them and concatenate their resources.
func foldResources(rules []v1.PolicyRule) []v1.PolicyRule {
	var out []v1.PolicyRule

	ruleMap := make(map[string]v1.PolicyRule)

	for _, rule := range rules {
		if len(rule.NonResourceURLs) != 0 {
			out = append(out, rule)
			continue
		}

		key := "groups=" + strings.Join(rule.APIGroups, ",")
		//key += ";resources=" + strings.Join(rule.Resources, ",")
		key += ";resourceNames=" + strings.Join(rule.ResourceNames, ",")
		key += ";verbs=" + strings.Join(rule.Verbs, ",")

		existing, found := ruleMap[key]
		if !found {
			ruleMap[key] = rule
			continue
		}

		existing.Resources = append(existing.Resources, rule.Resources...)
		sort.Strings(existing.Resources)

		ruleMap[key] = existing
	}

	for _, rule := range ruleMap {
		out = append(out, rule)
	}

	return out
}

// foldVerbs combines rules that are the same other than in the verbs field, so we can join them and concatenate their verbs.
func foldVerbs(rules []v1.PolicyRule) []v1.PolicyRule {
	var out []v1.PolicyRule

	ruleMap := make(map[string]v1.PolicyRule)

	for _, rule := range rules {
		if len(rule.NonResourceURLs) != 0 {
			out = append(out, rule)
			continue
		}

		key := "groups=" + strings.Join(rule.APIGroups, ",")
		key += ";resources=" + strings.Join(rule.Resources, ",")
		key += ";resourceNames=" + strings.Join(rule.ResourceNames, ",")
		// key += ";verbs=" + strings.Join(rule.Verbs, ",")

		existing, found := ruleMap[key]
		if !found {
			ruleMap[key] = rule
			continue
		}

		existing.Verbs = append(existing.Verbs, rule.Verbs...)
		sort.Strings(existing.Verbs)

		ruleMap[key] = existing
	}

	for _, rule := range ruleMap {
		out = append(out, rule)
	}

	return out
}

// normalizeRules sorts and deduplicates the values within rules, for easier folding and stable output.
func normalizeRules(rules []v1.PolicyRule) []v1.PolicyRule {
	for i := range rules {
		rule := &rules[i]

		rule.APIGroups = normalizeStringSlice(rule.APIGroups)
		rule.NonResourceURLs = normalizeStringSlice(rule.NonResourceURLs)
		rule.ResourceNames = normalizeStringSlice(rule.ResourceNames)
		rule.Resources = normalizeStringSlice(rule.Resources)
		rule.Verbs = normalizeStringSlice(rule.Verbs)
	}

	return rules
}

// normalizeStringSlice deduplicates and sorts a []string
func normalizeStringSlice(in []string) []string {
	var out []string

	done := make(map[string]bool)
	for _, s := range in {
		if done[s] {
			continue
		}
		out = append(out, s)
	}

	sort.Strings(out)
	return out
}
