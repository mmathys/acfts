package common

import "errors"

/**
Checks the format of the request
- Fields must be defined
- there must be at least one input and one output field
- all values must be greater than zero
- all inputs must have at least one signature
- the mode of the signature must be valid TODO
- all signatures must have the same mode
*/
func CheckFormat(tx *Transaction) error {
	if tx == nil || tx.Outputs == nil || tx.Inputs == nil {
		return errors.New("invalid format")
	}

	if len(tx.Outputs) == 0 {
		return errors.New("there must be at least one output")
	}

	if len(tx.Inputs) == 0 {
		return errors.New("there must be a least on input")
	}

	// signature mode
	var mode int
	if len(tx.Inputs[0].Signatures) > 0 {
		mode = tx.Inputs[0].Signatures[0].Mode
	} else {
		return errors.New("encountered input without signatures")
	}

	for _, input := range tx.Inputs {
		if input.Amount <= 0 {
			return errors.New("input must be greater than zero")
		}
		if len(input.Signatures) == 0 {
			return errors.New("encountered input without signatures")
		}
		for _, sig := range input.Signatures {
			if sig.Mode != mode {
				return errors.New("encountered inconsistent (or invalid) signature mode")
			}
		}
	}

	for _, output := range tx.Outputs {
		if output.Amount <= 0 {
			return errors.New("output must be greater than zero")
		}
	}

	return nil
}

/**
Checks constraints, such as:
- Sum of inputs = sum of outputs
*/
func CheckConstraints(tx *Transaction) error {
	sumInputs := 0
	sumOutputs := 0

	for _, input := range tx.Inputs {
		sumInputs += input.Amount
	}

	for _, output := range tx.Outputs {
		sumOutputs += output.Amount
	}

	if sumInputs != sumOutputs {
		return errors.New("sum of inputs must equal sum of outputs")
	}

	return nil
}
