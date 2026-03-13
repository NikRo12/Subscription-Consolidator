package commonerrors

import (
	"errors"
)

/*
This error says that the root directory doesn't have .env-file
*/
var NoENVFile = errors.New("Cannot find .env file in project's root directory")
