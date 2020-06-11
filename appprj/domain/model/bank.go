package model

import (
	"fmt"

	"go-prj-skeleton/appprj/strutil"
)

var (
	Banks = []string{"VCB", "ACB", "VIB"}

	ErrInvalidBank = fmt.Errorf("invalid bank")
)

func ValidateBank(bank string) error {
	if bank == "" {
		return fmt.Errorf("%s: %w", bank, ErrInvalidBank)
	}

	if !strutil.Include(Banks, bank) {
		return fmt.Errorf("%s: %w", bank, ErrInvalidBank)
	}

	return nil
}
