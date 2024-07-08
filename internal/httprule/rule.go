package httprule

import (
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Get(m protoreflect.MethodDescriptor) (*annotations.HttpRule, bool) {
	descriptor, ok := proto.GetExtension(m.Options(), annotations.E_Http).(*annotations.HttpRule)
	if !ok || descriptor == nil {
		return nil, false
	}
	return descriptor, true
}

type Rule struct {
	// The HTTP method to use.
	Method string
	// The template describing the URL to use.
	Template        Template
	Body            string
	AdditionalRules []Rule
}

func ParseRule(httpRule *annotations.HttpRule) (Rule, error) {
	method, err := httpRuleMethod(httpRule)
	if err != nil {
		return Rule{}, err
	}
	url, err := httpRuleURL(httpRule)
	if err != nil {
		return Rule{}, err
	}
	template, err := ParseTemplate(url)
	if err != nil {
		return Rule{}, err
	}
	additional := make([]Rule, len(httpRule.GetAdditionalBindings()))
	for i, r := range httpRule.GetAdditionalBindings() {
		a, err := ParseRule(r)
		if err != nil {
			return Rule{}, fmt.Errorf("parse additional binding %d: %w", i, err)
		}
		additional[i] = a
	}
	return Rule{
		Method:          method,
		Template:        template,
		Body:            httpRule.GetBody(),
		AdditionalRules: additional,
	}, nil
}

func httpRuleURL(rule *annotations.HttpRule) (string, error) {
	switch v := rule.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		return v.Get, nil
	case *annotations.HttpRule_Post:
		return v.Post, nil
	case *annotations.HttpRule_Delete:
		return v.Delete, nil
	case *annotations.HttpRule_Patch:
		return v.Patch, nil
	case *annotations.HttpRule_Put:
		return v.Put, nil
	case *annotations.HttpRule_Custom:
		return v.Custom.GetPath(), nil
	default:
		return "", fmt.Errorf("http rule does not have an URL defined")
	}
}

func httpRuleMethod(rule *annotations.HttpRule) (string, error) {
	switch v := rule.GetPattern().(type) {
	case *annotations.HttpRule_Get:
		return http.MethodGet, nil
	case *annotations.HttpRule_Post:
		return http.MethodPost, nil
	case *annotations.HttpRule_Delete:
		return http.MethodDelete, nil
	case *annotations.HttpRule_Patch:
		return http.MethodPatch, nil
	case *annotations.HttpRule_Put:
		return http.MethodPut, nil
	case *annotations.HttpRule_Custom:
		return v.Custom.GetKind(), nil
	default:
		return "", fmt.Errorf("http rule does not have an URL defined")
	}
}
