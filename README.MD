# ANTHILL

This is a small research project in building HA event-driven business
apps for retail with specific technologies and trade-offs.


## Trade-offs

1. **Throughput over latency** - as long as response latency stays
   within the limit. If we start exceeding the limit, then new
   requests will be bounced back.
2. **Consistency over availability** - for operations where the cost
   of inconsistency or staleness is big for the business
   (e.g. overpicking or overselling).
3. **Availability over consistency** - preferring eventual
   consistency, when it is part of the business or throughput is more
   important (e.g. overselling items which could be back-ordered
   quickly).
   
## Tech

Currently evaluated tech stack:

1. LMDB - embedded DB based on B-Tree and designed for read-heavy operations.
2. golang - simple language with good concurrency and performance,
   works well with code-generation (gofmt and fast compilation cycle).
3. Cap'n Proto/FlatBuffers - serialization format that avoids some
   memory allocations.
