/*
Copyright (C) 2018 Expedia Group.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package template

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomIntFunc(t *testing.T) {
	defer func(of func(int) int) { randomizer = of }(randomizer)
	randomizer = func(ub int) int { require.Equal(t, ub, 100); return 42 }

	val, err := Resolve("{{ randomInt(100) }}", nil)
	require.NoError(t, err)

	assert.Equal(t, "42", val)
}

func TestRandomAlphaShouldGenerateRandomAlphaStringOfRequestedSize(t *testing.T) {
	val, err := Resolve("{{ randomAlpha(10) }}", nil)
	require.NoError(t, err)

	assert.Regexp(t, "^[A-Za-z]{10}$", val.(string))
}

func TestRandomAlphaShouldReturnErrorIfLengthIsNegative(t *testing.T) {
	_, err := Resolve("{{ randomAlpha(-5) }}", nil)

	assert.EqualError(t, err, "error while evaluating expression: '{{ randomAlpha(-5) }}': word length must be non-negative")
}

func TestTemplateFunctionShouldResolveTemplateWithProvidedContext(t *testing.T) {

	context := Context{
		"ctx": map[string]interface{}{"team": map[string]string{"name": "Flyte"}},
		"tpl": `#!/bin/sh
echo "Hello {{ team.name }}">>{}`,
	}

	resolved, err := Resolve("{{ template(tpl, ctx) }}", context)
	require.NoError(t, err)

	want := `#!/bin/sh
echo "Hello Flyte">>{}`

	assert.Equal(t, want, resolved)
}

func TestTemplateFunctionInvalidEntry(t *testing.T) {
	template := `{{ template("{{ randomAlpha(-4) }}", ctx) }}`
	_, err := Resolve(template, Context{"ctx": map[string]interface{}{}})

	assert.EqualError(t, err, "error while evaluating expression: '{{ template(\"{{ randomAlpha(-4) }}\", ctx) }}': word length must be non-negative")
}

func TestDatastore(t *testing.T) {

	fn, ok := staticContext["datastore"]
	require.True(t, ok)
	defer func() { AddStaticContextEntry("datastore", fn) }()

	AddStaticContextEntry("datastore", func(key string) interface{} {
		return map[string]interface{}{"flyte": "flyte@flyte.com"}
	})

	template := `{email: "{{ datastore("product_teams").flyte }}"}`
	resolved, err := Resolve(template, Context{})

	require.NoError(t, err)
	assert.Equal(t, `{email: "flyte@flyte.com"}`, resolved.(string))
}

func TestMatchesCronExpression(t *testing.T) {
	matched, err := Resolve(`{{ "2018-02-14T23:18:09.0481031Z" | matchesCron: "18 23 * * *" }}`, nil)
	require.NoError(t, err)
	require.Equal(t, "True", matched)
}

func TestMatchesCronExpressionReturnsErrorWhenAnInvalidDateIsGiven(t *testing.T) {
	_, err := Resolve(`{{ "not a date" | matchesCron: "18 23 * * *" }}`, nil)
	require.Error(t, err)
	assert.EqualError(t, err, "error while evaluating expression: '{{ \"not a date\" | matchesCron: \"18 23 * * *\" }}': invalid date: 'not a date'")
}

func TestMatchesCronExpressionReturnsErrorWhenAnInvalidCronExpressionIsGiven(t *testing.T) {
	_, err := Resolve(`{{ "2018-02-14T23:18:09.0481031Z" | matchesCron: "not a cron expression" }}`, nil)
	require.Error(t, err)
	assert.EqualError(t, err, "error while evaluating expression: '{{ \"2018-02-14T23:18:09."+
		"0481031Z\" | matchesCron: \"not a cron expression\" }}': [Error | Line 1 Col 57 near 'matchesCron'] missing field(s)")
}
