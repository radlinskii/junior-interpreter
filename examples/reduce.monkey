const reduce = fun(arr, initial, fn) {
    const iter = fun(arr, result) {
        if(len(arr) == 0) {
            return result;
        }

        return iter(rest(arr), fn(result, first(arr)));
    };

    return iter(arr, initial);
};

const sum = fun(arr) {
    return reduce(arr, 0, fun(initial, el) {
        return initial + el;
    });
};

print(sum([1,2,3,4,5]));
