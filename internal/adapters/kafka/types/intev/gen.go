package intev

//go:generate protoc -I=. -I../../../../../ --go_out=. --go_opt=paths=source_relative --dispatcher_out=. internal_events.proto
