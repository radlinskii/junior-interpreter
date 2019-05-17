*G* = < *N*,*T*,*P*,*S* >

*T* = {`EOF`, `const`, `=`, `;`, `a`, `b`, ..., `z`, `A`, `B`, ..., `Z`, `true`, `false`, 
`0`, `1`, ..., `9`, `:`, `;`, `,`, `{`, `}`, `[`, `]`, `(`, `)`, `==`, `!=`,  `<=`,  `>=`,  `<`,
`?`,  `+`,  `/`, `"`, `if`, `else`, `return`, `fun`}


*N* = {
**Statements**, **Statement**, **Expression**, **ConstStatement**, **ExpressionStatement**, **BlockStatement**
**Identifier**, **Letters**, **Letter**, **IntegerLiteral**, **Digits**, **Digit**, **BooleanLiteral**,
**StringLiteral**, **PrefixExpression**, **OperatorPrefix**, **InfixExpression**, **OperatorInfix**, **BANG**,
**MINUS**, **EQ**, **NEQ**,**LTE**, **GTE**, **LT**, **GT**, **PLUS**, **SLASH**, **ASTERISK**, **IfStatement**,
**FunctionLiteral**, **Identifiers**, **ReturnStatement**, **CallExpression**, **Expressions**, **ArrayLiteral**,
**IndexExpression**, **HashLiteral**, **ExpressionPairs** 
}

*S* = ****Statements****

*P* = {  
&nbsp;&nbsp; **Statements** &rarr; `EOF` | **Statement** | **Statements**,  
&nbsp;&nbsp; **Statement** &rarr; **ConstStatement** | **ReturnStatement** | **BlockStatement** | **IfStatement** | **ExpressionStatement**,  
&nbsp;&nbsp; **ConstStatement** &rarr; `const` **Identifier** `=` **Expression**`;`,  
&nbsp;&nbsp; **ReturnStatement** &rarr; `return`&nbsp;`;` | `return` **Expression**`;`,  
&nbsp;&nbsp; **IfStatement** &rarr; `if`&nbsp;`(`**Expression**`)`&nbsp;`{`**BlockStatement**`}` |
`if`&nbsp;`(`**Expression**`)``{` **BlockStatement**`}`&nbsp;`else`&nbsp;`{` **BlockStatement** `}`,  
&nbsp;&nbsp; **BlockStatement** &rarr; **Statement**`;`**BlockStatement** | **Statement**`;`,  
&nbsp;&nbsp; **ExpressionStatement** &rarr; **Expression**`;`,  
&nbsp;&nbsp; **Expression** &rarr; **Identifier** | **IntegerLiteral** | **BooleanLiteral** | **StringLiteral** |
**PrefixExpression** | **FunctionLiteral** | **InfixExpression** | **CallExpression** | **ArrayLiteral** |
**IndexExpression** | **HashLiteral** | `(`**Expression**`)`,  
&nbsp;&nbsp; **Identifier** &rarr; **Letters**,  
&nbsp;&nbsp; **Letters** &rarr; **Letter** | **Letter****Letters**,  
&nbsp;&nbsp; **Letter** &rarr; `a` | `b` | .. | `z` | `A` | `B` | .. | `Z`,  
&nbsp;&nbsp; **IntegerLiteral** &rarr; **Digits**,  
&nbsp;&nbsp; **Digits** &rarr; **Digit** | **Digit****Digits**,  
&nbsp;&nbsp; **Digit** &rarr; `0` | `1` | .. | `9`,  
&nbsp;&nbsp; **BooleanLiteral** &rarr; `true` | `false`,  
&nbsp;&nbsp; **StringLiteral** &rarr; `"`**Letters**`"` | `""`,  
&nbsp;&nbsp; **PrefixExpression** &rarr; **OperatorPrefix** **Expression**,  
&nbsp;&nbsp; **OperatorPrefix** &rarr; **MINUS** | **BANG**,  
&nbsp;&nbsp; **InfixExpression** &rarr; **Expression** **OperatorInfix** **Expression**,  
&nbsp;&nbsp; **OperatorInfix** &rarr; **EQ** | **NEQ** | **LTE** | **GTE** | **LT** | **GT** | **PLUS** |**MINUS** |
**SLASH** | **ASTERISK**,  
&nbsp;&nbsp; **BANG** &rarr; `!`,  
&nbsp;&nbsp; **MINUS** &rarr; `-`,  
&nbsp;&nbsp; **EQ** &rarr; `==`,  
&nbsp;&nbsp; **NEQ**&rarr; `!=`,  
&nbsp;&nbsp; **LTE** &rarr; `<=`,  
&nbsp;&nbsp; **GTE** &rarr; `>=`,  
&nbsp;&nbsp; **LT** &rarr; `<`,  
&nbsp;&nbsp; **GT** &rarr; `?`,  
&nbsp;&nbsp; **PLUS** &rarr; `+`,  
&nbsp;&nbsp; **SLASH** &rarr; `/`,  
&nbsp;&nbsp; **ASTERISK** &rarr; `*`,  
&nbsp;&nbsp; **FunctionLiteral** &rarr; `fun`&nbsp;`(`**Identifiers**`)`&nbsp;`{`**BlockStatement**&nbsp;**ReturnStatement**`}` |
`fun`&nbsp;`()`&nbsp;`{`**BlockStatement**&nbsp;**ReturnStatement**`}`,  
&nbsp;&nbsp; **Identifiers** &rarr; **Identifier** | **Identifier**`,`**Identifiers**,  
&nbsp;&nbsp; **CallExpression** &rarr; **Identifier**`()` | **Identifier**`(`**Expressions**`)`,  
&nbsp;&nbsp; **Expressions** &rarr; **Expression** | **Expression**`,`&nbsp;**Expressions**,  
&nbsp;&nbsp; **ArrayLiteral** &rarr; `[`**Expressions**`]`,  
&nbsp;&nbsp; **IndexExpression** &rarr; **Identifier**`[`**Expression**`]`,  
&nbsp;&nbsp; **HashLiteral** &rarr; `{`**ExpressionPairs**`}`,  
&nbsp;&nbsp; **ExpressionPairs** &rarr; **Expression**`:`&nbsp;**Expression** |
**Expression**`:`&nbsp;**Expression**`,`**ExpressionPairs**,  
}
