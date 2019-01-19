# Complexity Analysis

This document presents the time/space complexities for each of the algorithms
provided in the library.

## Table of Contents

- [Complexity Analysis](#complexity-analysis)
  - [Table of Contents](#table-of-contents)
  - [Statistics](#statistics)
    - [Quantile](#quantile)
      - [OSTQuantile/OSTMedian](#ostquantileostmedian)
      - [HeapMedian](#heapmedian)
    - [Min/Max](#minmax)
      - [Min](#min)
      - [Max](#max)
    - [Moment-Based Statistics](#moment-based-statistics)
      - [Mean](#mean)
      - [Moment](#moment)
      - [Variance](#variance)
      - [Std](#std)
      - [Skewness](#skewness)
      - [Kurtosis](#kurtosis)
      - [Core (Univariate)](#core-univariate)
    - [Joint Distribution Statistics](#joint-distribution-statistics)
      - [Covariance](#covariance)
      - [Correlation](#correlation)
      - [Autocorrelation](#autocorrelation)
      - [Core (Multivariate)](#core-multivariate)
  - [References](#references)

## Statistics

### [Quantile](https://godoc.org/github.com/alexander-yu/stream/quantile)

#### OSTQuantile/OSTMedian

Let `n` be the size of the window, or the stream if tracking the global quantile. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |

#### HeapMedian

Let `n` be the size of the window, or the stream if tracking the global quantile. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |

### [Min/Max](https://godoc.org/github.com/alexander-yu/stream/minmax)

#### Min

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time)        | Value (time)       | Space                         |
| :----------------: | :----------------: | :---------------------------: |
| `O(1)` (amortized) | `O(1)` (amortized) | `O(1)` if global, else `O(n)` |

#### Max

Let `n` be the size of the window, or the stream if tracking the global maximum. Then we have the following complexities:

| Push (time)        | Value (time)       | Space                         |
| :----------------: | :----------------: | :---------------------------: |
| `O(1)` (amortized) | `O(1)` (amortized) | `O(1)` if global, else `O(n)` |

### [Moment-Based Statistics](https://godoc.org/github.com/alexander-yu/stream/moment)

#### Mean

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Moment

Let `n` be the size of the window, or the stream if tracking the global minimum; let `k` be the moment being tracked. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(k^2)`    | `O(1)`       | `O(1)` if global, else `O(n)` |

See [Core](#Core) for an explanation of why `Push` has a time complexity of `O(k^2)`, rather than `O(k)`.

#### Variance

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Std

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Skewness

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Kurtosis

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Core (Univariate)

Let `n` be the size of the window, or the stream if tracking the global minimum; let `k` be the maximum exponent of the power sums that is being tracked. Then we have the following complexities:

| Push (time) | Sum (time) | Count (time) | Space                                     |
| :---------: | :--------: | :----------: | :---------------------------------------: |
| `O(k^2)`    | `O(1)`     | `O(1)`       | `O(k)` if global, else `O(k + n)` |

The reason that the `Push` method has a time complexity of `O(k^2)` is due to the algorithm being used to update the power sums; while the traditional `O(k)` method involves simply keeping track of raw power sums (i.e. non-centralized) and then representing the centralized power sum as a linear combination of the raw power sums and the mean (by doing binomial expansion), this is prone to underflow/overflow and as a result is much less numerically stable. See [1] for the paper whose algorithm this library uses, and a more in-depth explanation of the above.

### [Joint Distribution Statistics](https://godoc.org/github.com/alexander-yu/stream/joint)

#### Covariance

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Correlation

Let `n` be the size of the window, or the stream if tracking the global minimum. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Autocorrelation

Let `n` be the size of the window, or the stream if tracking the global minimum; let `l` be the lag of the autocorrelation. Then we have the following complexities:

| Push (time) | Value (time) | Space                             |
| :---------: | :----------: | :-------------------------------: |
| `O(1)`      | `O(1)`       | `O(l)` if global, else `O(l + n)` |

#### Core (Multivariate)

Let `n` be the size of the window, or the stream if tracking the global minimum. Moreover, let `t` be the number of tuples that are configured, let `d` be the number of variables being tracked. Now for a given tuple `m`, define

    p(m) = (m_1 + 1) * ... * (m_k + 1)

 and let `a` be the maximum such `p(m)` over all tuples `m` that are configured. Then we have the following complexities:

| Push (time)  | Sum (time) | Count (time) | Space                                                 |
| :----------: | :--------: | :----------: | :---------------------------------------------------: |
| `O(tda^2)` | `O(d)`     | `O(1)`       | `O(d + ta^2)` if global, else `O(d + ta^2 + n)` |

## References

1: P. Pebay, T. B. Terriberry, H. Kolla, J. Bennett, Numerically stable, scalable formulas for parallel and online computation of higher-order multivariate central moments with arbitrary weights, Computational Statistics 31 (2016) 1305–1325.
