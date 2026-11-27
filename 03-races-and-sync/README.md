# 03 - Races and sync

## Which I would ship

Atomic counter, because in the current implementation with a mutex we also burden the caller with the mutex: locking and unlocking it. Atomic is a more convenient way to do this kind of solution.
