package bot

type state = int

const (
	listening state = iota
	configuring
)
