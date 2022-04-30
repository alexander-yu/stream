# Complexity Analysis

This document presents the time/space complexities for each of the algorithms
provided in the library.

## Table of Contents

- [Complexity Analysis](#complexity-analysis)
  - [Table of Contents](#table-of-contents)
  - [Statistics](#statistics)
    - [Quantile](#quantile)
      - [Quantile](#quantile-1)
      - [Median](#median)
      - [IQR](#iqr)
      - [HeapMedian](#heapmedian)
    - [Min/Max](#minmax)
      - [Min](#min)
      - [Max](#max)
    - [Moment-Based Statistics](#moment-based-statistics)
      - [Mean](#mean)
      - [EWMA](#ewma)
      - [Moment](#moment)
      - [EWMMoment](#ewmmoment)
      - [Std](#std)
      - [EWMStd](#ewmstd)
      - [Skewness](#skewness)
      - [Kurtosis](#kurtosis)
      - [Core (Univariate)](#core-univariate)
    - [Joint Distribution Statistics](#joint-distribution-statistics)
      - [Cov](#cov)
      - [EWMCov](#ewmcov)
      - [Corr](#corr)
      - [EWMCorr](#ewmcorr)
      - [Autocorr](#autocorr)
      - [Autocov](#autocov)
      - [Core (Multivariate)](#core-multivariate)
  - [References](#references)

## Statistics

### [Quantile](https://godoc.org/github.com/alexander-yu/stream/quantile)

#### Quantile

Let `n` be the size of the window, or the stream if tracking the global quantile. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |

#### Median

Let `n` be the size of the window, or the stream if tracking the global median. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |

#### IQR

Let `n` be the size of the window, or the stream if tracking the global interquartile range. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(log n)`   | `O(n)` |

#### HeapMedian

Let `n` be the size of the window, or the stream if tracking the global median. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(log n)`  | `O(1)`       | `O(n)` |

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

Let `n` be the size of the window, or the stream if tracking the global mean. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### EWMA

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Moment

Let `n` be the size of the window, or the stream if tracking the global moment; let `k` be the moment being tracked. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(k^2)`    | `O(1)`       | `O(1)` if global, else `O(n)` |

See [Core](#Core) for an explanation of why `Push` has a time complexity of `O(k^2)`, rather than `O(k)`.

#### EWMMoment

Let `k` be the moment being tracked. Then we have the following complexities:

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(k^2)`    | `O(1)`       | `O(1)` |

See [Core](#Core) for an explanation of why `Push` has a time complexity of `O(k^2)`, rather than `O(k)`.

#### Std

Let `n` be the size of the window, or the stream if tracking the global standard deviation. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### EWMStd

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Skewness

Let `n` be the size of the window, or the stream if tracking the global skewness. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Kurtosis

Let `n` be the size of the window, or the stream if tracking the global kurtosis. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### Core (Univariate)

Let `n` be the size of the window, or the stream if tracking the global sums; let `k` be the maximum exponent of the power sums that is being tracked. Then we have the following complexities:

| Push (time) | Sum (time) | Count (time) | Space                                     |
| :---------: | :--------: | :----------: | :---------------------------------------: |
| `O(k^2)`    | `O(1)`     | `O(1)`       | `O(k)` if global, else `O(k + n)` |

The reason that the `Push` method has a time complexity of `O(k^2)` is due to the algorithm being used to update the power sums; while the traditional `O(k)` method involves simply keeping track of raw power sums (i.e. non-centralized) and then representing the centralized power sum as a linear combination of the raw power sums and the mean (by doing binomial expansion), this is prone to underflow/overflow and as a result is much less numerically stable. See [1] for the paper whose algorithm this library uses, and a more in-depth explanation of the above.

### [Joint Distribution Statistics](https://godoc.org/github.com/alexander-yu/stream/joint)

#### Cov

Let `n` be the size of the window, or the stream if tracking the global covariance. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### EWMCov

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Corr

Let `n` be the size of the window, or the stream if tracking the global correlation. Then we have the following complexities:

| Push (time) | Value (time) | Space                         |
| :---------: | :----------: | :---------------------------: |
| `O(1)`      | `O(1)`       | `O(1)` if global, else `O(n)` |

#### EWMCorr

| Push (time) | Value (time) | Space  |
| :---------: | :----------: | :----: |
| `O(1)`      | `O(1)`       | `O(1)` |

#### Autocorr

Let `n` be the size of the window, or the stream if tracking the global autocorrelation; let `l` be the lag of the autocorrelation. Then we have the following complexities:

| Push (time) | Value (time) | Space                             |
| :---------: | :----------: | :-------------------------------: |
| `O(1)`      | `O(1)`       | `O(l)` if global, else `O(l + n)` |

#### Autocov

Let `n` be the size of the window, or the stream if tracking the global autocovariance; let `l` be the lag of the autocovariance. Then we have the following complexities:

| Push (time) | Value (time) | Space                             |
| :---------: | :----------: | :-------------------------------: |
| `O(1)`      | `O(1)`       | `O(l)` if global, else `O(l + n)` |

#### Core (Multivariate)

Let `n` be the size of the window, or the stream if tracking the global sums. Moreover, let `t` be the number of tuples that are configured, let `d` be the number of variables being tracked. Now for a given tuple `m`, define

    p(m) = (m_1 + 1) * ... * (m_k + 1)

 and let `a` be the maximum such `p(m)` over all tuples `m` that are configured. Then we have the following complexities:

| Push (time)  | Sum (time) | Count (time) | Space                                                 |
| :----------: | :--------: | :----------: | :---------------------------------------------------: |
| `O(tda^2)` | `O(d)`     | `O(1)`       | `O(d + ta^2)` if global, else `O(d + ta^2 + n)` |

## References

1: P. Pebay, T. B. Terriberry, H. Kolla, J. Bennett, Numerically stable, scalable formulas for parallel and online computation of higher-order multivariate central moments with arbitrary weights, Computational Statistics 31 (2016) 1305–1325.
