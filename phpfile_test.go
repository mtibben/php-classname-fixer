package main

import (
	"testing"
)

func TestDisplaynameSkel(t *testing.T) {

	t1 := `
use Edge_Database_ConnectionPool as ConnectionPool;
use Edge_Database_LoadBalancedDatabase as LoadBalancedDatabase;
`
	matches := useasRe.FindAllStringSubmatch(t1, -1)

}
