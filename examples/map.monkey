const map = fun(arr, fn) {
    const iter = fun(arr, accumulator) {
        if (len(arr) == 0) {
            return accumulator;
        }

        return iter(rest(arr), push(accumulator, fn(first(arr))));
    };

    return iter(arr, []);
};

const a = [1,2,3,4,5];
const triple = fun(x) { return x*3; };


map(a, triple);
