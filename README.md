# MatchingWildcardsInGo
UTF-8-ready and fast ASCII-only routines for matching wildcards

Matching Wildcards in Go

This file set includes ASCII and UTF-8-ready routines for matching wildcards in Go (fastwildcompare.go), based on the Rust+ implementation here: https://developforperformance.com/MatchingWildcardsInRust.html

It also includes Go implementations of ASCII testcases for correctness and performance originally implemented in C/C++, plus a new set of UTF-8 testcases originally implemented in Rust.

A description of the algorithm's implementation and testing strategies, performance findings, and thoughts about how to choose one routine over another appear here: https://developforperformance.com/MatchingWildcardsUTF8ReadyInGoSwiftAndCpp.html#MatchingWildcardsInGoAsciiVersion
