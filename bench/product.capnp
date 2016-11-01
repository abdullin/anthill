using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("main");
$Go.import("main");

struct Classification {
    id @0 :UInt64;
    name @1 :Text;
}

struct Product {
    id @0 :UInt64;
    code @1 :Text;
    sku @2 :Text;
    description @3 :Text;
    classification @4: Classification;
}
