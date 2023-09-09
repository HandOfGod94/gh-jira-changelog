package fsm_util

type State = string
type Event = string

func Before(event Event) string {
	return "before_" + event
}
