"use strict";
function factorial(n) {
    let result = 1;
    for (let i = 1; i <= n; i++) {
        result *= i;
    }
    return result;
}
function f(n) {
    if (!Number.isInteger(n) || n < 0) {
        throw new Error("Input harus integer positif.");
    }
    const pembilang = factorial(n);
    const penyebut = Math.pow(2, n);
    return Math.ceil(pembilang / penyebut);
}
console.log(f(0));
console.log(f(1));
console.log(f(2));
console.log(f(3));
console.log(f(4));
console.log(f(5));
console.log(f(10));
