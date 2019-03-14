var a = 5;

var makeAdder = fun(x) {
    return fun (y) {
        return x + y;
    };
};

var addTwo = makeAdder(2);

addTwo(a); // Output: 7

/*
    1. a == 5
    2. addTwo adds 2
    3. addTwo(a) => 5 + 2 == 7
*/
