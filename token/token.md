### Metasymbols
- `d`- digit (0-9)
- `c` - character (a-zA-Z)
- x`|`y - x or y
- `{`x`}` - 0 or more occurrences of x
- x - literal interpretation 

### Tokens

| # | token | literal |
| :---: | :---: | :---: |
| 1	| *INT* | `d`{`d`} |
| 2	| *STRING* | `"`...`"` |
| 3	| *BOOLEAN* | `true` &#124; `false` |
| 4	| *ASSIGN* | `=` |
| 5	| *PLUS* | `+` |
| 6	| *MINUS* | `-` |
| 7	| *BANG* | `!` |
| 8	| *ASTERISK* | `*` |
| 9	| *SLASH* | `/` |
| 10	| *LT* | `<` |
| 11	| *GT* | `>` |
| 12	| *LTE* | `<=` |
| 13	| *GTE* | `>=` |
| 14	| *EQ* | `==` |
| 15	| *NEQ* | `!=` |
| 16	| *COMMA* | `,` |
| 17	| *SEMICOLON* | `;` |
| 18	| *COLON* | `:` |
| 19	| *LPAREN* | `(` |
| 20	| *RPAREN* | `)` |
| 21	| *LBRACE* | `{` |
| 22	| *RBRACE* | `}` |
| 23	| *LBRACKET* | `[` |
| 24	| *RBRACKET* | `]` |
| 25	| *IDENT* | `c`{`c`} |
| 26	| *FUNCTION* | `fun` |
| 27	| *RETURN* | `return` |
| 28	| *CONST* | `const` |
| 29	| *IF* | `if` |
| 30	| *ELSE* | `else` |
| 31	| *EOF* | `EOF` |
| 32	| *ILLEGAL* |  |
