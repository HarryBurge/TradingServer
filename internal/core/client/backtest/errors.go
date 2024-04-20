package backtest

type SymbolNotFoundError struct {
}

func (e SymbolNotFoundError) Error() string {
	return "symbol not found"
}

type IncorrectTimeError struct {
}

func (e IncorrectTimeError) Error() string {
    return "incorrect time"
}

type IdNotFoundError struct {
}

func (e IdNotFoundError) Error() string {
    return "id not found"
}