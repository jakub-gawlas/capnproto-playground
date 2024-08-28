using Go = import "/go.capnp";

@0xc8e5b86f8a9fb4d0;

$Go.package("arith");
$Go.import("src/arith");

interface Arith {
    multiply @0 (a: Int64, b: Int64) -> (product: Int64);
    divide @1 (num: Int64, denom: Int64) -> (quotient: Int64, remainder: Int64);
}

struct Number {
    value @0 :Int64;

    interface Ops {
        plus @0 (valueToAdd: Int64) -> (sum: Number);
    }
}
