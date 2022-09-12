package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
)

func handleLabels() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		labelReq := &struct {
			Labels []string `json:"labels"`
		}{Labels: []string{}}

		if err := readJSON(ctx, req, labelReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		matchSet := makeMatchSet(labelReq.Labels)

		labelRes := &struct {
			Labels map[string][]string `json:"labels"`
		}{}

		// Iterate over all resources and cache the labels.
		// We may have to optimize this query by adding a special
		// store that will keep this cache intact as labels are
		// added and removed to the resources. For now it a naive
		// o(n)*o(m) implementation, where n is number of resources
		// and m is amortized number of labels per resource
		rMap := api.Store.Query()
		labelRes.Labels = matchLabels(matchSet, rMap)

		writeJSON(ctx, res, labelRes)
	}
}

func makeMatchSet(labels []string) map[string]struct{} {
	matchSet := make(map[string]struct{}, len(labels))
	for _, l := range labels {
		matchSet[l] = struct{}{}
	}

	return matchSet
}

func matchLabels(matchSet map[string]struct{}, rMap *zebra.ResourceMap) map[string][]string {
	matchAny := len(matchSet) == 0
	labelSet := make(map[string]map[string]struct{})

	for _, resList := range rMap.Resources {
		for _, r := range resList.Resources {
			labels := r.GetMeta().Labels
			for k, v := range labels {
				_, matchOk := matchSet[k] // check if key is list
				matched := matchAny || matchOk

				if matched {
					valueSet, ok := labelSet[k]
					if !ok {
						valueSet = make(map[string]struct{})
						labelSet[k] = valueSet
					}

					// add the value
					valueSet[v] = struct{}{}
				}
			}
		}
	}

	return makeLabelValues(labelSet)
}

func makeLabelValues(labelSet map[string]map[string]struct{}) map[string][]string {
	labelVals := make(map[string][]string)

	for k, valueSet := range labelSet {
		valueList := make([]string, 0, len(valueSet))
		for v := range valueSet {
			valueList = append(valueList, v)
		}

		labelVals[k] = valueList
	}

	return labelVals
}
