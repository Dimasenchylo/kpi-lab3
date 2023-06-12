package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Dimasenchylo/kpi-lab3/painter"
)

type UIState struct {
	BackgroundColor painter.Operation
	BackgroundRect  *painter.BgRectangle
	Figures         []painter.Figure
	MoveOps         []painter.Operation
	UpdateOp        painter.Operation
}

func (s *UIState) Reset() {
	s.BackgroundColor = nil
	s.BackgroundRect = nil
	s.Figures = nil
	s.MoveOps = nil
	s.UpdateOp = nil
}

type Parser struct {
	State *UIState
}

func NewParser() *Parser {
	return &Parser{
		State: &UIState{},
	}
}

// Parse reads and parses input from the provided io.Reader and returns the corresponding list of painter.Operation.
func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.State.Reset()
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() { // loop through the input stream using the scanner
		commandLine := scanner.Text()

		err := p.parse(commandLine) // parse the command line into an operation
		if err != nil {
			return nil, err
		}
	}
	return p.finalResult(), nil
}

func (p *Parser) finalResult() []painter.Operation {
	var res []painter.Operation
	if p.State.BackgroundColor != nil {
		res = append(res, p.State.BackgroundColor)
	}
	if p.State.BackgroundRect != nil {
		res = append(res, p.State.BackgroundRect)
	}
	if len(p.State.MoveOps) != 0 {
		res = append(res, p.State.MoveOps...)
	}
	p.State.MoveOps = nil
	if len(p.State.Figures) != 0 {
		for _, figure := range p.State.Figures {
			res = append(res, &figure)
		}
	}
	if p.State.UpdateOp != nil {
		res = append(res, p.State.UpdateOp)
	}
	return res
}

func (p *Parser) parse(commandLine string) error {
	parts := strings.Split(commandLine, " ")
	instruction := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	var iArgs []int
	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err == nil {
			iArgs = append(iArgs, i)
		}
	}

	switch instruction {
	case "white":
		p.State.BackgroundColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.State.BackgroundColor = painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		p.State.BackgroundRect = &painter.BgRectangle{X1: iArgs[0], Y1: iArgs[1], X2: iArgs[2], Y2: iArgs[3]}
	case "figure":
		figure := painter.Figure{X: iArgs[0], Y: iArgs[1]}
		p.State.Figures = append(p.State.Figures, figure)
	case "move":
		moveOp := painter.Move{X: iArgs[0], Y: iArgs[1], Figures: p.State.Figures}
		p.State.MoveOps = append(p.State.MoveOps, &moveOp)
	case "reset":
		p.State.Reset()
		p.State.BackgroundColor = painter.OperationFunc(painter.ClearScreen)
	case "update":
		p.State.UpdateOp = painter.UpdateOp
	default:
		return fmt.Errorf("error with parse %v", commandLine)
	}
	return nil
}
