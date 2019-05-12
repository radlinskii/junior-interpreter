### Metasymbols
- d - digit (0-9)
- l - letter (a-zA-Z)
- character - any character
- x|y - x or y
- {x} - 0 or more occurences of x
- 'x' - literal interpretation 

### Tokens

| Lp | token | description |
| :---: | :---: | :---: |
| 1	| INT | d{d} |
| 2	| STRING | "{character}" |
| 3	| BOOLEAN | 'true'\|'false' |
| 4	| ASSIGN | '=' |
| 5	| PLUS | '+' |
| 6	| MINUS | '-' |
| 7	| BANG | '!' |
| 8	| ASTERISK | '*' |
| 9	| SLASH | '/' |
| 10	| LT | '<' |
| 11	| GT | '>' |
| 12	| LTE | '<=' |
| 13	| GTE | '>=' |
| 14	| EQ | '==' |
| 15	| NEQ | '!=' |
| 16	| COMMA | ',' |
| 17	| SEMICOLON | ';' |
| 18	| COLON | ':' |
| 19	| LPAREN | '(' |
| 20	| RPAREN | ')' |
| 21	| LBRACE | '{' |
| 22	| RBRACE | '}' |
| 23	| LBRACKET | '[' |
| 24	| RBRACKET | ']' |
| 25	| FUNCTION | 'fun' |
| 26	| RETURN | 'return' |
| 27	| CONST | 'const' |
| 28	| IF | 'if' |
| 29	| ELSE | 'else' |
| 30	| ILLEGAL |  |
| 31	| EOF | 'EOF' |
| 32	| IDENT | l{l} |
