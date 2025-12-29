# High Water Mark (HWM)
25.11.2025

The High Water Mark (HWM) mechanism describes how KSeF manages data completeness over time for the `PermanentStorage` date.

At any given moment, the system knows a point in time (`HWM`) up to which it is certain that all invoices have been saved and no new documents with a `PermanentStorage` date earlier than or equal to this moment will appear.

![HWM](hwm.png)

- For time <= `HWM` - all invoices with a `PermanentStorage` date in this range have already been permanently saved in KSeF.
The system guarantees that no new invoice with a `PermanentStorage` date <= `HWM` will appear in the future.
- In the range (`HWM`, `Now`):
    - some invoices are already visible and can be returned in the query,
    - due to the asynchronous and multithreaded nature of the saving process, new invoices may still appear in this range, i.e., with a `PermanentStorage` date falling within the range (`HWM`, `Now`].

Conclusion:
- everything that is <= `HWM` can be treated as a **closed** and **complete** set,
- everything that is > `HWM` is **potentially incomplete** and requires careful handling during synchronization.

## Scenario 1 - synchronization "only up to HWM"

![HWM-1](hwm-1.png)

With each query, the system retrieves invoices **from the "last known point" only up to the current `HWM` value**. The new `HWM` value becomes the start of the next range.

Advantages:
- data up to `HWM` is definitive - there is no need to recheck the same range,
- the number of duplicates between consecutive downloads is minimal.

Consequences:
- some of the most recent invoices from the range `(HWM, Now]` are not visible in the local system - they will appear only after `HWM` shifts in the next cycle.

This scenario is recommended for incremental, automatic data synchronization where optimizing traffic and minimizing the number of duplicates is more important than immediate availability of the most recent invoices.

## Scenario 2 - synchronization "up to Now"

![HWM-2](hwm-2.png)

The system integrating with KSeF performs cyclic, incremental queries **from the last starting point up to `Now`** and saves all returned invoices, including those from the range `(HWM, Now]`.

Since data in this range may be incomplete, **the next query repeats part of the range** - at least from the previous `HWM` to the new `Now`. Deduplication is required on the local system side (e.g., by KSeF number).

Advantages:
- the local system (and user) sees the latest invoices as quickly as possible, without waiting for `HWM` to "catch up".

Consequences:
- the range `(HWM, Now]` must be checked again in the next query,
- duplicates will appear that need to be removed on the local system side.

The same mechanism can also be used **ad hoc**, when the user manually requests a data refresh - the system then retrieves "here and now" the most recent available invoices from the last known date up to `Now`.

## Related documents

- [Incremental invoice retrieval](przyrostowe-pobieranie-faktur.md)
