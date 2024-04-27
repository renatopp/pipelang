package object_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestList_Size(t *testing.T) {
	common.AssertCode(t, `[1,2,3].Size()`, `3`)
}

func TestList_Set(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Set(1, 4); a`, `[1, 4, 3]`)
	common.AssertCode(t, `a := [1,2,3]; a.Set(-1, 4); a`, `[1, 2, 4]`)
	common.AssertCodeError(t, `a := [1,2,3]; a.Set(5, 4); a`)
}

func TestList_Assign(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a[1] = 4; a`, `[1, 4, 3]`)
	common.AssertCode(t, `a := [1,2,3]; a[-1] = 4; a`, `[1, 2, 4]`)
	common.AssertCodeError(t, `a := [1,2,3]; a[5] = 4; a`)
	common.AssertCodeError(t, `a := [1,2,3]; a[5, 3] = 4; a`)
}

func TestList_Push(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Push(4); a`, `[1, 2, 3, 4]`)
	common.AssertCode(t, `a := [1,2,3]; a.Push(4); a.Push(0); a`, `[1, 2, 3, 4, 0]`)
}

func TestList_Pop(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Pop(); a`, `[1, 2]`)
	common.AssertCode(t, `a := [1,2,3]; a.Pop(); a.Pop(); a`, `[1]`)
}

func TestList_Insert(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Insert(1, 4); a`, `[1, 4, 2, 3]`)
}

func TestList_Remove(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Remove(2); a`, `[1, 3]`)
	common.AssertCode(t, `a := [1,2,3]; a.Remove(2); a.Remove(1); a`, `[3]`)
}

func TestList_RemoveAt(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.RemoveAt(1); a`, `[1, 3]`)
	common.AssertCode(t, `a := [1,2,3]; a.RemoveAt(1); a.RemoveAt(0); a`, `[3]`)
}

func TestList_Copy(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; b := a.Copy(); b`, `[1, 2, 3]`)
}

func TestList_Clear(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Clear(); a`, `[]`)
}

func TestList_Concat(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Concat([4,5])`, `[1, 2, 3, 4, 5]`)
	common.AssertCode(t, `a := [1,2,3]; a.Concat([4,5], [3, 4])`, `[1, 2, 3, 4, 5, 3, 4]`)
	common.AssertCodeError(t, `a := [1,2,3]; a.Concat(4, 5)`)
}

func TestList_Split(t *testing.T) {
	common.AssertCode(t, `[1,2,3,2, 4,5].Split(2)`, `[[1], [3], [4, 5]]`)
}

func TestList_SplitAt(t *testing.T) {
	common.AssertCode(t, `[1,2,3,4,5].SplitAt(2)`, `[[1, 2], [3, 4, 5]]`)
}

func TestList_SplitFn(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3,2,4,5]; a.SplitFn(x :x > 2)`, `[[1, 2], [3, 2], [4], [5]]`)
}
func TestList_Sorted(t *testing.T) {
	common.AssertCode(t, `a := [3,1,2]; a.Sorted(); a`, `[1, 2, 3]`)
}

func TestList_SortedFn(t *testing.T) {
	common.AssertCode(t, `a := [3,1,2]; a.SortedFn((x, y): x > y); a`, `[3, 2, 1]`)
}

func TestList_Reversed(t *testing.T) {
	common.AssertCode(t, `[1,2,3].Reversed()`, `[3, 2, 1]`)
}

func TestList_Indexing(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a[1]`, `2`)
	common.AssertCode(t, `a := [1,2,3]; a[-1]`, `3`)
	common.AssertCodeError(t, `a := [1,2,3]; a[5]`)
	common.AssertCode(t, `a := [1,2,3]; a[1, 2]`, `[2]`)
	common.AssertCode(t, `a := [1,2,3]; a[-3, -1]`, `[1, 2]`)
	common.AssertCode(t, `a := [1,2,3]; a[-1, 1]`, `[]`)
}

func TestList_Get(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Get(1)`, `2`)
	common.AssertCode(t, `a := [1,2,3]; a.Get(-1)`, `3`)
	common.AssertCodeError(t, `a := [1,2,3]; a.Get(5)`)
}

func TestList_GetOr(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.GetOr(1, 0)`, `2`)
	common.AssertCode(t, `a := [1,2,3]; a.GetOr(-1, 0)`, `3`)
	common.AssertCode(t, `a := [1,2,3]; a.GetOr(5, 0)`, `0`)
}

func TestList_Sub(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Sub(1, 2)`, `[2]`)
	common.AssertCode(t, `a := [1,2,3]; a.Sub(-3, -1)`, `[1, 2]`)
	common.AssertCode(t, `a := [1,2,3]; a.Sub(-1, 1)`, `[]`)
}

func TestList_Find(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Find(2)`, `1`)
	common.AssertCode(t, `a := [1,2,3]; a.Find(4)`, `-1`)
}

func TestList_FindFn(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.FindFn(x: x > 1)`, `1`)
	common.AssertCode(t, `a := [1,2,3]; a.FindFn(x: x > 3)`, `-1`)
}

func TestList_FindLast(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3,2]; a.FindLast(2)`, `3`)
	common.AssertCode(t, `a := [1,2,3]; a.FindLast(4)`, `-1`)
}

func TestList_FindLastFn(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3,2]; a.FindLastFn(x: x > 1)`, `3`)
	common.AssertCode(t, `a := [1,2,3]; a.FindLastFn(x: x > 3)`, `-1`)
}

func TestList_FindAll(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3,2]; a.FindAll(2)`, `[1, 3]`)
	common.AssertCode(t, `a := [1,2,3]; a.FindAll(4)`, `[]`)
}

func TestList_FindAllFn(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3,2]; a.FindAllFn(x: x > 1)`, `[1, 2, 3]`)
	common.AssertCode(t, `a := [1,2,3]; a.FindAllFn(x: x > 3)`, `[]`)
}

func TestList_Contains(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Contains(2)`, `true`)
	common.AssertCode(t, `a := [1,2,3]; a.Contains(4)`, `false`)
}

func TestList_ContainsFn(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.ContainsFn(x: x > 1)`, `true`)
	common.AssertCode(t, `a := [1,2,3]; a.ContainsFn(x: x > 3)`, `false`)
}

func TestList_IsEmpty(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.IsEmpty()`, `false`)
	common.AssertCode(t, `a := []; a.IsEmpty()`, `true`)
}

func TestList_Count(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Count(2)`, `1`)
	common.AssertCode(t, `a := [1,2,3]; a.Count(4)`, `0`)
}

func TestList_CountFn(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.CountFn(x: x > 1)`, `2`)
	common.AssertCode(t, `a := [1,2,3]; a.CountFn(x: x > 3)`, `0`)
}

func TestList_Join(t *testing.T) {
	common.AssertCode(t, `a := [1,2,3]; a.Join(",")`, `1,2,3`)
}

func TestList_Elements(t *testing.T) {
	common.AssertCode(t, `
		a := [1,2,3];
		el := a.Elements()
		sum := 0
		for e in el {
			sum += e
		}	
	`, `6`)
	common.AssertCode(t, `
		a := [1,2,3];
		sum := 0
		for e in a {
			sum += e
		}	
	`, `6`)
}
