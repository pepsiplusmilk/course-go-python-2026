package roller

import "errors"

type Request struct {
	NumberOfDices      int
	NumberOfFaces      int
	NumberOfIterations int
	Seed               int64
}

func (r Request) Validate() error {
	if r.NumberOfDices < 1 {
		return errors.New("can't roll less than one dice")
	}

	if r.NumberOfDices > 100 {
		return errors.New("can't roll more than 100 dices")
	}

	if r.NumberOfFaces < 1 {
		return errors.New("can't roll dice that doesn't have faces")
	}

	if r.NumberOfFaces > 100 {
		return errors.New("can't roll dice with too much faces")
	}

	return nil
}
