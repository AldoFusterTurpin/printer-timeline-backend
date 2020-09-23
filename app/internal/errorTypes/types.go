// package errorTypes defines the different errorTypes of the application.
package errorTypes

type ConstError string

func (err ConstError) Error() string {
	return string(err)
}
