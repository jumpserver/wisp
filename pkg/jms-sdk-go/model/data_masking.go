package model

import "sort"

var _ sort.Interface = DataMaskingRules{}

type DataMaskingRules []DataMaskingRule

func (f DataMaskingRules) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f DataMaskingRules) Len() int {
	return len(f)
}

func (f DataMaskingRules) Less(i, j int) bool {
	return f[i].Priority < f[j].Priority
}

type DataMaskingRule struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	FieldsPattern string `json:"fields_pattern"`
	MaskingMethod string `json:"masking_method"`
	MaskPattern   string `json:"mask_pattern"`
	Priority      int    `json:"priority"`
	IsActive      bool   `json:"is_active"`
}
