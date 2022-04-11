package assembler

type Category int8

const (
	CategoryEOF Category = iota
	CategoryNewLine
	CategoryWhitespace
	CategoryColon
	CategoryComma
	CategoryFullStop
	CategoryComment
	CategoryUpperName
	CategoryLowerName
	CategoryNumber
	CategoryAlign
	CategoryAllocate
	CategoryData
	CategoryStatement
	CategoryLabel
)

func (c Category) String() string {
	switch c {
	case CategoryEOF:
		return "EOF"
	case CategoryNewLine:
		return "NewLine"
	case CategoryWhitespace:
		return "Whitespace"
	case CategoryColon:
		return "Colon"
	case CategoryComma:
		return "Comma"
	case CategoryFullStop:
		return "FullStop"
	case CategoryComment:
		return "Comment"
	case CategoryUpperName:
		return "UpperName"
	case CategoryLowerName:
		return "LowerName"
	case CategoryNumber:
		return "Number"
	case CategoryAlign:
		return "Align"
	case CategoryAllocate:
		return "Allocate"
	case CategoryData:
		return "Data"
	case CategoryStatement:
		return "Statement"
	case CategoryLabel:
		return "Label"
	default:
		return "Unrecognized Category"
	}
}
