const chooseBigger = fun(x,y) {
    print("choosing bigger betweeen", x, "and", y);

    if (x > y) {
        return x;
    }

    return y;
};

const max = fun(arr) {
    const findMax = fun(arr, max) {
        if (len(arr) == 0) {
            return max;
        }

        return findMax(rest(arr), chooseBigger(first(arr), max));
    };

    return findMax(arr, -99999999);
};


print(max([1,2,43,5,21,121]));
