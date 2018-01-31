## Benchmarks for various decimal programs.

## General notes

- Times are measured in seconds, unless otherwise noted
- Measured on a MacBook Pro, 2.9 GHz Intel Core i5, 8 GB 2133 MHz LPDDR3
- Some benchmarks are adapted from www.bytereef.org/mpdecimal/benchmarks.html

The benchmarks aim to be as fair as possible, but ultimately they do compare
different libraries with different feature sets. For example, [go-inf/inf][8]
boasts the fastest overall runtime of any library, but it also doesn't implement
the GDA spec: it lacks contexts, non-finite NaN/Inf/±zero values, conditions,
etc. Further, programs like [cocroachdb/apd][2] sacrifice speed to ensure strict
compatibility with the GDA spec.

In general, libraries that cannot fully complete a challenge will be unranked.
For example, Go's `float64` type cannot provide 19 or more digits of precision,
so it's unranked in the Pi test. Similarly so with [apmckinlay/dnum][3].

## Pi

|    Program (version)                      | 9 digits | 19 digits | 38 digits | 100 digits | average |
|-------------------------------------------|----------|-----------|-----------|------------|---------|
| [go-inf/inf][8] (Go 1.9)                       | 0.10     | 0.23      | 0.53      | 1.43       | 0.572   |
| [JDK BigDecimal][4] (Java 1.8, warm)           | 0.049    | 0.19      | 0.60      | 3.29       | 1.05    |
| [ericlagergren/decimal][1] (Go 1.9, mode Go)   | 0.046    | 0.29      | 0.81      | 3.09       | 1.10    |
| [ericlagergren/decimal][1] (Go 1.9, mode GDA)  | 0.048    | 0.34      | 0.97      | 3.70       | 1.26    |
| [Python decimal][5] (Python 3.6.2)             | 0.27     | 0.58      | 1.32      | 4.52       | 1.67    |
| [JDK BigDecimal][4] (Java 1.8)                 | 0.29     | 0.96      | 1.79      | 3.99       | 1.76    |
| [shopspring/decimal][7] decimal (Go 1.9)       | 0.38     | 0.94      | 1.95      | 5.26       | 2.13    |
| [cockroachdb/apd][2] (Go 1.9)                  | 0.52     | 2.14      | 9.01      | 71.62      | 20.81   |
| [Python decimal][6] (Python 2.7.10)            | 12.93    | 28.91     | 64.96     | 192.58     | 74.84   |
| float64 (Go 1.9)                          | 0.057    | -         | -         | -          | -       |
| double (C LLVM 9.0.0 -O3)                 | 0.057    | -         | -         | -          | -       |
| [apmckinlay/dnum][3] (Go 1.9)                  | 0.091    | -         | -         | -          | -       |
| float (Python 2.7.10)                     | 0.59     | -         | -         | -          | -       |

## Mandelbrot

|    Program (version)                      | 9 digits | 16 digits | 19 digits | 34 digits | 38 digits | average |
|-------------------------------------------|----------|-----------|-----------|-----------|-----------|---------|
| [ericlagergren/decimal][1] (Go 1.9, mode GDA)  | 2.73     | 9.07      | 14.54     | 24.95     | 25.09     | 15.27   |
| [ericlagergren/decimal][1] (Go 1.9, mode Go)   | 2.73     | 9.70      | 15.02     | 26.13     | 26.62     | 16.04   |
| float64 (Go 1.9)                          | 0.0034   | -         | -         | -         | -         | -       |

[1]: https://github.com/ericlagergren/decimal
[2]: https://github.com/cockroachdb/apd
[3]: https://github.com/apmckinlay/gsuneido/util/dnum
[4]: https://docs.oracle.com/javase/8/docs/api/java/math/BigDecimal.html
[5]: https://docs.python.org/3.6/library/decimal.html
[6]: https://docs.python.org/2/library/decimal.html
[7]: https://github.com/shopspring/decimal
[8]: https://github.com/go-inf/inf
