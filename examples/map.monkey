var map = fun(arr, fn) {
    var iter = fun(arr, accumulator) {
        if (len(arr) == 0) {
            return accumulator;
        }

        return iter(rest(arr), push(accumulator, fn(first(arr))));
    };

    return iter(arr, [])
};

var a = [1,2,3,4,5];
var triple = fun(x) { return x*3; };

map(a, triple)
