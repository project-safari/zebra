package zebra

import (
	"reflect"
	"strings"
)

// Filter given map by uuids.
func FilterUUID(uuids []string, resMap *ResourceMap) (*ResourceMap, error) {
	retMap := NewResourceMap(resMap.GetFactory())

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			if IsIn(res.GetID(), uuids) {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Filter given map by types.
func FilterType(types []string, resMap *ResourceMap) (*ResourceMap, error) {
	f := resMap.GetFactory()
	retMap := NewResourceMap(f)

	for _, t := range types {
		l, ok := resMap.Resources[t]
		if !ok {
			continue
		}

		copyL := NewResourceList(f)

		CopyResourceList(copyL, l)
		retMap.Resources[t] = copyL
	}

	return retMap, nil
}

// Filter given map by label name and val.
func FilterLabel(query Query, resMap *ResourceMap) (*ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return resMap, err
	}

	retMap := NewResourceMap(resMap.GetFactory())

	inVals := false

	if query.Op == MatchEqual || query.Op == MatchIn {
		inVals = true
	}

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			labels := res.GetLabels()
			matchIn := labels.MatchIn(query.Key, query.Values...)

			if (inVals && matchIn) || (!inVals && !matchIn) {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Filter given map by property name (case insensitive) and val.
func FilterProperty(query Query, resMap *ResourceMap) (*ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return resMap, err
	}

	retMap := NewResourceMap(resMap.GetFactory())

	inVals := false

	if query.Op == MatchEqual || query.Op == MatchIn {
		inVals = true
	}

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			val := FieldByName(reflect.ValueOf(res).Elem(), query.Key).String()
			matchIn := IsIn(val, query.Values)

			if (inVals && matchIn) || (!inVals && !matchIn) {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Ignore case in returning value of given field.
func FieldByName(v reflect.Value, field string) reflect.Value {
	field = strings.ToLower(field)

	return v.FieldByNameFunc(
		func(found string) bool {
			return strings.ToLower(found) == field
		})
}
