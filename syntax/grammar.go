package syntax

import (
	"github.com/zyedidia/gpeg/charset"
	p "github.com/zyedidia/gpeg/pattern"
)

// Pattern    <- Spacing_ (Grammar / Expression) EndOfFile_
// Grammar    <- Definition+
// Definition <- Identifier LEFTARROW Expression
//
// Expression <- Sequence (SLASH Sequence)*
// Sequence   <- Prefix*
// Prefix     <- (AND / NOT)? Suffix
// Suffix     <- Primary (QUESTION / STAR / PLUS)?
// Primary    <- CAP BRACEO Expression COMMA String BRACEC
//             / WORDS BRACEO (!BRACEC (String COMMA?))* BRACEC
//             / INCLUDE BRACEO String BRACEC
//             / REF BRACEO Expression COMMA String BRACEC
//             / BACK BRACEO String BRACEC
//             / Identifier !LEFTARROW
//             / '(' Expression ')'
//             / Literal / Class
//             / DOT
//
// Identifier <- IdentStart IdentCont* Spacing_
// IdentStart <- [a-zA-Z_]
// IdentCont  <- IdentStart / [0-9]
//
// Literal    <- ['] (!['] Char)* ['] Spacing_
// 			   / ["] (!["] Char)* ["] Spacing_
// Class      <- '[' CARAT? (!']' Range)* ']' Spacing_
// Range      <- Char '-' Char / Char
// Char       <- '\\' [nrt'"\[\]\\\-]
// 			/ '\\' [0-2][0-7][0-7]
// 			/ '\\' [0-7][0-7]?
// 			/ !'\\' .
//
// AND        <- '&' Spacing_
// NOT        <- '!' Spacing_
// QUESTION   <- '?' Spacing_
// STAR       <- '*' Spacing_
// PLUS       <- '+' Spacing_
// DOT        <- '.' Spacing_
// CARAT      <- '^' Spacing_
// BRACEO     <- '{' Spacing_
// BRACEC     <- '}' Spacing_
// LEFTARROW  <- '<-' Spacing_
// OPEN       <- '(' Spacing_
// CLOSE      <- ')' Spacing_
// SLASH      <- '/' Spacing_
// CAP        <- 'cap' Spacing_
// WORDS      <- 'words' Spacing_
// INCLUDE    <- 'include' Spacing_
// REF        <- 'ref' Spacing_
// BACK       <- 'back' Spacing_
// COMMA      <- ',' Spacing_
//
// Spacing_   <- (Space_ / Comment_)*
// Comment_   <- '#' (!EndOfLine_ .)* EndOfLine_
// Space_     <- ' ' / '\t' / EndOfLine_
// EndOfLine_ <- '\r\n' / '\n' / '\r'
// EndOfFile_ <- !.

const (
	idPattern = iota
	idGrammar
	idDefinition
	idExpression
	idSequence
	idPrefix
	idSuffix
	idPrimary
	idLiteral
	idRange
	idClass
	idIdentifier
	idIdentStart
	idIdentCont
	idChar
	idAND
	idNOT
	idQUESTION
	idSTAR
	idPLUS
	idDOT
	idCARAT
	idOPEN
	idCAP
	idWORDS
	idINCLUDE
	idREF
	idBACK
	idBRACEPO
)

var grammar = map[string]p.Pattern{
	"Pattern": p.Cap(p.Concat(
		p.NonTerm("Spacing"),
		p.Or(
			p.NonTerm("Grammar"),
			p.NonTerm("Expression"),
		),
		p.NonTerm("EndOfFile"),
	), idPattern),
	"Grammar": p.Cap(p.Plus(p.NonTerm("Definition")), idGrammar),
	"Definition": p.Cap(p.Concat(
		p.NonTerm("Identifier"),
		p.NonTerm("LEFTARROW"),
		p.NonTerm("Expression"),
	), idDefinition),

	"Expression": p.Cap(p.Concat(
		p.NonTerm("Sequence"),
		p.Star(p.Concat(
			p.NonTerm("SLASH"),
			p.NonTerm("Sequence"),
		)),
	), idExpression),
	"Sequence": p.Cap(p.Star(p.NonTerm("Prefix")), idSequence),
	"Prefix": p.Cap(p.Concat(
		p.Optional(p.Or(
			p.NonTerm("AND"),
			p.NonTerm("NOT"),
		)),
		p.NonTerm("Suffix"),
	), idPrefix),
	"Suffix": p.Cap(p.Concat(
		p.NonTerm("Primary"),
		p.Optional(p.Or(
			p.NonTerm("QUESTION"),
			p.NonTerm("STAR"),
			p.NonTerm("PLUS"),
		)),
	), idSuffix),
	"Primary": p.Cap(p.Or(
		p.Concat(
			p.NonTerm("CAP"),
			p.NonTerm("BRACEO"),
			p.NonTerm("Expression"),
			p.NonTerm("COMMA"),
			p.NonTerm("Literal"),
			p.NonTerm("BRACEC"),
		),
		p.Concat(
			p.NonTerm("WORDS"),
			p.NonTerm("BRACEO"),
			p.Star(p.Concat(
				p.Not(p.NonTerm("BRACEC")),
				p.NonTerm("Literal"),
				p.Optional(p.NonTerm("COMMA")),
			)),
			p.NonTerm("BRACEC"),
		),
		p.Concat(
			p.NonTerm("INCLUDE"),
			p.NonTerm("BRACEO"),
			p.NonTerm("Literal"),
			p.NonTerm("BRACEC"),
		),
		p.Concat(
			p.NonTerm("REF"),
			p.NonTerm("BRACEO"),
			p.NonTerm("Expression"),
			p.NonTerm("COMMA"),
			p.NonTerm("Literal"),
			p.NonTerm("BRACEC"),
		),
		p.Concat(
			p.NonTerm("BACK"),
			p.NonTerm("BRACEO"),
			p.NonTerm("Literal"),
			p.NonTerm("BRACEC"),
		),
		p.Concat(
			p.NonTerm("Identifier"),
			p.Not(p.NonTerm("LEFTARROW")),
		),
		p.Concat(
			p.NonTerm("OPEN"),
			p.NonTerm("Expression"),
			p.NonTerm("CLOSE"),
		),
		p.Concat(
			p.NonTerm("BRACEPO"),
			p.NonTerm("Expression"),
			p.NonTerm("BRACEPC"),
		),
		p.NonTerm("Literal"),
		p.NonTerm("Class"),
		p.NonTerm("DOT"),
	), idPrimary),

	"Identifier": p.Cap(p.Concat(
		p.NonTerm("IdentStart"),
		p.Star(p.NonTerm("IdentCont")),
		p.NonTerm("Spacing"),
	), idIdentifier),
	"IdentStart": p.Cap(
		p.Set(charset.Range('a', 'z').
			Add(charset.Range('A', 'Z')).
			Add(charset.New([]byte{'_'})),
		), idIdentStart),
	"IdentCont": p.Cap(p.Or(
		p.NonTerm("IdentStart"),
		p.Set(charset.Range('0', '9')),
	), idIdentCont),

	"Literal": p.Cap(p.Or(
		p.Concat(
			p.Literal("'"),
			p.Star(p.Concat(
				p.Not(p.Literal("'")),
				p.NonTerm("Char"),
			)),
			p.Literal("'"),
			p.NonTerm("Spacing"),
		),
		p.Concat(
			p.Literal("\""),
			p.Star(p.Concat(
				p.Not(p.Literal("\"")),
				p.NonTerm("Char"),
			)),
			p.Literal("\""),
			p.NonTerm("Spacing"),
		),
	), idLiteral),
	"Class": p.Cap(p.Concat(
		p.Literal("["),
		p.Optional(p.NonTerm("CARAT")),
		p.Star(p.Concat(
			p.Not(p.Literal("]")),
			p.NonTerm("Range"),
		)),
		p.Literal("]"),
		p.NonTerm("Spacing"),
	), idClass),
	"Range": p.Cap(p.Or(
		p.Concat(
			p.NonTerm("Char"),
			p.Literal("-"),
			p.NonTerm("Char"),
		),
		p.NonTerm("Char"),
	), idRange),
	"Char": p.Cap(p.Or(
		p.Concat(
			p.Literal("\\"),
			p.Set(charset.New([]byte{'n', 'r', 't', '\'', '"', '[', ']', '\\', '-'})),
		),
		p.Concat(
			p.Literal("\\x"),
			p.Set(charset.Range('0', '9').
				Add(charset.Range('a', 'f')).
				Add(charset.Range('A', 'F')),
			),
			p.Set(charset.Range('0', '9').
				Add(charset.Range('a', 'f')).
				Add(charset.Range('A', 'F')),
			),
		),
		p.Concat(
			p.Literal("\\"),
			p.Set(charset.Range('0', '2')),
			p.Set(charset.Range('0', '7')),
			p.Set(charset.Range('0', '7')),
		),
		p.Concat(
			p.Literal("\\"),
			p.Set(charset.Range('0', '7')),
			p.Optional(p.Set(charset.Range('0', '7'))),
		),
		p.Concat(
			p.Not(p.Literal("\\")),
			p.Any(1),
		),
	), idChar),

	"AND": p.Cap(p.Concat(
		p.Literal("&"),
		p.NonTerm("Spacing"),
	), idAND),
	"NOT": p.Cap(p.Concat(
		p.Literal("!"),
		p.NonTerm("Spacing"),
	), idNOT),
	"QUESTION": p.Cap(p.Concat(
		p.Literal("?"),
		p.NonTerm("Spacing"),
	), idQUESTION),
	"STAR": p.Cap(p.Concat(
		p.Literal("*"),
		p.NonTerm("Spacing"),
	), idSTAR),
	"PLUS": p.Cap(p.Concat(
		p.Literal("+"),
		p.NonTerm("Spacing"),
	), idPLUS),
	"DOT": p.Cap(p.Concat(
		p.Literal("."),
		p.NonTerm("Spacing"),
	), idDOT),
	"CARAT": p.Cap(p.Concat(
		p.Literal("^"),
		p.NonTerm("Spacing"),
	), idCARAT),
	"OPEN": p.Cap(p.Concat(
		p.Literal("("),
		p.NonTerm("Spacing"),
	), idOPEN),
	"CLOSE": p.Concat(
		p.Literal(")"),
		p.NonTerm("Spacing"),
	),
	"BRACEO": p.Concat(
		p.Literal("{"),
		p.NonTerm("Spacing"),
	),
	"BRACEC": p.Concat(
		p.Literal("}"),
		p.NonTerm("Spacing"),
	),
	"BRACEPO": p.Cap(p.Concat(
		p.Literal("{{"),
		p.NonTerm("Spacing"),
	), idBRACEPO),
	"BRACEPC": p.Concat(
		p.Literal("}}"),
		p.NonTerm("Spacing"),
	),
	"SLASH": p.Concat(
		p.Literal("/"),
		p.NonTerm("Spacing"),
	),
	"LEFTARROW": p.Concat(
		p.Literal("<-"),
		p.NonTerm("Spacing"),
	),
	"CAP": p.Cap(p.Concat(
		p.Literal("cap"),
		p.NonTerm("Spacing"),
	), idCAP),
	"WORDS": p.Cap(p.Concat(
		p.Literal("words"),
		p.NonTerm("Spacing"),
	), idWORDS),
	"INCLUDE": p.Cap(p.Concat(
		p.Literal("include"),
		p.NonTerm("Spacing"),
	), idINCLUDE),
	"REF": p.Cap(p.Concat(
		p.Literal("ref"),
		p.NonTerm("Spacing"),
	), idREF),
	"BACK": p.Cap(p.Concat(
		p.Literal("back"),
		p.NonTerm("Spacing"),
	), idBACK),
	"COMMA": p.Concat(
		p.Literal(","),
		p.NonTerm("Spacing"),
	),

	"Spacing": p.Star(p.Or(
		p.NonTerm("Space"),
		p.NonTerm("Comment"),
	)),
	"Comment": p.Concat(
		p.Literal("#"),
		p.Star(p.Concat(
			p.Not(p.NonTerm("EndOfLine")),
			p.Any(1),
		)),
		p.NonTerm("EndOfLine"),
	),
	"Space": p.Or(
		p.Set(charset.New([]byte{' ', '\t'})),
		p.NonTerm("EndOfLine"),
	),
	"EndOfLine": p.Or(
		p.Literal("\r\n"),
		p.Literal("\n"),
		p.Literal("\r"),
	),
	"EndOfFile": p.Not(p.Any(1)),
}
