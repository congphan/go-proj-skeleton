package postgre

import (
	"github.com/pkg/errors"

	"go-prj-skeleton/app/pgutil"
)

type pgHelperStruct struct {
}

var pgHelper = pgHelperStruct{}

func (helper pgHelperStruct) delete(model interface{}) error {
	db := pgutil.DB()
	err := db.Delete(model)
	if err != nil {
		return errors.Wrap(err, "delete failed")
	}

	return nil
}
