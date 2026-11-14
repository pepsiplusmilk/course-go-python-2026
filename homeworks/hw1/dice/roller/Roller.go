package roller

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type PrintMode int

const (
	PrintModeDefault PrintMode = iota
	PrintModeOnlySum
	PrintModeVerbose
)

type RollResult struct {
	DiceScores        []int
	Sum               []int
	Average           []float64
	Mode              PrintMode
	IterationsDone    int
	NumberOfDiceFaces int
}

func (r *RollResult) Roll(request *Request) error {
	if r == nil {
		return errors.New("structure that must holds a roll result isn't defined(got a nil instead)")
	}

	gen := rand.New(rand.NewSource(request.Seed))

	r.IterationsDone = request.NumberOfIterations
	r.NumberOfDiceFaces = request.NumberOfFaces

	r.DiceScores = make([]int, request.NumberOfDices*request.NumberOfIterations)
	r.Sum = make([]int, request.NumberOfIterations)
	r.Average = make([]float64, request.NumberOfIterations)

	for i := 0; i < request.NumberOfIterations; i++ {
		for j := 0; j < request.NumberOfDices; j++ {
			newRoll := gen.Intn(request.NumberOfFaces) + 1
			r.DiceScores[i*request.NumberOfDices+j] = newRoll

			r.Sum[i] += newRoll
		}

		r.Average[i] = float64(r.Sum[i]) / float64(request.NumberOfDices)
	}

	return nil
}

func (r RollResult) printOneIter(iter int, singeIterationWasDone bool) string {
	sep := ""
	if !singeIterationWasDone {
		sep = "  "
	}

	var builder strings.Builder
	numberOfDices := len(r.DiceScores) / r.IterationsDone

	if r.Mode == PrintModeDefault {
		builder.WriteString(sep)
		builder.WriteString("Rolls: ")

		if numberOfDices == 1 {
			builder.WriteString(strconv.Itoa(r.Sum[iter]))
		} else {
			jter := iter * numberOfDices

			for i, roll := range r.DiceScores[jter : jter+numberOfDices] {
				builder.WriteString(strconv.Itoa(roll))

				if i != jter+numberOfDices-1 {
					builder.WriteString("+")
				} else {
					builder.WriteString("=")
				}
			}

			builder.WriteString(strconv.Itoa(r.Sum[iter]))
		}
	} else if r.Mode == PrintModeOnlySum {
		builder.WriteString(sep)
		builder.WriteString("Sum: ")
		builder.WriteString(strconv.Itoa(r.Sum[iter]))
	} else if r.Mode == PrintModeVerbose {
		jter := iter * numberOfDices

		for i, roll := range r.DiceScores[jter : jter+numberOfDices] {
			builder.WriteString(sep)
			builder.WriteString(fmt.Sprintf("Dice %d: ", i+1))
			builder.WriteString(strconv.Itoa(roll))
			builder.WriteString("\n")
		}

		// print sum
		builder.WriteString(sep)
		builder.WriteString("Sum: ")

		builder.WriteString(strconv.Itoa(r.Sum[iter]))

		builder.WriteString("\n")
		// print mean value
		builder.WriteString(sep)
		builder.WriteString("Average: ")

		builder.WriteString(strconv.FormatFloat(r.Average[iter], 'f', -1, 64))
	}

	builder.WriteString("\n")
	return builder.String()
}

func (r RollResult) ToString() string {
	numberOfDices := len(r.DiceScores) / r.IterationsDone

	var builder strings.Builder

	if r.Mode == PrintModeVerbose {
		builder.WriteString(fmt.Sprintf("Rolling %dd%d", numberOfDices, r.NumberOfDiceFaces))

		if r.IterationsDone > 1 {
			builder.WriteString(fmt.Sprintf(" (%d iterations):\n", r.IterationsDone))
		} else {
			builder.WriteString(":\n")
		}
	}

	if r.IterationsDone == 1 {
		builder.WriteString(r.printOneIter(0, true))
	} else {
		for i := range r.IterationsDone {
			builder.WriteString(fmt.Sprintf("Iteration %d:\n", i+1))
			builder.WriteString(r.printOneIter(i, false))
		}
	}

	return builder.String()
}
