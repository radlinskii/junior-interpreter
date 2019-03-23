const min = fun(arr) {
    const findMin = fun(arr, min) {
        if (len(arr) == 0) {
            return min;
        }

        if (first(arr) < min) {
            return findMin(rest(arr), first(arr));
        }

        return findMin(rest(arr), min);
    };

    return findMin(arr, 99999999);
};


print(min([1,-2,13,-1,4]));
