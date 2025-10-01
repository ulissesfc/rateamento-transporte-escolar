var x1 >= 0, <= 3, integer;
var x2 >= 0, <= 1, integer;
var x3 >= 0, <= 2, integer;
var x4 >= 0, <= 4, integer;
var x5 >= 0, <= 2, integer;
var x6 >= 0, <= 1, integer;
minimize obj: x1 + x2 + x3 + x4 + x5 + x6;
subject to demand: 16*x1 + 16*x2 + 17*x3 + 19*x4 + 22*x5 + 24*x6 >= 121;
s.t. x1_min: x1 >= 2;
