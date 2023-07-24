
Notice  this behavior with the regex parser:
once it starts to match it can not go back 
therefore parsing "4.567" will never succeed here:
    o Input: "/4/ | '4.567'"

GPT:
In general, it is recommended to avoid using regular expressions (regex) within Parsing Expression Grammar (PEG) rules whenever possible. PEG is a parsing technique that differs from traditional regular expressions in its ability to handle more complex grammars and provide better error reporting.

While regex can be useful for simple pattern matching tasks, they have limitations when it comes to handling nested or recursive structures. PEG, on the other hand, allows you to define more expressive and context-sensitive grammars.

Using regex within PEG rules can lead to several issues:

1. Ambiguity: Regular expressions are inherently greedy, meaning they will match as much as possible. This can lead to ambiguity in parsing when used within PEG rules, as the PEG parser may not be able to determine the correct interpretation.

2. Readability and maintainability: Mixing regex with PEG rules can make the grammar harder to read and understand. PEG rules are typically more explicit and self-contained, making them easier to maintain and modify.

3. Error reporting: PEG parsers are designed to provide detailed error messages when parsing fails. When regex is used within PEG rules, error reporting becomes more challenging, as regex does not provide the same level of error information.

That being said, there may be cases where using regex within PEG rules is necessary or appropriate, especially for simple patterns or tokenization. However, it is generally advisable to leverage the full power of PEG by defining explicit rules instead of relying heavily on regex.

