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

package translate

import (
	"fmt"

	"bytes"

	"strings"

	"github.com/DE-labtory/koa/opcode"
)

// Asm is generated by compiling.
type Asm struct {
	AsmCodes []AsmCode
}

type AsmCode struct {
	RawByte []byte
	Value   string
}

// Emerge() translates instruction to bytecode
// An operand of operands should be 4 bytes.
func (a *Asm) Emerge(operator opcode.Type, operands ...[]byte) int {
	asmCode, err := convert(operator, operands...)
	if err != nil {
		return 0
	}

	a.AsmCodes = append(a.AsmCodes, asmCode...)
	return len(a.AsmCodes)
}

// EmergeAt() translates instruction to bytecode and append at index
// An operand of operands should be 4 bytes.
func (a *Asm) EmergeAt(index int, operator opcode.Type, operands ...[]byte) int {
	asmCode, err := convert(operator, operands...)
	if err != nil {
		return 0
	}

	a.AsmCodes = append(a.AsmCodes[:index], append(asmCode, a.AsmCodes[index:]...)...)
	return len(a.AsmCodes)
}

func (a *Asm) ReplaceOperandAt(index int, operands []byte) error {
	a.AsmCodes[index] = AsmCode{
		Value:   fmt.Sprintf("%x", operands),
		RawByte: operands,
	}
	return nil
}

func (a *Asm) ReplaceOperatorAt(index int, operator opcode.Type) error {
	opStr, err := operator.String()
	if err != nil {
		return err
	}

	a.AsmCodes[index] = AsmCode{
		Value:   opStr,
		RawByte: []byte{byte(operator)},
	}
	return nil
}

func (a *Asm) Equal(a1 Asm) bool {
	if len(a.AsmCodes) != len(a1.AsmCodes) {
		return false
	}

	for i, asm := range a.AsmCodes {
		if asm.Value != a1.AsmCodes[i].Value {
			return false
		}

		if !bytes.Equal(asm.RawByte, a1.AsmCodes[i].RawByte) {
			return false
		}
	}

	return true
}

func (a *Asm) ToRawByteCode() []byte {
	result := make([]byte, 0)
	for _, code := range a.AsmCodes {
		result = append(result, code.RawByte...)
	}
	return result
}

func (a *Asm) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	s := []string{}
	for _, code := range a.AsmCodes {
		s = append(s, code.Value)
	}
	out.WriteString(strings.Join(s, " "))
	out.WriteString("]")

	return out.String()
}

func convert(operator opcode.Type, operands ...[]byte) ([]AsmCode, error) {
	// Translate operator to byte
	asmCodes := make([]AsmCode, 0)

	// Translate operator to assembly
	opStr, err := operator.String()
	if err != nil {
		return nil, err
	}

	asmCodes = append(asmCodes, AsmCode{
		Value:   opStr,
		RawByte: []byte{byte(operator)},
	})

	for _, o := range operands {
		asmCodes = append(asmCodes, AsmCode{
			Value:   fmt.Sprintf("%x", o),
			RawByte: o,
		})
	}

	return asmCodes, nil
}
