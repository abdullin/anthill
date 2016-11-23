using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("main");
$Go.import("main");

struct Classification {
    id @0 :UInt64;
    name @1 :Text;
}

struct Product {
    code @0 :Text;
    sku @1 :Text;
    description @2 :Text;
    classification @3: Classification;
}
