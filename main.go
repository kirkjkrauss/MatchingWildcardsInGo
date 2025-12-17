// Go testcases for matching wildcards.
//
// Copyright 2025 Kirk J Krauss.  This is a Derivative Work based on
// material that is copyright 2018 IBM Corporation and available at
//
//	https://developforperformance.com/MatchingWildcardsInRust.html
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// This file provides sets of correctness and performance tests, for
// matching wildcards in Go, along with a main() routine that invokes the
// testcases and outputs the results.
package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Package-scope testcase selection flags.
//
// For a fair comparison involving implementations that aren't UTF-8-ready,
// set bTestUtf8 = false.
const (
	bComparePerformance   = false // Compares using ASCII tests
	bTestWild             = true
	bTestTame             = true
	bTesteEmpty           = true
	bTestUtf8             = true  // Skips ASCII test timings
	bTestCaseInsensitive  = true
)

// Package-scope variables for low-latency accumulation of performance data.
var (
	iAccumulatedTimeAscii int64
	iAccumulatedTimeUTF8  int64
	// Can add accumulator variables for more performance comparisons here...
	bTestingUtf8 bool
)

// This function compares a tame/wild string pair via each included routine.
func test(tame_string, wild_string string, bExpectedResult bool) bool {
	bPassed := true
	timeStart := time.Now()
	timeFinish := time.Now()

	if bComparePerformance {
		if !bTestingUtf8 {
			// Get execution times for our two matching wildcards routines.
			timeStart = time.Now()

			if bExpectedResult != FastWildCompareAscii(
				wild_string, tame_string) {
				bPassed = false
			}

			timeFinish = time.Now()
			iAccumulatedTimeAscii += timeFinish.Sub(timeStart).Nanoseconds()
		}

		timeStart = time.Now()

		// Allocate array-style memory and initialize with each input string's
		// 32-bit UTF-8 code points.
		//
		// A memory allocation failure can be associated with a panic.  In a
		// situation involving many calls to this routine, arrangements to
		// catch allocation failures may be placed around that entire set of
		// calls.
		//
		if bExpectedResult != FastWildCompareRuneSlices(
			[]rune(wild_string), []rune(tame_string)) {
			bPassed = false
		}

		timeFinish = time.Now()
		iAccumulatedTimeUTF8 += timeFinish.Sub(timeStart).Nanoseconds()

		// Can add more performance comparisons here...
	} else if bTestUtf8 && bCompareCaseInsensitive {
		// Case-insensitive matching:
		// Allocate array-style memory and initialize with each input string's
		// lowercased 32-bit UTF-8 code points.
		//
		// A memory allocation failure can be associated with a panic.  See
		// above comment regarding catching that situation in production code.
		//
		if bExpectedResult != FastWildCompareRuneSlices(
			[]rune(strings.ToLower(wild_string)),
			[]rune(strings.ToLower(tame_string))) {
			bPassed = false
		}
		// Can add tests for more matching wildcards routines here...
	} else if bExpectedResult != FastWildCompareAscii(
		wild_string, tame_string) {
		bPassed = false
	}

	return bPassed
}

// A set of wildcard comparison tests.
func testWild() {
	var iReps int
	bAllPassed := true
	bTestingUtf8 = false

	if bComparePerformance {
		// Can choose as many repetitions as you might expect in production.
		iReps = 1000000
	} else {
		iReps = 1
	}

	for iReps > 0 {
		iReps--

		// Case with first wildcard after total match.
		bAllPassed = bAllPassed && test("Hi", "Hi*", true)

		// Case with mismatch after '*'.
		bAllPassed = bAllPassed && test("abc", "ab*d", false)

		// Cases with repeating character sequences.
		bAllPassed = bAllPassed && test("abcccd", "*ccd", true)
		bAllPassed = bAllPassed && test("mississipissippi", "*issip*ss*", true)
		bAllPassed = bAllPassed && test("xxxx*zzzzzzzzy*f", "xxxx*zzy*fffff", false)
		bAllPassed = bAllPassed && test("xxxx*zzzzzzzzy*f", "xxx*zzy*f", true)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "xxxx*zzy*fffff", false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "xxxx*zzy*f", true)
		bAllPassed = bAllPassed && test("xyxyxyzyxyz", "xy*z*xyz", true)
		bAllPassed = bAllPassed && test("mississippi", "*sip*", true)
		bAllPassed = bAllPassed && test("xyxyxyxyz", "xy*xyz", true)
		bAllPassed = bAllPassed && test("mississippi", "mi*sip*", true)
		bAllPassed = bAllPassed && test("ababac", "*abac*", true)
		bAllPassed = bAllPassed && test("ababac", "*abac*", true)
		bAllPassed = bAllPassed && test("aaazz", "a*zz*", true)
		bAllPassed = bAllPassed && test("a12b12", "*12*23", false)
		bAllPassed = bAllPassed && test("a12b12", "a12b", false)
		bAllPassed = bAllPassed && test("a12b12", "*12*12*", true)

		if !bComparePerformance {
			// From DDJ reader Andy Belf: a case of repeating text matching 
			// the different kinds of wildcards in order of '*' and then '?'.
			bAllPassed = bAllPassed && test("caaab", "*a?b", true)

			// This similar case was found, probably independently, by Dogan 
			// Kurt.
			bAllPassed = bAllPassed && test("aaaaa", "*aa?", true)
		}

		// Additional cases where the '*' char appears in the tame string.
		bAllPassed = bAllPassed && test("*", "*", true)
		bAllPassed = bAllPassed && test("a*abab", "a*b", true)
		bAllPassed = bAllPassed && test("a*r", "a*", true)
		bAllPassed = bAllPassed && test("a*ar", "a*aar", false)

		// More double wildcard scenarios.
		bAllPassed = bAllPassed && test("XYXYXYZYXYz", "XY*Z*XYz", true)
		bAllPassed = bAllPassed && test("missisSIPpi", "*SIP*", true)
		bAllPassed = bAllPassed && test("mississipPI", "*issip*PI", true)
		bAllPassed = bAllPassed && test("xyxyxyxyz", "xy*xyz", true)
		bAllPassed = bAllPassed && test("miSsissippi", "mi*sip*", true)
		bAllPassed = bAllPassed && test("abAbac", "*Abac*", true)
		bAllPassed = bAllPassed && test("abAbac", "*Abac*", true)
		bAllPassed = bAllPassed && test("aAazz", "a*zz*", true)
		bAllPassed = bAllPassed && test("A12b12", "*12*23", false)
		bAllPassed = bAllPassed && test("a12B12", "*12*12*", true)
		bAllPassed = bAllPassed && test("oWn", "*oWn*", true)

		// Completely tame (no wildcards) cases.
		bAllPassed = bAllPassed && test("bLah", "bLah", true)

		// Simple mixed wildcard tests suggested by Marlin Deckert.
		bAllPassed = bAllPassed && test("a", "*?", true)
		bAllPassed = bAllPassed && test("ab", "*?", true)
		bAllPassed = bAllPassed && test("abc", "*?", true)

		// More mixed wildcard tests including coverage for false positives.
		bAllPassed = bAllPassed && test("a", "??", false)
		bAllPassed = bAllPassed && test("ab", "?*?", true)
		bAllPassed = bAllPassed && test("ab", "*?*?*", true)
		bAllPassed = bAllPassed && test("abc", "?**?*?", true)
		bAllPassed = bAllPassed && test("abc", "?**?*&?", false)
		bAllPassed = bAllPassed && test("abcd", "?b*??", true)
		bAllPassed = bAllPassed && test("abcd", "?a*??", false)
		bAllPassed = bAllPassed && test("abcd", "?**?c?", true)
		bAllPassed = bAllPassed && test("abcd", "?**?d?", false)
		bAllPassed = bAllPassed && test("abcde", "?*b*?*d*?", true)

		// Single-character-match cases.
		bAllPassed = bAllPassed && test("bLah", "bL?h", true)
		bAllPassed = bAllPassed && test("bLaaa", "bLa?", false)
		bAllPassed = bAllPassed && test("bLah", "bLa?", true)
		bAllPassed = bAllPassed && test("bLaH", "?Lah",
			bCompareCaseInsensitive)
		bAllPassed = bAllPassed && test("bLaH", "?LaH", true)

		// Many-wildcard scenarios.
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab",
			"a*a*a*a*a*a*aa*aaa*a*a*b", true)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"*a*b*ba*ca*a*aa*aaa*fa*ga*b*", true)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"*a*b*ba*ca*a*x*aaa*fa*ga*b*", false)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"*a*b*ba*ca*aaaa*fa*ga*gggg*b*", false)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"*a*b*ba*ca*aaaa*fa*ga*ggg*b*", true)
		bAllPassed = bAllPassed && test("aaabbaabbaab", "*aabbaa*a*", true)
		bAllPassed = bAllPassed && test("a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*",
			"a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*", true)
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaaa",
			"*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*", true)
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaa",
			"*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*", false)
		bAllPassed = bAllPassed && test("abc*abcd*abcde*abcdef*abcdefg*abcdefgh*abcdefghi*abcdefghij*abcdefghijk*abcdefghijkl*abcdefghijklm*abcdefghijklmn",
			"abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*a            bc*", false)
		bAllPassed = bAllPassed && test("abc*abcd*abcde*abcdef*abcdefg*abcdefgh*abcdefghi*abcdefghij*abcdefghijk*abcdefghijkl*abcdefghijklm*abcdefghijklmn",
			"abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*", true)
		bAllPassed = bAllPassed && test("abc*abcd*abcd*abc*abcd",
			"abc*abc*abc*abc*abc", false)
		bAllPassed = bAllPassed && test(
			"abc*abcd*abcd*abc*abcd*abcd*abc*abcd*abc*abc*abcd",
			"abc*abc*abc*abc*abc*abc*abc*abc*abc*abc*abcd", true)
		bAllPassed = bAllPassed && test("abc",
			"********a********b********c********", true)
		bAllPassed = bAllPassed && test("********a********b********c********",
			"abc", false)
		bAllPassed = bAllPassed && test("abc",
			"********a********b********b********", false)
		bAllPassed = bAllPassed && test("*abc*", "***a*b*c***", true)

		// Case-insensitive algorithm tests.
		if (bCompareCaseInsensitive) {
			bAllPassed = bAllPassed && test("mississippi", "*issip*PI",
				true)
			bAllPassed = bAllPassed && test("miSsissippi", "mi*Sip*",
				true)
			bAllPassed = bAllPassed && test("bLah", "bLaH",
				true)
		}

		// Tests suggested by other DDJ readers.
		bAllPassed = bAllPassed && test("", "?", false)
		bAllPassed = bAllPassed && test("", "*?", false)
		bAllPassed = bAllPassed && test("", "", true)
		bAllPassed = bAllPassed && test("a", "", false)
	}

	if bAllPassed {
		fmt.Println("Passed wildcard tests")
	} else {
		fmt.Println("Failed wildcard tests")
	}
}

// A set of tests with (almost) no '*' wildcards.
func testTame() {
	var iReps int
	bAllPassed := true
	bTestingUtf8 = false

	if bComparePerformance {
		// Can choose as many repetitions as you might expect in production.
		iReps = 1000000
	} else {
		iReps = 1
	}

	for iReps > 0 {
		iReps--

		// Case with last character mismatch.
		bAllPassed = bAllPassed && test("abc", "abd", false)

		// Cases with repeating character sequences.
		bAllPassed = bAllPassed && test("abcccd", "abcccd", true)
		bAllPassed = bAllPassed && test("mississipissippi",
			"mississipissippi", true)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf",
			"xxxxzzzzzzzzyfffff", false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "xxxxzzzzzzzzyf",
			true)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "xxxxzzy.fffff",
			false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "xxxxzzzzzzzzyf",
			true)
		bAllPassed = bAllPassed && test("xyxyxyzyxyz", "xyxyxyzyxyz", true)
		bAllPassed = bAllPassed && test("mississippi", "mississippi", true)
		bAllPassed = bAllPassed && test("xyxyxyxyz", "xyxyxyxyz", true)
		bAllPassed = bAllPassed && test("m ississippi", "m ississippi",
			true)
		bAllPassed = bAllPassed && test("ababac", "ababac?", false)
		bAllPassed = bAllPassed && test("dababac", "ababac", false)
		bAllPassed = bAllPassed && test("aaazz", "aaazz", true)
		bAllPassed = bAllPassed && test("a12b12", "1212", false)
		bAllPassed = bAllPassed && test("a12b12", "a12b", false)
		bAllPassed = bAllPassed && test("a12b12", "a12b12", true)

		// A mix of cases
		bAllPassed = bAllPassed && test("n", "n", true)
		bAllPassed = bAllPassed && test("aabab", "aabab", true)
		bAllPassed = bAllPassed && test("ar", "ar", true)
		bAllPassed = bAllPassed && test("aar", "aaar", false)
		bAllPassed = bAllPassed && test("XYXYXYZYXYz", "XYXYXYZYXYz", true)
		bAllPassed = bAllPassed && test("missisSIPpi", "missisSIPpi", true)
		bAllPassed = bAllPassed && test("mississipPI", "mississipPI", true)
		bAllPassed = bAllPassed && test("xyxyxyxyz", "xyxyxyxyz", true)
		bAllPassed = bAllPassed && test("miSsissippi", "miSsissippi",
			true)
			
		if bCompareCaseInsensitive {
			bAllPassed = bAllPassed && test("miSsissippi", "miSsisSippi",
				true)
			bAllPassed = bAllPassed && test("abAbac", "abAbac",
				true)
			bAllPassed = bAllPassed && test("abAbac", "abAbac",
				true)
			bAllPassed = bAllPassed && test("bLah", "bLaH",
				true)
		}

		bAllPassed = bAllPassed && test("aAazz", "aAazz", true)
		bAllPassed = bAllPassed && test("A12b12", "A12b123", false)
		bAllPassed = bAllPassed && test("a12B12", "a12B12", true)
		bAllPassed = bAllPassed && test("oWn", "oWn", true)
		bAllPassed = bAllPassed && test("bLah", "bLah", true)

		// Single '?' cases.
		bAllPassed = bAllPassed && test("a", "a", true)
		bAllPassed = bAllPassed && test("ab", "a?", true)
		bAllPassed = bAllPassed && test("abc", "ab?", true)

		// Mixed '?' cases.
		bAllPassed = bAllPassed && test("a", "??", false)
		bAllPassed = bAllPassed && test("ab", "??", true)
		bAllPassed = bAllPassed && test("abc", "???", true)
		bAllPassed = bAllPassed && test("abcd", "????", true)
		bAllPassed = bAllPassed && test("abc", "????", false)
		bAllPassed = bAllPassed && test("abcd", "?b??", true)
		bAllPassed = bAllPassed && test("abcd", "?a??", false)
		bAllPassed = bAllPassed && test("abcd", "??c?", true)
		bAllPassed = bAllPassed && test("abcd", "??d?", false)
		bAllPassed = bAllPassed && test("abcde", "?b?d*?", true)

		// Longer string scenarios.
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab", true)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab", true)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"abababababababababababababababababababaacacacacacacacadaeafagahaiajaxalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab", false)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaggggagaaaaaaaab", false)
		bAllPassed = bAllPassed && test("abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab",
			"abababababababababababababababababababaacacacacacacacadaeafagahaiajakalaaaaaaaaaaaaaaaaaffafagaagggagaaaaaaaab", true)
		bAllPassed = bAllPassed && test("aaabbaabbaab", "aaabbaabbaab",
			true)
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true)
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaaa",
			"aaaaaaaaaaaaaaaaa", true)
		bAllPassed = bAllPassed && test("aaaaaaaaaaaaaaaa",
			"aaaaaaaaaaaaaaaaa", false)
		bAllPassed = bAllPassed && test("abcabcdabcdeabcdefabcdefgabcdefghabcdefghiabcdefghijabcdefghijkabcdefghijklabcdefghijklmabcdefghijklmn",
			"abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabc",
			false)
		bAllPassed = bAllPassed && test("abcabcdabcdeabcdefabcdefgabcdefghabcdefghiabcdefghijabcdefghijkabcdefghijklabcdefghijklmabcdefghijklmn",
			"abcabcdabcdeabcdefabcdefgabcdefghabcdefghiabcdefghijabcdefghijkabcdefghijklabcdefghijklmabcdefghijklmn",
			true)
		bAllPassed = bAllPassed && test("abcabcdabcdabcabcd",
			"abcabc?abcabcabc", false)
		bAllPassed = bAllPassed && test(
			"abcabcdabcdabcabcdabcdabcabcdabcabcabcd",
			"abcabc?abc?abcabc?abc?abc?bc?abc?bc?bcd", true)
		bAllPassed = bAllPassed && test("?abc?", "?abc?", true)
	}

	if bAllPassed {
		fmt.Println("Passed tame string tests")
	} else {
		fmt.Println("Failed tame string tests")
	}
}

// A set of tests with empty input.
func testEmpty() {
	var iReps int
	bAllPassed := true
	bTestingUtf8 = false

	if bComparePerformance {
		// Can choose as many repetitions as you might expect in production.
		iReps = 1000000
	} else {
		iReps = 1
	}

	for iReps > 0 {
		iReps--

		// A simple case.
		bAllPassed = bAllPassed && test("", "abd", false)

		// Cases with repeating character sequences.
		bAllPassed = bAllPassed && test("", "abcccd", false)
		bAllPassed = bAllPassed && test("", "mississipissippi", false)
		bAllPassed = bAllPassed && test("", "xxxxzzzzzzzzyfffff", false)
		bAllPassed = bAllPassed && test("", "xxxxzzzzzzzzyf", false)
		bAllPassed = bAllPassed && test("", "xxxxzzy.fffff", false)
		bAllPassed = bAllPassed && test("", "xxxxzzzzzzzzyf", false)
		bAllPassed = bAllPassed && test("", "xyxyxyzyxyz", false)
		bAllPassed = bAllPassed && test("", "mississippi", false)
		bAllPassed = bAllPassed && test("", "xyxyxyxyz", false)
		bAllPassed = bAllPassed && test("", "m ississippi", false)
		bAllPassed = bAllPassed && test("", "ababac*", false)
		bAllPassed = bAllPassed && test("", "ababac", false)
		bAllPassed = bAllPassed && test("", "aaazz", false)
		bAllPassed = bAllPassed && test("", "1212", false)
		bAllPassed = bAllPassed && test("", "a12b", false)
		bAllPassed = bAllPassed && test("", "a12b12", false)

		// A mix of cases.
		bAllPassed = bAllPassed && test("", "n", false)
		bAllPassed = bAllPassed && test("", "aabab", false)
		bAllPassed = bAllPassed && test("", "ar", false)
		bAllPassed = bAllPassed && test("", "aaar", false)
		bAllPassed = bAllPassed && test("", "XYXYXYZYXYz", false)
		bAllPassed = bAllPassed && test("", "missisSIPpi", false)
		bAllPassed = bAllPassed && test("", "mississipPI", false)
		bAllPassed = bAllPassed && test("", "xyxyxyxyz", false)
		bAllPassed = bAllPassed && test("", "miSsissippi", false)
		bAllPassed = bAllPassed && test("", "miSsisSippi", false)
		bAllPassed = bAllPassed && test("", "abAbac", false)
		bAllPassed = bAllPassed && test("", "abAbac", false)
		bAllPassed = bAllPassed && test("", "aAazz", false)
		bAllPassed = bAllPassed && test("", "A12b123", false)
		bAllPassed = bAllPassed && test("", "a12B12", false)
		bAllPassed = bAllPassed && test("", "oWn", false)
		bAllPassed = bAllPassed && test("", "bLah", false)
		bAllPassed = bAllPassed && test("", "bLaH", false)

		// Both strings empty.
		bAllPassed = bAllPassed && test("", "", true)

		// Another simple case.
		bAllPassed = bAllPassed && test("abc", "", false)

		// More cases with repeating character sequences.
		bAllPassed = bAllPassed && test("abcccd", "", false)
		bAllPassed = bAllPassed && test("mississipissippi", "", false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "", false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "", false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "", false)
		bAllPassed = bAllPassed && test("xxxxzzzzzzzzyf", "", false)
		bAllPassed = bAllPassed && test("xyxyxyzyxyz", "", false)
		bAllPassed = bAllPassed && test("mississippi", "", false)
		bAllPassed = bAllPassed && test("xyxyxyxyz", "", false)
		bAllPassed = bAllPassed && test("m ississippi", "", false)
		bAllPassed = bAllPassed && test("ababac", "", false)
		bAllPassed = bAllPassed && test("dababac", "", false)
		bAllPassed = bAllPassed && test("aaazz", "", false)
		bAllPassed = bAllPassed && test("a12b12", "", false)
		bAllPassed = bAllPassed && test("a12b12", "", false)
		bAllPassed = bAllPassed && test("a12b12", "", false)

		// Another mix of cases.
		bAllPassed = bAllPassed && test("n", "", false)
		bAllPassed = bAllPassed && test("aabab", "", false)
		bAllPassed = bAllPassed && test("ar", "", false)
		bAllPassed = bAllPassed && test("aar", "", false)
		bAllPassed = bAllPassed && test("XYXYXYZYXYz", "", false)
		bAllPassed = bAllPassed && test("missisSIPpi", "", false)
		bAllPassed = bAllPassed && test("mississipPI", "", false)
		bAllPassed = bAllPassed && test("xyxyxyxyz", "", false)
		bAllPassed = bAllPassed && test("miSsissippi", "", false)
		bAllPassed = bAllPassed && test("miSsissippi", "", false)
		bAllPassed = bAllPassed && test("abAbac", "", false)
		bAllPassed = bAllPassed && test("abAbac", "", false)
		bAllPassed = bAllPassed && test("aAazz", "", false)
		bAllPassed = bAllPassed && test("A12b12", "", false)
		bAllPassed = bAllPassed && test("a12B12", "", false)
		bAllPassed = bAllPassed && test("oWn", "", false)
		bAllPassed = bAllPassed && test("bLah", "", false)
		bAllPassed = bAllPassed && test("bLah", "", false)
	}

	if bAllPassed {
		fmt.Println("Passed empty string tests")
	} else {
		fmt.Println("Failed empty string tests")
	}
}

// Correctness tests for a case-sensitive arrangement for invoking a
// UTF-8-enabled routine for matching wildcards.  See relevant code /
// comments in test().
func testUtf8() {
	bAllPassed := true
	bTestingUtf8 = true

	// Simple correctness tests involving various UTF-8 symbols and
	// international content.
	bAllPassed = bAllPassed && test("🐂🚀♥🍀貔貅🦁★□√🚦€¥☯🐴😊🍓🐕🎺🧊☀☂🐉",
		"*☂🐉", true)

	if bCompareCaseInsensitive {
		bAllPassed = bAllPassed && test("AbCD", "abc?", true)
		bAllPassed = bAllPassed && test("AbC★", "abc?", true)
		bAllPassed = bAllPassed && test("⚛⚖☁o", "⚛⚖☁O", true)
	}

	bAllPassed = bAllPassed && test("▲●🐎✗🤣🐶♫🌻ॐ", "▲●☂*", false)
	bAllPassed = bAllPassed && test("𓋍𓋔𓎍", "𓋍𓋔?", true)
	bAllPassed = bAllPassed && test("𓋍𓋔𓎍", "𓋍?𓋔𓎍", false)
	bAllPassed = bAllPassed && test("♅☌♇", "♅☌♇", true)
	bAllPassed = bAllPassed && test("⚛⚖☁", "⚛🍄☁", false)
	bAllPassed = bAllPassed && test("⚛⚖☁O", "⚛⚖☁0", false)
	bAllPassed = bAllPassed && test("गते गते पारगते पारसंगते बोधि स्वाहा",
		"गते गते पारगते प????गते बोधि स्वाहा", true)
	bAllPassed = bAllPassed && test(
		"Мне нужно выучить русский язык, чтобы лучше оценить Пушкина.",
		"Мне нужно выучить * язык, чтобы лучше оценить *.", true)
	bAllPassed = bAllPassed && test(
		"אני צריך ללמוד אנגלית כדי להעריך את גינסברג",
		" אני צריך ללמוד אנגלית כדי להעריך את ???????", false)
	bAllPassed = bAllPassed && test(
		"ગિન્સબર્ગની શ્રેષ્ઠ પ્રશંસા કરવા માટે મારે અંગ્રેજી શીખવું પડશે.",
		"* શ્રેષ્ઠ પ્રશંસા કરવા માટે મારે * શીખવું પડશે.", true)
	bAllPassed = bAllPassed && test(
		"ગિન્સબર્ગની શ્રેષ્ઠ પ્રશંસા કરવા માટે મારે અંગ્રેજી શીખવું પડશે.",
		"??????????? શ્રેષ્ઠ પ્રશંસા કરવા માટે મારે * શીખવું પડશે.", true)
	bAllPassed = bAllPassed && test(
		"ગિન્સબર્ગની શ્રેષ્ઠ પ્રશંસા કરવા માટે મારે અંગ્રેજી શીખવું પડશે.",
		"ગિન્સબર્ગની શ્રેષ્ઠ પ્રશંસા કરવા માટે મારે હિબ્રુ ભાષા શીખવી પડશે.", false)

	// These tests involve multiple=byte code points that contain bytes
	// identical to the single-byte code points for '*' and '?'.
	bAllPassed = bAllPassed && test("ḪؿꜪἪꜿ", "ḪؿꜪἪꜿ", true)
	bAllPassed = bAllPassed && test("ḪؿUἪꜿ", "ḪؿꜪἪꜿ", false)
	bAllPassed = bAllPassed && test("ḪؿꜪἪꜿ", "ḪؿꜪἪꜿЖ", false)
	bAllPassed = bAllPassed && test("ḪؿꜪἪꜿ", "ЬḪؿꜪἪꜿ", false)
	bAllPassed = bAllPassed && test("ḪؿꜪἪꜿ", "?ؿꜪ*ꜿ", true)

	if bAllPassed {
		fmt.Println("Passed UTF-8 tests")
	} else {
		fmt.Println("Failed UTF-8 tests")
	}
}

// Entry point for the Rust executable.  Performance findings (if any) are
// displayed here, once all tests have run.
func main() {
	// Accumulate timing data for all implementations invoked in test().
	if bTestTame {
		testTame()
	}

	if bTestEmpty {
		testEmpty()
	}

	if bTestWild {
		testWild()
	}

	if bTestUtf8 {
		testUtf8()
	}

	if bComparePerformance {
		// Timings have been accumulated via package-scope data.
		fBase := 10.0
		fExpNanoseconds := 9.0
		fExpMilliseconds := 3.0

		// Represent the timings in seconds, to millisecond precision.
		fTimeCumulativeAsciiVersion := (float64(iAccumulatedTimeAscii) /
			math.Pow(fBase, fExpNanoseconds)) * math.Pow(fBase, fExpMilliseconds)
		fTimeCumulativeUtf8Version := (float64(iAccumulatedTimeUTF8) /
			math.Pow(fBase, fExpNanoseconds)) * math.Pow(fBase, fExpMilliseconds)
		// Can add similar calculations for more performance comparisons...

		fUtf8VersionTimeInSeconds := fTimeCumulativeUtf8Version / 1000
		fAsciiVersionTimeInSeconds := fTimeCumulativeAsciiVersion / 1000

		// Show the timing results.
		fmt.Printf(
			"FastWildCompareAscii() - for ASCII strings: %.3f seconds\n",
			fAsciiVersionTimeInSeconds)
		// Can add results for more performance comparisons here...
		fmt.Printf(
			"FastWildCompareRuneSlices() - for UTF-8-encoded strings: %.3f seconds\n",
			fUtf8VersionTimeInSeconds)
	}
}
