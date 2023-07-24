
Notice  this behavior with the regex parser:
once it starts to match it can not go back 
therefore parsing "4.567" will never succeed here:
    o Input: "/4/ | '4.567'"


