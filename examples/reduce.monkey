var reduce = fun(arr, initial, fn) {
    var iter = fun(arr, result) {
        return if(len(arr) == 0) {
            return result
        } else {
            return iter(rest(arr), fn(result, first(arr)));
        }
    };

    return iter(arr, initial);
};

var sum = fun(arr) {
    return reduce(arr, 0, fun(initial, el) {
        return initial + el
    });
};

sum([1,2,3,4,5]);
