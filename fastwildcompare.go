// Go routines for matching wildcards.
//
// Copyright 2025 Kirk J Krauss.  This is a Derivative Work based on 
// material that is copyright 2025 Kirk J Krauss and available at
//
//     https://developforperformance.com/MatchingWildcardsInRust.html
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     https://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Go implementation of fast_wild_compare_ascii(), for ASCII text.
//
// Compares two ASCII strings.  Accepts '?' as a single-character wildcard.
// For each '*' wildcard, seeks out a matching sequence of any characters 
// beyond it.  Otherwise compares the strings a character at a time. 
//
package main

func FastWildCompareAscii(strWild, strTame string) bool {
	var iWild int = 0     // Index for both input strings in upper loop
	var iTame int         // Index for tame content, used in lower loop
	var iWildSequence int // Index for prospective match after '*'
	var iTameSequence int // Index for match in tame content

    // Find a first wildcard, if one exists, and the beginning of any  
    // prospectively matching sequence after it.
    for {
		// Check for the end from the start.  Get out fast, if possible.
		if len(strTame) <= iWild {
			if len(strWild) > iWild {
				for strWild[iWild] == '*' {
					iWild++
					
					if len(strWild) <= iWild {
						return true        // "ab" matches "ab*".
					}
				}

			    return false               // "abcd" doesn't match "abc".
			} else {
				return true                // "abc" matches "abc".
			}
		} else if len(strWild) <= iWild {
		    return false                   // "abc" doesn't match "abcd".
		} else if strWild[iWild] == '*' {
			// Got wild: set up for the second loop and skip on down there.
			iTame = iWild

			for {
				iWild++

				if len(strWild) <= iWild {
					return true            // "abc*" matches "abcd".
				}
				
				if strWild[iWild] != '*' {
					break
				}
			}

			// Search for the next prospective match.
			if strWild[iWild] != '?' {
				for strWild[iWild] != strTame[iTame] {
					iTame++

					if len(strTame) <= iTame {
						return false       // "a*bc" doesn't match "ab".
					}
				}
			}

			// Keep fallback positions for retry in case of incomplete match.
			iWildSequence = iWild
			iTameSequence = iTame
			break
		} else if strWild[iWild] != strTame[iWild] && strWild[iWild] != '?' {
			return false                   // "abc" doesn't match "abd".
		}

		iWild++                            // Everything's a match, so far.
	}

    // Find any further wildcards and any further matching sequences.
    for {
		if len(strWild) > iWild && strWild[iWild] == '*' {
            // Got wild again.
			for {
				iWild++

				if len(strWild) <= iWild {
					return true            // "ab*c*" matches "abcd".
				}
				
				if strWild[iWild] != '*' {
					break
				}
			}

			if len(strTame) <= iTame {
                return false               // "*bcd*" doesn't match "abc".
            }

            // Search for the next prospective match.
            if strWild[iWild] != '?' {
                for len(strTame) > iTame && 
				    strWild[iWild] != strTame[iTame] {
					iTame++

                    if len(strTame) <= iTame {
                        return false       // "a*b*c" doesn't match "ab".
                    }
                }
            }

            // Keep the new fallback positions.
			iWildSequence = iWild
			iTameSequence = iTame
        } else {
            // The equivalent portion of the upper loop is really simple.
            if len(strTame) <= iTame {
				if len(strWild) <= iWild {
					return true            // "*b*c" matches "abc".
				}
			
                return false               // "*bcd" doesn't match "abc".
            }
			
			if len(strWild) <= iWild ||
		       (strWild[iWild] != strTame[iTame] && 
		        strWild[iWild] != '?') {
				// A fine time for questions.
				for len(strWild) > iWildSequence && 
				    strWild[iWildSequence] == '?' {
					iWildSequence++
					iTameSequence++
				}

				iWild = iWildSequence

				// Fall back, but never so far again.
				for	{
					iTameSequence++

					if len(strTame) <= iTameSequence {
						if len(strWild) <= iWild {
							return true  // "*a*b" matches "ab".
						} else {
							return false // "*a*b" doesn't match "ac".
						}
					}

					if len(strWild) > iWild && 
					   strWild[iWild] == strTame[iTameSequence] {
						break
					}
				}

	            iTame = iTameSequence
			}
        }

        // Another check for the end, at the end.
        if len(strTame) <= iTame {
			if len(strWild) <= iWild {
				return true                // "*bc" matches "abc".
			}

			return false                   // "*bc" doesn't match "abcd".
		}

        iWild++                            // Everything's still a match.
        iTame++
    }
}

// Go implementation of fast_wild_compare_utf8(), for UTF-8-encoded rune 
// slices.
//
// Compares two rune slices.  Accepts '?' as a single-rune wildcard.  For 
// each '*' wildcard, seeks out a matching sequence of any runes beyond it.  
// Otherwise compares the slices a rune at a time. 
//
func FastWildCompareRuneSlices(rslcWild, rslcTame []rune) bool {
	var iWild int = 0     // Index for both input strings in upper loop
	var iTame int         // Index for tame content, used in lower loop
	var iWildSequence int // Index for prospective match after '*'
	var iTameSequence int // Index for match in tame content
	
    // Find a first wildcard, if one exists, and the beginning of any  
    // prospectively matching sequence after it.
    for {
		// Check for the end from the start.  Get out fast, if possible.
		if len(rslcTame) <= iWild {
			if len(rslcWild) > iWild {
				for rslcWild[iWild] == '*' {
					iWild++
					
					if len(rslcWild) <= iWild {
						return true        // "ab" matches "ab*".
					}
				}

			    return false               // "abcd" doesn't match "abc".
			} else {
				return true                // "abc" matches "abc".
			}
		} else if len(rslcWild) <= iWild {
		    return false                   // "abc" doesn't match "abcd".
		} else if rslcWild[iWild] == '*' {
			// Got wild: set up for the second loop and skip on down there.
			iTame = iWild

			for {
				iWild++

				if len(rslcWild) <= iWild {
					return true            // "abc*" matches "abcd".
				}
				
				if rslcWild[iWild] != '*' {
					break
				}
			}

			// Search for the next prospective match.
			if rslcWild[iWild] != '?' {
				for rslcWild[iWild] != rslcTame[iTame] {
					iTame++

					if len(rslcTame) <= iTame {
						return false       // "a*bc" doesn't match "ab".
					}
				}
			}

			// Keep fallback positions for retry in case of incomplete match.
			iWildSequence = iWild
			iTameSequence = iTame
			break
		} else if rslcWild[iWild] != rslcTame[iWild] && rslcWild[iWild] != '?' {
			return false                   // "abc" doesn't match "abd".
		}

		iWild++                            // Everything's a match, so far.
	}

    // Find any further wildcards and any further matching sequences.
    for {
		if len(rslcWild) > iWild && rslcWild[iWild] == '*' {
            // Got wild again.
			for {
				iWild++

				if len(rslcWild) <= iWild {
					return true            // "ab*c*" matches "abcd".
				}
				
				if rslcWild[iWild] != '*' {
					break
				}
			}

			if len(rslcTame) <= iTame {
                return false               // "*bcd*" doesn't match "abc".
            }

            // Search for the next prospective match.
            if rslcWild[iWild] != '?' {
                for len(rslcTame) > iTame && 
				    rslcWild[iWild] != rslcTame[iTame] {
					iTame++

                    if len(rslcTame) <= iTame {
                        return false       // "a*b*c" doesn't match "ab".
                    }
                }
            }

            // Keep the new fallback positions.
			iWildSequence = iWild
			iTameSequence = iTame
        } else {
            // The equivalent portion of the upper loop is really simple.
            if len(rslcTame) <= iTame {
				if len(rslcWild) <= iWild {
					return true            // "*b*c" matches "abc".
				}
			
                return false               // "*bcd" doesn't match "abc".
            }
			
			if len(rslcWild) <= iWild ||
		       rslcWild[iWild] != rslcTame[iTame] && 
		       rslcWild[iWild] != '?' {
				// A fine time for questions.
				for len(rslcWild) > iWildSequence && 
				    rslcWild[iWildSequence] == '?' {
					iWildSequence++
					iTameSequence++
				}

				iWild = iWildSequence

				// Fall back, but never so far again.
				for	{
					iTameSequence++

					if len(rslcTame) <= iTameSequence {
						if len(rslcWild) <= iWild {
							return true    // "*a*b" matches "ab".
						} else {
							return false   // "*a*b" doesn't match "ac".
						}
					}

					if len(rslcWild) > iWild && 
					   rslcWild[iWild] == rslcTame[iTameSequence] {
						break
					}
				}

	            iTame = iTameSequence
			}
        }

        // Another check for the end, at the end.
        if len(rslcTame) <= iTame {
			if len(rslcWild) <= iWild {
				return true                // "*bc" matches "abc".
			}

			return false                   // "*bc" doesn't match "abcd".
		}

        iWild++                            // Everything's still a match.
        iTame++
    }
}
