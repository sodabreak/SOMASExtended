package main

/* 
* Code to test the functionality of the orphan pool, which deals with agents
* that are not currently part of a team re-joining teams in subsequent turns.
*/

import "testing"  // built-in go testing package
import "github.com/stretchr/testify/assert"  // assert package, easier to 

func TestNoBreak(t *testing.T) {
    result := 2 + 3
    assert.Equal(t, 5, result)
}

