/*
 * Copyright 2018 De-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parse_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/DE-labtory/koa/ast"
	"github.com/DE-labtory/koa/parse"
)

// expectedFnArg is used to verifing parsed function args data
type expectedFnArg struct {
	argIdent string
	argType  ast.DataStructure
}

// expectedFnArg is used to verifing parsed function header data
type expectedFnHeader struct {
	retType ast.DataStructure
	args    []expectedFnArg
}

// fnTmplData represents smart contract function which is going to
// be injected into contractTmpl
type fnTmplData struct {
	FuncName string
	Args     string
	RetType  string
	Stmts    []string
}

// contractTmplData represents smart contract which is going to
// be injected into contractTmpl
type contractTmplData struct {
	Fns []fnTmplData
}

// contractTmpl is template for creating smart contract code
// contractTmplData injects data to this template
const contractTmpl = `
contract {
	{{range .Fns -}}
		func {{.FuncName}}({{.Args}}) {{.RetType}} {
        {{range .Stmts -}}
			{{.}}
        {{end}}
    }
	{{end}}
}
`

// createTestContractCode creates and returns string code, which is made by
// contractTmpl and contractTmplData
func createTestContractCode(c contractTmplData) string {
	out := bytes.NewBufferString("")
	instance, _ := template.New("ContractTemplate").Parse(contractTmpl)
	instance.Execute(out, c)

	return out.String()
}

func parseTestContract(input string) (*ast.Contract, error) {
	l := parse.NewLexer(input)
	buf := parse.NewTokenBuffer(l)
	return parse.Parse(buf)
}

// chkFnHeader verify smart contract's function header
func chkFnHeader(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader) {
	if fn.ReturnType != efh.retType {
		t.Errorf("testReturn() has wrong return type expected=%v, got=%v",
			ast.VoidType.String(), fn.ReturnType.String())
	}

	if len(fn.Parameters) != len(efh.args) {
		t.Errorf("testReturn() has wrong parameters length got=%v",
			len(fn.Parameters))
	}

	for i, a := range fn.Parameters {
		testFnParameters(t, a, efh.args[i].argType, efh.args[i].argIdent)
	}
}

// TODO: 1. add multi function contract test cases
// TODO: 2. test edge cases - type mismatch, return undefined variable
func TestReturnStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.ReturnStatement
	}{
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement1",
						Args:     "",
						RetType:  "",
						Stmts: []string{
							"return 1",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.VoidType,
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.IntegerLiteral{Value: 1},
				},
			},
		},
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement2",
						Args:     "a int",
						RetType:  "int",
						Stmts: []string{
							"return a",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Value: "a"},
				},
			},
		},
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement3",
						Args:     "a int, b int",
						RetType:  "int",
						Stmts: []string{
							"return a + b + 1",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
						{"b", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.InfixExpression{
						Left: &ast.InfixExpression{
							Left:     &ast.Identifier{Value: "a"},
							Operator: ast.Plus,
							Right:    &ast.Identifier{Value: "b"},
						},
						Operator: ast.Plus,
						Right:    &ast.IntegerLiteral{Value: 1},
					},
				},
			},
		},
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement4",
						Args:     "a int",
						RetType:  "int",
						Stmts: []string{
							"return (",
							"a)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Value: "a"},
				},
			},
		},
		{
			contractTmpl: contractTmplData{
				Fns: []fnTmplData{
					{
						FuncName: "returnStatement4",
						Args:     "a int",
						RetType:  "int",
						Stmts: []string{
							"return (",
							"a)",
						},
					},
					{
						FuncName: "returnStatement4_1",
						Args:     "b int",
						RetType:  "int",
						Stmts: []string{
							"return (",
							"b)",
						},
					},
				},
			},
			expectedFnHeaders: []expectedFnHeader{
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"a", ast.IntType},
					},
				},
				{
					retType: ast.IntType,
					args: []expectedFnArg{
						{"b", ast.IntType},
					},
				},
			},
			expected: []ast.ReturnStatement{
				{
					ReturnValue: &ast.Identifier{Value: "a"},
				},
				{
					ReturnValue: &ast.Identifier{Value: "b"},
				},
			},
		},
	}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runReturnStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runReturnStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efhs expectedFnHeader, tt ast.ReturnStatement) {
	t.Logf("test ReturnStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efhs)

	for _, stmt := range fn.Body.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", stmt)
		}

		testExpression2(t, returnStmt.ReturnValue, tt.ReturnValue)
	}
}

// TODO: add test cases
func TestAssignStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.AssignStatement
	}{}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runAssignStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}

	}
}

func runAssignStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.AssignStatement) {
	t.Logf("test AssignStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.ReturnStatement. got=%T", stmt)
		}

		testAssignStatement(t, assignStmt, tt.Type, tt.Variable, tt.Value)
	}
}

// TODO: add test cases
func TestIfElseStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.IfStatement
	}{}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runIfStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runIfStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.IfStatement) {
	t.Logf("test IfStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		ifStmt, ok := stmt.(*ast.IfStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.IfStatement. got=%T", stmt)
		}
		testIfStatement(t, ifStmt, tt.Condition, tt.Consequence, tt.Alternative)
	}
}

// TODO: add test cases
func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		contractTmpl      contractTmplData
		expectedFnHeaders []expectedFnHeader
		expected          []ast.ExpressionStatement
	}{}

	for _, tt := range tests {
		input := createTestContractCode(tt.contractTmpl)
		contract, err := parseTestContract(input)

		if err != nil {
			t.Errorf("parser error: %q", err)
			t.FailNow()
		}

		for i, fn := range contract.Functions {
			runExpressionStatementTestCases(t, fn, tt.expectedFnHeaders[i], tt.expected[i])
		}
	}
}

func runExpressionStatementTestCases(t *testing.T, fn *ast.FunctionLiteral, efh expectedFnHeader, tt ast.ExpressionStatement) {
	t.Logf("test ExpressionStatement - [%s]", fn.Name)

	chkFnHeader(t, fn, efh)

	for _, stmt := range fn.Body.Statements {
		exprStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("function body stmt is not *ast.IfStatement. got=%T", stmt)
		}
		testExpression2(t, exprStmt.Expr, tt.Expr)
	}
}

func testExpression2(t *testing.T, exp ast.Expression, expected ast.Expression) {
	if exp.String() != expected.String() {
		t.Errorf(`expression is not "%s", bot got "%s"`, exp.String(), expected.String())
	}
}

func testFnParameters(t *testing.T, p *ast.ParameterLiteral, ds ast.DataStructure, id string) {
	if p.Type != ds {
		t.Errorf("wrong parameter type expected=%s, got=%s",
			p.Type.String(), ds.String())
	}
	if p.Identifier.Value != id {
		t.Errorf("wrong parameter identifier expected=%s, got=%s",
			p.Type.String(), ds.String())
	}
}

func testAssignStatement(t *testing.T, stmt *ast.AssignStatement, ds ast.DataStructure, ident ast.Identifier, value ast.Expression) {
	if stmt.Type != ds {
		t.Errorf("wrong assign statement type expected=%T, got=%T",
			ds, stmt.Type)
	}
	if stmt.Variable.String() != ident.String() {
		t.Errorf("wrong assign statement variable expected=%s, got=%s",
			ident.String(), stmt.Variable.String())
	}
	if stmt.Value.String() != value.String() {
		t.Errorf("wrong assign statement value expected=%s, got=%s",
			value.String(), stmt.Value.String())
	}
}

func testIfStatement(t *testing.T, stmt *ast.IfStatement, condition ast.Expression, consequences *ast.BlockStatement, alternatives *ast.BlockStatement) {
	if stmt.Condition.String() != condition.String() {
		t.Errorf("wrong condition statement type expected=%s, got=%s",
			condition, stmt.Condition.String())
	}

	if len(stmt.Consequence.Statements) != len(consequences.Statements) {
		t.Errorf("wrong condition statement consequences length expected=%d, got=%d",
			len(consequences.Statements), len(stmt.Consequence.Statements))
	}
	for i, csq := range stmt.Consequence.Statements {
		if csq.String() != consequences.Statements[i].String() {
			t.Errorf("wrong condition statement consequences literal expected=%s, got=%s",
				csq.String(), consequences.Statements[i].String())
		}
	}

	if stmt.Alternative == nil {
		return
	}
	if len(stmt.Alternative.Statements) != len(alternatives.Statements) {
		t.Errorf("wrong condition statement alternatives length expected=%d, got=%d",
			len(alternatives.Statements), len(stmt.Alternative.Statements))
	}
	for i, alt := range stmt.Alternative.Statements {
		if alt.String() != alternatives.Statements[i].String() {
			t.Errorf("wrong condition statement alternative literal expected=%s, got=%s",
				alt.String(), alternatives.Statements[i].String())
		}
	}
}
