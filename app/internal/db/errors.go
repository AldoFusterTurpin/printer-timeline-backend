package db

import (
	. "bitbucket.org/aldoft/printer-timeline-backend/app/internal/errorTypes"
)

const (
	NotFoundErr       = ConstError("element not found in database")
	ConditionalPutErr = ConstError("ConstError trying to add duplicated item")
	ConditionalDelErr = ConstError("ConstError deleting Item. It not exist")

	ExpressionBuilderErr = ConstError("element not found in database")
	MarshallErr          = ConstError("ConstError marshalling data")
	UnmarshallErr        = ConstError("ConstError unmarshalling data")
)
