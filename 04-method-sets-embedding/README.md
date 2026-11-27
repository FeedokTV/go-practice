# 04 - Method sets and embedding

When `Speak` is defined on `*Dog`, `d.Speak()` is automatically converted to `(&d).Speak()`. With an interface it doesn't work like that, because the interface takes the method set of `Dog`. `Dog`'s method set is empty here, but `*Dog`'s contains `Speak()`.
