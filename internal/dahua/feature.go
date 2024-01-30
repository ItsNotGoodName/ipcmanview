package dahua

import (
	"cmp"
	"slices"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

func init() {
	featureMap = make(map[string]models.DahuaFeature)
	for _, feature := range FeatureList {
		featureMap[feature.Value] = feature.DahuaFeature
	}
	slices.SortFunc(FeatureList, func(a Feature, b Feature) int { return cmp.Compare(a.Value, b.Value) })
}

var FeatureList []Feature = []Feature{
	{models.DahuaFeatureCamera, "camera", "Camera", ""},
}

type Feature struct {
	models.DahuaFeature
	Value       string
	Name        string
	Description string
}

var featureMap map[string]models.DahuaFeature

func FeatureFromStrings(featureStrings []string) models.DahuaFeature {
	var f models.DahuaFeature
	for _, featureString := range featureStrings {
		feature, ok := featureMap[featureString]
		if !ok {
			continue
		}
		f = f | feature
	}
	return f
}

func FeatureToStrings(feature models.DahuaFeature) []string {
	var strings []string
	for _, v := range FeatureList {
		if v.DahuaFeature.EQ(feature) {
			strings = append(strings, v.Value)
		}
	}
	return strings
}
