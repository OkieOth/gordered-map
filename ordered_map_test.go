package omap_test

import (
	"maps"
	"os"
	"strings"
	"testing"
	"text/template"

	omap "github.com/okieoth/gordered-map"
	"github.com/stretchr/testify/require"
)

func TestNewFromJSONFile(t *testing.T) {
	fileName := "./resources/test_schema.json"
	mapThing, err := omap.NewFromJSONFile(fileName)
	require.Nil(t, err, "error while create object from file")
	require.NotNil(t, mapThing, "created instance is nil")
	title, found := omap.GetValue[string](mapThing, "title")
	require.True(t, found, "couldn't find title entry")
	require.NotNil(t, title, "title object is nil")
	require.Equal(t, "Person", title, "title has wrong value")

	title2, found := omap.GetValue[string](mapThing, "title2")
	require.False(t, found, "found title2 entry")
	require.Equal(t, "", title2, "although title2 isn't exsiting, was the result not nil")
	found = omap.HasValue(mapThing, "title2")
	require.False(t, found, "found wrong title2 entry")

	// create a new entry
	title2Val := "I am a new entry"
	err = omap.Set(mapThing, "title2", title2Val)
	require.Nil(t, err, "error while setting title2")
	found = omap.HasValue(mapThing, "title2")
	require.True(t, found, "couldn't find title2 entry (1)")

	title2, found = omap.GetValue[string](mapThing, "title2")
	require.True(t, found, "couldn't find title2 entry (2)")
	require.Equal(t, title2Val, title2, "title2 has wrong value")

	// Check that title wasn't overwritten
	title, found = omap.GetValue[string](mapThing, "title")
	require.True(t, found, "couldn't find title entry")
	require.NotNil(t, title, "title object is nil")
	require.Equal(t, "Person", title, "title has wrong value")

	// create a new entry
	omap.Set(mapThing, "title2", 12)
	err = omap.Set(mapThing, "title2", 23)
	require.NotNil(t, err, "error wasn't returned when updating value with wrong type")

	// return map sub type
	propMap, found := omap.GetChildMap(mapThing, "properties")
	require.True(t, found, "couldn't find properties map")
	require.NotNil(t, propMap, "propMap is nil")

	// get value from child map
	nameMap, found := omap.GetChildMap(propMap, "name")
	require.True(t, found, "couldn't find properties map")
	require.NotNil(t, nameMap, "propMap is nil")

	typeStr, found := omap.GetValue[string](nameMap, "type")
	require.True(t, found, "can't find type attrib in 'name' map")
	require.Equal(t, "object", typeStr, "typeStr has wrong content")

	_, found = omap.GetChildArray(nameMap, "required2")
	require.False(t, found, "didn't respond false for non existing array property")
	arrayThing, found := omap.GetChildArray(nameMap, "required")
	require.True(t, found, "didn't find required array")
	require.NotNil(t, arrayThing, "although required array was found it has nil as value")

	l := omap.GetArrayLen(arrayThing)
	require.Equal(t, 2, l, "required array has the wrong length")
	s, err := omap.GetValueAt[string](arrayThing, 1)
	require.Nil(t, err, "error while getting last elem")
	require.Equal(t, "last", s, "wrong last element in array")
	_, err = omap.GetValueAt[string](arrayThing, 2)
	require.NotNil(t, err, "no error in case of wrong array index")
	_, err = omap.GetValueAt[float32](arrayThing, 1)
	require.NotNil(t, err, "no error when requesting array elem with wrong type")

	requiredArray, found := omap.GetTypedChildArray[string](nameMap, "required")
	require.True(t, found, "didn't find required array (2)")
	require.Equal(t, []string{"first", "last"}, requiredArray, "requiredArray has wrong content")

	outputFile := "./temp/test_schema_output.json"
	err = mapThing.SerializeJSONFile(outputFile)
}

func TestIterators(t *testing.T) {
	fileName := "./resources/test_schema.json"
	mapThing, err := omap.NewFromJSONFile(fileName)
	require.Nil(t, err, "error while create object from file")
	require.NotNil(t, mapThing, "created instance is nil")
	found := maps.Collect(mapThing.IterateToValue(omap.OrderedPair2Value))
	require.NotNil(t, found, "iterator doesn't fill maps.Collect")
	require.Len(t, found, 7, "wrong number of items out of the iterator")

	requireArray, ok := found["required"]
	require.True(t, ok, "couldn't find 'required' key")

	a := make([]string, 0)
	iExpected := 0
	for i, v := range omap.IterateOverArray(requireArray, omap.OrderedValue2Value[string]) {
		require.Equal(t, iExpected, i, "wrong index")
		a = append(a, v)
		iExpected++
	}
	require.Equal(t, []string{"id", "name", "contact", "roles"}, a, "wrong extracted array")
}

func TestTemplateUse(t *testing.T) {
	fileName := "./resources/test_schema.json"
	mapThing, err := omap.NewFromJSONFile(fileName)
	require.Nil(t, err, "error while create object from file")
	require.NotNil(t, mapThing, "created instance is nil")

	templateFile := "./resources/test_template.tmpl"
	templateBytes, err := os.ReadFile(templateFile)

	tmpl := template.Must(template.New("TemplateTest").Funcs(
		template.FuncMap{
			"toText": omap.ToText,
		}).Parse(string(templateBytes)))
	var stringBuilder strings.Builder

	err = tmpl.Execute(&stringBuilder, mapThing)
	require.Nil(t, err, "error while process template", err)

	outputFile := "./temp/template_output.txt"
	file, err := os.Create(outputFile)
	require.Nil(t, err)
	defer file.Close()
	file.WriteString(stringBuilder.String())
}
