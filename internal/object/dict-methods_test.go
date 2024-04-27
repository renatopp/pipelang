package object_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestDict_GetSet(t *testing.T) {
	common.AssertCode(t, `a := {}; a.Set("key", "value"); a.Get("key")`, "value")
	common.AssertCode(t, `a := {}; a['key'] = 1; a['key']`, "1")
	common.AssertCode(t, `a := {}; a[3] = 1; a['3']`, "1")
	common.AssertCode(t, `a := {}; a[3] = 1; a[3]`, "1")
	common.AssertCodeError(t, `a := {}; a['key']`)
}

func TestDict_GetOr(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.GetOr("key", "default")`, "44")
	common.AssertCode(t, `a := {}; a.GetOr("key", "default")`, "default")
}

func TestDict_Has(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.Has(3)`, "false")
	common.AssertCode(t, `a := { key=44 }; a.Has(44)`, "false")
	common.AssertCode(t, `a := { key=44 }; a.Has('key')`, "true")
}

func TestDict_Remove(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.Remove('key')`, "44")
	common.AssertCode(t, `a := { key=44 }; a.Remove('key'); a `, "{}")
	common.AssertCodeError(t, `a := { key=44 }; a.Remove('unknown')`)
}

func TestDict_Contains(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.Contains('key')`, "false")
	common.AssertCode(t, `a := { key=44 }; a.Contains('unknown')`, "false")
	common.AssertCode(t, `a := { key=44 }; a.Contains(44)`, "true")
}

func TestDict_ContainsFn(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.ContainsFn(fn (x) { x == 44 })`, "true")
	common.AssertCode(t, `a := { key=44 }; a.ContainsFn(fn (x) { x == 45 })`, "false")
}

func TestDict_Size(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.Size()`, "1")
	common.AssertCode(t, `a := { key=44, key2=45 }; a.Size()`, "2")
	common.AssertCode(t, `a := { }; a.Size()`, "0")
}

func TestDict_Clear(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; a.Clear(); a.Size()`, "0")
	common.AssertCode(t, `a := { key=44 }; a.Clear(); a`, "{}")
}

func TestDict_Copy(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; b := a.Copy(); b`, "{key=44}")
	common.AssertCode(t, `a := { key=44 }; b := a.Copy(); b.Set('key', 45); a`, "{key=44}")
}

func TestDict_Concat(t *testing.T) {
	common.AssertCode(t, `a := { key=44 }; b := { key2=45 }; a.Concat(b)['key']`, "44")
	common.AssertCode(t, `a := { key=44 }; b := { key2=45 }; a.Concat(b)['key2']`, "45")
	common.AssertCode(t, `a := { key=44 }; b := { key2=45 }; a.Concat(b); b`, "{key2=45}")
	common.AssertCode(t, `a := { key=44 }; b := { key2=45 }; a.Concat(b); a.Size()`, "2")
	common.AssertCode(t, `a := { key=44 }; b := { key=45 }; a.Concat(b); a`, "{key=45}")
}

func TestDict_Find(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Find(1)`, "x")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Find(3)`, "z")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Find(4)`, "false")
}

func TestDict_FindFn(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindFn(fn (x) { x == 1 })`, "x")
	common.AssertCode(t, `a := { x=1, z=3 }; a.FindFn(fn (x) { x > 1 })`, "z")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindFn(fn (x) { x > 3 })`, "false")
}

func TestDict_FindAll(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindAll(1)`, "['x']")
	common.AssertCode(t, `a := { x=1, y=2, z=2 }; a.FindAll(2).Sorted()`, "['y', 'z']")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindAll(2, 3).Sorted()`, "['y', 'z']")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindAll(4)`, "[]")
}

func TestDict_FindAllFn(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindAllFn(fn (x) { x == 1 })`, "['x']")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindAllFn(fn (x) { x > 1 }).Sorted()`, "['y', 'z']")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.FindAllFn(fn (x) { x > 3 })`, "[]")
}

func TestDict_IsEmpty(t *testing.T) {
	common.AssertCode(t, `a := { x=1 }; a.IsEmpty()`, "false")
	common.AssertCode(t, `a := { }; a.IsEmpty()`, "true")
}

func TestDict_Count(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Count(1)`, "1")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Count(2, 3)`, "2")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Count(4)`, "0")
}

func TestDict_CountFn(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.CountFn(fn (x) { x == 1 })`, "1")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.CountFn(fn (x) { x > 1 })`, "2")
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.CountFn(fn (x) { x > 3 })`, "0")
}

func TestDict_Keys(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Keys().Sorted()`, "['x', 'y', 'z']")
	common.AssertCode(t, `a := { }; a.Keys()`, "[]")
}

func TestDict_Values(t *testing.T) {
	common.AssertCode(t, `a := { x=1, y=2, z=3 }; a.Values().Sorted()`, "[1, 2, 3]")
	common.AssertCode(t, `a := { }; a.Values()`, "[]")
}

func TestDict_Items(t *testing.T) {
	common.AssertCode(t, `a := { }; a.Items()`, "[]")
}
