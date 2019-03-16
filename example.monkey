// recursive function
var factorial = fun(x) {
    if (x < 1) {
        return 1;
    }

    return factorial(x - 1) * x;
};

var fiveFactorial = factorial(5); // 5!

// (120 / 4) / (3 * 2) == 5
var five = fiveFactorial / 4 / (3 * 2);

// array with objects of different type
var array = [0, 1, 1 + 1, true, "fifth element of an array"];

array[five - 1];
