package object

// Active records are used to store the current state of the evaluation in case
// of interruptions. Particularly used for yield interruptions, so the
// evaluator can resume the evaluation from the last active record.
type ActiveRecord interface {
}
