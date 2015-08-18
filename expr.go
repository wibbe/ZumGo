package main

type OperandType int

const (
	OpConstant OperandType = iota
	OpCellRef
	OpRangeRef
	OpFunction
)

type Operand interface {
	Value(doc *Document) float64
	Type() OperandType
}

type Expr struct {
	expression []Operand
}

func ParseExpression(text string) (*Expr, error) {
	return nil, nil
}

func (e *Expr) String() string {
	return "="
}

func (e *Expr) Evaluate(doc *Document) (float64, error) {
	argStack := make([]Operand, 0, len(e.expression)/2)

	for _, op := range e.expression {
		switch op.Type() {
		case OpConstant:
		case OpCellRef:
		case OpRangeRef:
			argStack = append(argStack, op)
		case OpFunction:
			//f := op.(*Function)
		}
	}

	return 0.0, nil
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
	return OpConstant
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
	return OpConstant
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
	return OpRangeRef
}

type Function struct {
}

func (f *Function) Value(doc *Document) float64 {
	return 0.0
}

func (f *Function) Type() OperandType {
	return OpFunction
}
