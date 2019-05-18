package scopes

import (
	"strings"
	"sync"
)

//strech connverts a slice to bigger length by adding toAdd to it
func strech(scopeSlice []string, toAdd string, toLen int) []string {
	lenDiff := toLen - len(scopeSlice)
	for i := 0; i < lenDiff; i++ {
		scopeSlice = append(scopeSlice, toAdd)
	}
	return scopeSlice
}

//MatchScopes matches two scopes using Wildcard Scope Matching Strategy (asymetric)
func MatchScopes(scopeA, scopeB string) bool {
	scopeASplit := strings.Split(scopeA, ":")
	scopeBSplit := strings.Split(scopeB, ":")
	scopeALen := len(scopeASplit)
	scopeBLen := len(scopeBSplit)

	// If scopeBLen is smaller than scopeALen and last char of scopeB is not * return false
	if scopeBLen < scopeALen && scopeBSplit[scopeBLen-1] != "*" {
		return false
		// If scopeBLen is smaller than scopeALen and last char of scopeB is * stretch scopeB To Len Of ScopeA By Adding "*"
	} else if scopeBLen < scopeALen && scopeBSplit[scopeBLen-1] == "*" {
		scopeBSplit = strech(scopeBSplit, "*", scopeALen)
		// If scopeBLen is greater than scopeALen and last char of scopeA is not * return false
	} else if scopeBLen > scopeALen && scopeASplit[scopeALen-1] != "*" {
		return false
		// If scopeBLen is greater than scopeALen and last char of scopeA is * stretch scopeA To Len Of ScopeB By Adding "*"
	} else if scopeBLen > scopeALen && scopeASplit[scopeALen-1] == "*" {
		scopeASplit = strech(scopeASplit, "*", scopeBLen)
	}

	for i := 0; i < scopeALen; i++ {
		if !(scopeASplit[i] == scopeBSplit[i] || scopeBSplit[i] == "*") {
			return false
		}
	}

	return true
}

var wg sync.WaitGroup

/*ScopeInAllowedScopes checks if a scope is in a list of allowed scopes
then adds the result to a channel */
func ScopeInAllowedScopes(channel chan bool, scope string, allowedScopes []string) {
	for _, allowedScope := range allowedScopes {
		if MatchScopes(scope, allowedScope) {
			channel <- true
			wg.Done()
			return
		}
	}
	channel <- false
	wg.Done()
}

//MatchScopesRequestedToScopesAllowed Checks if scopes requested are allowed (concurrently)
func MatchScopesRequestedToScopesAllowed(scopesRequested, scopesAllowed []string) bool {
	matchedChan := make(chan bool)
	wg.Add(len(scopesRequested))

	/*for each scope requested it will run a goroutine where
	that scope is checked against the list of allowed scopes
	each goroutine will then add their results to the matchedChan*/
	for _, scopeRequested := range scopesRequested {
		go ScopeInAllowedScopes(matchedChan, scopeRequested, scopesAllowed)
	}

	//If the matchedChan channel has false message, the program will return false
	for i := 0; i < len(scopesRequested); i++ {
		select {
		case match := <-matchedChan:
			if !match {
				return false
			}
		}
	}

	//Return true if all goroutines are finished and matchedChan didn't contain false
	wg.Wait()
	close(matchedChan)
	return true
}
