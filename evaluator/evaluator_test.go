// Copyright 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evaluator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/evaluator"
)

type TestCase struct {
	name          string
	expression    string
	expectedValue int
	expectedError bool
}

// TestEvaluatorEmptyInput function checks the evaluator.Evaluate function for
// empty input
func TestEvaluatorEmptyInput(t *testing.T) {
	var values = make(map[string]int)
	expression := ""

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEvaluatorSingleToken function checks the evaluator.Evaluate function for
// single token input
func TestEvaluatorSingleToken(t *testing.T) {
	var values = make(map[string]int)
	expression := "42"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 42, result)
}

// TestEvaluatorArithmetic checks the evaluator.Evaluate function for simple
// arithmetic expression
func TestEvaluatorArithmetic(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "short expression",
			expression:    "1+2*3",
			expectedValue: 7,
		},
		{
			name:          "long expression",
			expression:    "4/2-1+5%2",
			expectedValue: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestEvaluatorParenthesis checks the evaluator.Evaluate function for simple
// arithmetic expression with parenthesis
func TestEvaluatorParenthesis(t *testing.T) {
	var values = make(map[string]int)
	expression := "(1+2)*3"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 9, result)
}

// TestEvaluatorRelational checks the evaluator.Evaluate function for simple
// relational expression
func TestEvaluatorRelational(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "less than",
			expression:    "1 < 2",
			expectedValue: 1,
		},
		{
			name:          "greater or equal",
			expression:    "1 >= 2",
			expectedValue: 0,
		},
		{
			name:          "long expression",
			expression:    "1 < 2 && 1 > 2 && 1 <= 2 && 1 >= 2 && 1==2 && 1 != 2",
			expectedValue: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestEvaluatorBoolean checks the evaluator.Evaluate function for simple
// boolean expressions
func TestEvaluatorBoolean(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "and",
			expression:    "1 && 0",
			expectedValue: 0,
		},
		{
			name:          "or",
			expression:    "1 || 0",
			expectedValue: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestEvaluatorValues checks the evaluator.Evaluate function for expression
// with named values
func TestEvaluatorValues(t *testing.T) {
	var values = make(map[string]int)
	values["x"] = 1
	values["y"] = 2
	expression := "x+y*2"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 5, result)
}

// TestEvaluatorWrongInput checks the evaluator.Evaluate function for
// expression that is not correct
func TestEvaluatorWrongInput(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "mul instead of right operand",
			expression:    "1**",
			expectedError: true,
		},
		{
			name:          "forgot closing parenthesis",
			expression:    "(1+2*",
			expectedError: true,
		},
		{
			name:          "no operands",
			expression:    "+",
			expectedError: true,
		},
		{
			name:          "no right operand",
			expression:    "2+",
			expectedError: true,
		},
		{
			name:          "no left operand",
			expression:    "+2",
			expectedError: true,
		},
		{
			name:          "no left operand (minus)",
			expression:    "-2",
			expectedError: true,
		},
		{
			name:          "== typo",
			expression:    "0=0",
			expectedError: true,
		},
		{
			name:          "zero division",
			expression:    "1/0",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedError {
				result, err := evaluator.Evaluate(tc.expression, values)
				assert.Error(t, err, "error is expected")
				assert.Equal(t, -1, result)
			}
		})
	}
}

// TestEvaluatorMissingValue checks the evaluator.Evaluate function for
// expression that use value not provided
func TestEvaluatorMissingValue(t *testing.T) {
	var values = make(map[string]int)
	expression := "value"

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEdgeCases tests expressions that rarely happen
// in the real world
func TestEdgeCases(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "useless parenthesis",
			expression:    "(2)*(2)",
			expectedValue: 4,
		},
		{
			name:          "multiple useless parenthesis",
			expression:    "(((0))) >= 0",
			expectedValue: 1,
		},
		{
			name:          "scrambled useless parenthesis",
			expression:    "((((0==0)))+1)",
			expectedValue: 2,
		},
		{
			name:          "0 addition idempotence",
			expression:    "1+0+0+0+0",
			expectedValue: 1,
		},
		{
			name:          "1 division idempotence",
			expression:    "5/1/1/1/1/1",
			expectedValue: 5,
		},
		{
			name:          "transitivity",
			expression:    "(3 > 2) && (2 > 1) == (3 > 1)",
			expectedValue: 1,
		},
		{
			name:          "big integer",
			expression:    "9223372036854775807+100-100",
			expectedValue: 9223372036854775807,
		},
		{
			name:          "overflow",
			expression:    "9223372036854775807+1",
			expectedValue: -9223372036854775808,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}