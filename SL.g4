grammar SL;

// Grammar Rules =============

start: e = expr EOF;

expr:
	e = conditionalOr (op = '?' e1 = conditionalOr ':' e2 = expr)?;

conditionalOr:
	e = conditionalAnd (ops += '||' e1 += conditionalAnd)*;

conditionalAnd: e = relation (ops += '&&' e1 += relation)*;

relation:
	calc
	| relation op = (
		'<'
		| '<='
		| '>='
		| '>'
		| '=='
		| '!='
		| 'in'
	) relation;

calc:
	unary
	| calc op = ('*' | '/' | '%') calc
	| calc op = ('+' | '-') calc;

unary:
	member					# MemberExpr
	| (ops += '!')+ member	# LogicalNot
	| (ops += '-')+ member	# Negate;

member:
	primary																# PrimaryExpr
	| member op = '.' (opt = '?')? id = escapeIdent						# Select
	| member op = '.' id = IDENTIFIER open = '(' args = exprList? ')'	# MemberCall
	| member op = '[' (opt = '?')? index = expr ']'						# Index;

primary:
	leadingDot = '.'? id = IDENTIFIER # Ident
	| leadingDot = '.'? id = IDENTIFIER (
		op = '(' args = exprList? ')'
	)													# GlobalCall
	| '(' e = expr ')'									# Nested
	| op = '[' elems = listInit? ','? ']'				# CreateList
	| op = '{' entries = mapInitializerList? ','? '}'	# CreateStruct
	| leadingDot = '.'? ids += IDENTIFIER (
		ops += '.' ids += IDENTIFIER
	)* op = '{' entries = fieldInitializerList? ','? '}'	# CreateMessage
	| literal												# ConstantLiteral;

exprList: e += expr (',' e += expr)*;

listInit: elems += optExpr (',' elems += optExpr)*;

fieldInitializerList:
	fields += optField cols += ':' values += expr (
		',' fields += optField cols += ':' values += expr
	)*;

optField: (opt = '?')? escapeIdent;

mapInitializerList:
	keys += optExpr cols += ':' values += expr (
		',' keys += optExpr cols += ':' values += expr
	)*;

escapeIdent:
	id = IDENTIFIER			# SimpleIdentifier
	| id = ESC_IDENTIFIER	# EscapedIdentifier;

optExpr: (opt = '?')? e = expr;

literal:
	sign = MINUS? tok = NUM_INT		# Int
	| tok = NUM_UINT				# Uint
	| sign = MINUS? tok = NUM_FLOAT	# Double
	| tok = STRING					# String
	| tok = BYTES					# Bytes
	| tok = TRUE					# BoolTrue
	| tok = FALSE					# BoolFalse
	| tok = NULL					# Null;

// Keywords
FALSE: 'false';
NULL: 'null';
TRUE: 'true';

// Not supported
AS: 'as';
BREAK: 'break';
CONST: 'const';
CONTINUE: 'continue';
ELSE: 'else';
FOR: 'for';
FUNCTION: 'function';
IF: 'if';
IMPORT: 'import';
IN: 'in';
LET: 'let';
LOOP: 'loop';
PACKAGE: 'package';
NAMESPACE: 'namespace';
RETURN: 'return';
VAR: 'var';
VOID: 'void';
WHILE: 'while';

// Lexer Rules ===========

EQUALS: '==';
NOT_EQUALS: '!=';
LESS: '<';
LESS_EQUALS: '<=';
GREATER_EQUALS: '>=';
GREATER: '>';
LOGICAL_AND: '&&';
LOGICAL_OR: '||';

LBRACKET: '[';
RPRACKET: ']';
LBRACE: '{';
RBRACE: '}';
LPAREN: '(';
RPAREN: ')';
DOT: '.';
COMMA: ',';
MINUS: '-';
EXCLAM: '!';
QUESTIONMARK: '?';
COLON: ':';
PLUS: '+';
STAR: '*';
SLASH: '/';
PERCENT: '%';

fragment BACKSLASH: '\\';
fragment LETTER: 'A' ..'Z' | 'a' ..'z';
fragment DIGIT: '0' ..'9';
fragment EXPONENT: ('e' | 'E') ( '+' | '-')? DIGIT+;
fragment HEXDIGIT: ('0' ..'9' | 'a' ..'f' | 'A' ..'F');
fragment RAW: 'r' | 'R';

fragment ESC_SEQ:
	ESC_CHAR_SEQ
	| ESC_BYTE_SEQ
	| ESC_UNI_SEQ
	| ESC_OCT_SEQ;

fragment ESC_CHAR_SEQ:
	BACKSLASH (
		'a'
		| 'b'
		| 'f'
		| 'n'
		| 'r'
		| 't'
		| 'v'
		| '"'
		| '\''
		| '\\'
		| '?'
		| '`'
	);

fragment ESC_OCT_SEQ:
	BACKSLASH ('0' ..'3') ('0' ..'7') ('0' ..'7');

fragment ESC_BYTE_SEQ: BACKSLASH ( 'x' | 'X') HEXDIGIT HEXDIGIT;

fragment ESC_UNI_SEQ:
	BACKSLASH 'u' HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT
	| BACKSLASH 'U' HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT;

WHITESPACE: ('\t' | ' ' | '\r' | '\n' | '\u000C')+ -> channel(HIDDEN);
COMMENT: '//' (~'\n')* -> channel(HIDDEN);

NUM_FLOAT: (
		DIGIT+ ('.' DIGIT+) EXPONENT?
		| DIGIT+ EXPONENT
		| '.' DIGIT+ EXPONENT?
	);

NUM_INT: ( DIGIT+ | '0x' HEXDIGIT+);

NUM_UINT: DIGIT+ ( 'u' | 'U') | '0x' HEXDIGIT+ ( 'u' | 'U');

STRING:
	'"' (ESC_SEQ | ~('\\' | '"' | '\n' | '\r'))* '"'
	| '\'' (ESC_SEQ | ~('\\' | '\'' | '\n' | '\r'))* '\''
	| '"""' (ESC_SEQ | ~('\\'))*? '"""'
	| '\'\'\'' (ESC_SEQ | ~('\\'))*? '\'\'\''
	| RAW '"' ~('"' | '\n' | '\r')* '"'
	| RAW '\'' ~('\'' | '\n' | '\r')* '\''
	| RAW '"""' .*? '"""'
	| RAW '\'\'\'' .*? '\'\'\'';

BYTES: ('b' | 'B') STRING;

IDENTIFIER: (LETTER | '_') ( LETTER | DIGIT | '_')*;
ESC_IDENTIFIER:
	'`' (LETTER | DIGIT | '_' | '.' | '-' | '/' | ' ')+ '`';