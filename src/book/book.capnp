using Go = import "/go.capnp";

@0xe445b2605be6742b;

$Go.package("book");
$Go.import("src/book");

struct Book {
  title @0 :Text;
  pages @1 :UInt32;
}
