package controller

import (
	"github.com/lance/js-function-operator/pkg/controller/jsfunction"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, jsfunction.Add)
}
