f := file.Read('invalid')?

if !f.Ok() {
	-- Treat error
}

f.Value()