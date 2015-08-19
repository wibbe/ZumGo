package main

import (
	"errors"
	//	"go/scanner"
	"go/token"
	//	"strconv"
)

type OperandType int

const (
	ConstantSymbol OperandType = iota
	CellSymbol
	RangeSymbol
	FunctionSymbol
)

type Symbol interface {
	Value(doc *Document) float64
	Type() OperandType
}

type SymbolFunc func(*Document, []Symbol) []Symbol

type Expr struct {
	expression []Symbol
}

type operatorPair struct {
	tok token.Token
	lit string
}

func ParseExpression(text string) (*Expr, error) {
	/*
		output := make([]Symbol)
		operatorStack := make([]operatorPair, 0, 10)

		var s scanner.Scanner

		fileSet := token.NewFileSet()
		file := fileSet.AddFile("", fileSet.Base(), len(text))
		s.Init(file, []byte(text), nil, scanner.ScanComments)

		for {
			_, tok, lit := s.Scan()

			switch tok {
			case token.INT:
				number, err := strconv.ParseFloat(lit, 64)
				if err != nil {
					return nil, err
				}
				output = append(output, NewConstant(number))

			case token.FLOAT:
				number, err := strconv.ParseFloat(lit, 64)
				if err != nil {
					return nil, err
				}
				output = append(output, NewConstant(number))

			case token.IDENT:
			}

			if tok == token.EOF {
				break
			}
		}

		return output, nil
	*/
	return nil, nil
}

func (e *Expr) String() string {
	return ""
}

func (e *Expr) Evaluate(doc *Document) (float64, error) {
	argStack := make([]Symbol, 0, len(e.expression)/2)

	for _, symbol := range e.expression {
		switch symbol.Type() {
		case ConstantSymbol:
		case CellSymbol:
		case RangeSymbol:
			argStack = append(argStack, symbol)
		case FunctionSymbol:
			f := symbol.(*Function)
			argStack = f.cmd(doc, argStack)
		}
	}

	if len(argStack) != 0 {
		return 0.0, errors.New("error while evaluating expression '" + e.String() + "'")
	}

	return argStack[0].Value(doc), nil
}

type Constant struct {
	value float64
}

func NewConstant(value float64) *Constant {

	return &Constant{value: value}
}

func (c *Constant) Value(doc *Document) float64 {
	return c.value
}

func (c *Constant) Type() OperandType {
	return ConstantSymbol
}

type CellRef struct {
	index Index
}

func NewCellRef(idx Index) *CellRef {
	return &CellRef{index: idx}
}

func (c *CellRef) Value(doc *Document) float64 {
	cell, err := doc.GetCell(c.index)
	if err != nil {
		return 0.0
	}

	value, err := cell.Eval(doc)
	if err != nil {
		return 0.0
	}

	return value
}

func (c *CellRef) Type() OperandType {
	return ConstantSymbol
}

type RangeRef struct {
	start Index
	end   Index
}

func NewRangeRef(start, end Index) *RangeRef {
	return &RangeRef{start: start, end: end}
}

func (r *RangeRef) Value(doc *Document) float64 {
	return 0.0
}

func (r *RangeRef) Type() OperandType {
	return RangeSymbol
}

type Function struct {
	cmd SymbolFunc
}

func (f *Function) Value(doc *Document) float64 {
	return 0.0
}

func (f *Function) Type() OperandType {
	return FunctionSymbol
}
