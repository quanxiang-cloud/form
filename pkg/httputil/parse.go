package httputil

import (
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// QueryToBody convert query to body json
func QueryToBody(query url.Values, pretty bool) string {
	val := value{typ: "object"}
	for k, v := range query {
		if len(v) > 0 {
			val.add(k, v[0], 0)
		}
	}

	q := val.value(0)
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(q, "", "  ")
	} else {
		b, err = json.Marshal(q)
	}
	if err != nil {
		return ""
	}

	return string(b)
}

type value struct {
	name string
	val  interface{}
	subs []*value
	typ  string
}

func (v *value) add(k, val string, depth int) {
	name, sub := splitKey(k)
	if !isArrayIndex(name) {
		v.typ = "object"
	}
	if sub == "" {
		v.getOrAddSub(name, "").setValue(name, val)
	} else {
		v.getOrAddSub(name, "array").add(sub, val, depth+1)
	}
}

func (v *value) setValue(name, val string) {
	v.name = name
	switch {
	case val == "true" || val == "false":
		v.typ = "boolean"
		v.val, _ = strconv.ParseBool(val)
	default:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			v.typ = "number"
			v.val = f
		} else {
			v.typ = "string"
			v.val = val
		}
	}
}

func (v *value) getOrAddSub(name string, typ string) *value {
	for _, p := range v.subs {
		if p.name == name {
			return p
		}
	}
	p := &value{
		name: name,
		typ:  typ,
	}
	v.subs = append(v.subs, p)
	return p
}

func (v *value) value(depth int) interface{} {
	switch v.typ {
	case "string", "number", "boolean":
		return v.val
	case "array":
		sort.Slice(v.subs, func(i, j int) bool {
			return v.subs[i].name < v.subs[j].name
		})
		r := []interface{}{}
		for _, v := range v.subs {
			r = append(r, v.value(depth+1))
		}
		return r
	case "object":
		r := map[string]interface{}{}
		for _, v := range v.subs {
			r[v.name] = v.value(depth + 1)
		}
		return r
	default:
	}
	return nil
}

func splitKey(key string) (parent, sub string) {
	parent = key
	if index := strings.Index(key, "."); index > 0 {
		parent = key[:index]
		sub = key[index+1:]
	}
	return
}

func isArrayIndex(key string) bool {
	if key == "" {
		return false
	}
	for _, v := range key {
		if !(v >= '0' && v <= '9') {
			return false
		}
	}
	return true
}
