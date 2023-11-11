package parser

func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}
