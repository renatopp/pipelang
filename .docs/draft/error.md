# Error

- `?` catch error and returning a Maybe
- `!` deconstruct the Maybe into `err, val`
- `??` shortcut to `<left> = Wrap<left>; <left>.Val() if <left>.Ok() else <right>`

	```go
	// The value if ok, raise if not
	f := file.Open()

	// Maybe
	f := file.Open()?
	if f.Err() {
		f = 'stub file' // invalid type change
	}

	// Unwrapping
	f := file.Open()?
	if f! {
		f = file.Open()
	}

  // Piping
	config := file.Open('txt')?
	| or :file.Open('.env')?
	| or :file.Open('.stuff')
	```


## reais

```go
	source, err := file.LoadSource()
	if err != nil {
		source = []byte{}
	}

	// viraria
	source := file.LoadSource() ?? []byte{}
```


```go

	ast, err := file.LoadAst()
	if err != nil {
		return nil, err
	}

	// viraria
	err, ast := file.LoadAst()?!
	if err {
		raise errorf('Blbalblaa')
	}

```


```go
	postLexer := internal.NewPostLexer(tokens)
	if err := postLexer.Optimize(); err != nil {
		logs.Print("[sourcefile] error post-lexing (%v)", err)
		return nil, err
	}

	if postLexer.Optimize()! as err {
		logs.Print("[sourcefile] error post-lexing (%v)", err)
		raise err
	}
```