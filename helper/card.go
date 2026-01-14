package helper

import (
	"strconv"
	"strings"
)

type StructCardBrandPattern struct {
	Start int
	End   int
}

type StructCardBrand struct {
	NiceType string
	Type     string
	Patterns []StructCardBrandPattern
}

var cardTypes = []StructCardBrand{
	{
		NiceType: "Visa",
		Type:     "visa",
		Patterns: []StructCardBrandPattern{
			{4, 4},
		},
	},
	{
		NiceType: "Mastercard",
		Type:     "mastercard",
		Patterns: []StructCardBrandPattern{
			{51, 55}, {2221, 2229}, {223, 229}, {23, 26}, {270, 271}, {2720, 2720},
		},
	},
	{
		NiceType: "Amex",
		Type:     "american-express",
		Patterns: []StructCardBrandPattern{
			{34, 34}, {37, 37},
		},
	},
	{
		NiceType: "Diners Club",
		Type:     "diners-club",
		Patterns: []StructCardBrandPattern{
			{300, 305}, {36, 36}, {38, 39},
		},
	},
	{
		NiceType: "Discover",
		Type:     "discover",
		Patterns: []StructCardBrandPattern{
			{6011, 6011}, {644, 649}, {65, 65},
		},
	},
	{
		NiceType: "JCB",
		Type:     "jcb",
		Patterns: []StructCardBrandPattern{
			{2131, 2131}, {1800, 1800}, {3528, 3589},
		},
	},
	{
		NiceType: "UnionPay",
		Type:     "unionpay",
		Patterns: []StructCardBrandPattern{
			{620, 620}, {62100, 62182}, {62184, 62187}, {62185, 62197},
			{62200, 62205}, {622010, 622999}, {622018, 622018},
			{62207, 62209}, {623, 626},
			{6270, 6270}, {6272, 6272}, {6276, 6276},
			{627700, 627779}, {627781, 627799},
			{6282, 6289}, {6291, 6291}, {6292, 6292},
			{810, 810}, {8110, 8131}, {8132, 8151},
			{8152, 8163}, {8164, 8171},
		},
	},
	{
		NiceType: "Maestro",
		Type:     "maestro",
		Patterns: []StructCardBrandPattern{
			{493698, 493698}, {500000, 504174}, {504176, 506698},
			{506779, 508999}, {56, 59}, {63, 63}, {67, 67}, {6, 6},
		},
	},
	{
		NiceType: "Elo",
		Type:     "elo",
		Patterns: []StructCardBrandPattern{
			{401178, 401178}, {401179, 401179}, {438935, 438935},
			{457631, 457631}, {457632, 457632}, {431274, 431274},
			{451416, 451416}, {457393, 457393}, {504175, 504175},
			{506699, 506778}, {509000, 509999}, {627780, 627780},
			{636297, 636297}, {636368, 636368},
			{650031, 650033}, {650035, 650051},
			{650405, 650439}, {650485, 650538},
			{650541, 650598}, {650700, 650718},
			{650720, 650727}, {650901, 650978},
			{651652, 651679}, {655000, 655019},
			{655021, 655058},
		},
	},
	{
		NiceType: "Mir",
		Type:     "mir",
		Patterns: []StructCardBrandPattern{
			{2200, 2204},
		},
	},
	{
		NiceType: "Hiper",
		Type:     "hiper",
		Patterns: []StructCardBrandPattern{
			{637095, 637095}, {63737423, 63737423}, {63743358, 63743358},
			{637568, 637568}, {637599, 637599}, {637609, 637609}, {637612, 637612},
		},
	},
	{
		NiceType: "Hipercard",
		Type:     "hipercard",
		Patterns: []StructCardBrandPattern{
			{606282, 606282},
		},
	},
}

func GetCreditCardBrand(cardNumber string) string {
	// remove spaces
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")

	for _, card := range cardTypes {
		for _, pattern := range card.Patterns {
			if matchPrefix(cardNumber, pattern) {
				return strings.ToUpper(card.NiceType)
			}
		}
	}

	return ""
}

func matchPrefix(cardNumber string, p StructCardBrandPattern) bool {
	startStr := strconv.Itoa(p.Start)
	prefixLen := len(startStr)

	if len(cardNumber) < prefixLen {
		return false
	}

	value, err := strconv.Atoi(cardNumber[:prefixLen])
	if err != nil {
		return false
	}

	return value >= p.Start && value <= p.End
}
