/* The ns package provides the net-score.org functionality.  */
package igo

// Unit tests for the LogList data structure.

import (
	"fmt"
	"math/rand"
	"testing"
)

// Struct for organizing test data.
type Node struct {
	m string // Method (GET, POST)
	r string // Resource (/, /feather/3456, /upload)
	i string // IP (192.168.1.1, 2620:0:1000:3801:a800:1ff:fe00:df)
	t int64  // Unix time (1337794942)
	v int64  // Value (1, 3, 5)
}

// Common references to test data.
var (
	exactPass *Node
	fuzzyPass *Node
	fuzzyFail *Node
	nodeList  []*Node
	logList   *LogList
	tolerance int64
)

/*Gets things ready to test.

This should be called first thing in every test.  It makes sure that the state
is set back to the starting state.  Nothing is placed in the LogList so that we
can test things on empty lists - see loadLogList for that.
*/
func setup() {
	logList = &LogList{}

	tolerance = 10 // seconds
	nodeList = make([]*Node, 6)

	// Search items
	exactPass = &Node{
		m: "GET",
		r: "/feather/1234",
		i: "192.168.1.25",
		t: 1337794942,
		v: 42,
	}
	fuzzyPass = &Node{
		m: exactPass.m,
		r: exactPass.r,
		i: exactPass.i,
		t: exactPass.t - (tolerance / 3),
		v: exactPass.v,
	}
	fuzzyFail = &Node{
		m: exactPass.m,
		r: exactPass.r,
		i: exactPass.i,
		t: exactPass.t + (tolerance + 1),
		v: exactPass.v,
	}

	// Search space
	nodeList[0] = exactPass
	nodeList[1] = &Node{
		m: "POST",
		r: "/upload/",
		i: "192.168.1.25",
		t: 1337794940,
		v: 7,
	}
	nodeList[2] = &Node{
		m: "POST",
		r: "/upload/",
		i: "192.168.100.2",
		t: 1337790000,
		v: 73,
	}
	nodeList[3] = &Node{
		m: "GET",
		r: "/feather/1234",
		i: "192.168.100.2",
		t: 1337791000,
		v: 75,
	}
	nodeList[4] = &Node{
		m: "GET",
		r: "/egg/1234",
		i: "192.168.100.2",
		t: 1337791100,
		v: 77,
	}
	nodeList[5] = &Node{
		m: "GET",
		r: "/",
		i: "192.168.99.2",
		t: 1337791100,
		v: 79,
	}
}

/*Loads the test data into the LogList.

Takes all of the Node structs in the nodeList and places them in the logList for
testing.
*/
func loadLogList() error {
	for i := range nodeList {
		err := logList.Push(nodeList[i].m, nodeList[i].r, nodeList[i].i,
			nodeList[i].t, nodeList[i].v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Makes sure that Init zeros things out properly.
func TestInit(t *testing.T) {
	setup()
	logList.front = &Element{}
	logList.length = 77
	logList.tolerance = 5
	logList.Init()
	if logList.front != nil {
		t.Fail()
	}
	if logList.length != 0 {
		t.Fail()
	}
	if logList.tolerance != 0 {
		t.Fail()
	}
}

// Makes sure that the tolerance settings are being set.
func TestSetTolerance(t *testing.T) {
	setup()
	if logList.tolerance != 0 {
		t.Fail()
	}
	logList.SetTolerance(tolerance)
	if logList.tolerance != 10 {
		t.Fail()
	}
}

// Makes sure that the Empty function works.
func TestEmpty(t *testing.T) {
	setup()
	testElement := &Element{
		Value: 999,
	}

	if !logList.Empty() {
		t.Fail()
	}
	logList.front = testElement
	if logList.Empty() {
		t.Fail()
	}
}

/*Testing the Push function.

This test makes sure that the Push function adds the Node to the list and that
it updates the length accordingly.
*/
func TestPushPass(t *testing.T) {
	setup()
	if logList.front != nil {
		t.Fail()
	}
	err := logList.Push(exactPass.m, exactPass.r, exactPass.i, exactPass.t,
		exactPass.v)
	if err != nil {
		t.Errorf("TestPushPass:logList.Push err = %v", err)
	}
	if logList.front == nil {
		t.Fail()
	}
	if logList.length != 1 {
		t.Fail()
	}
	err = logList.Push(fuzzyPass.m, fuzzyPass.r, fuzzyPass.i, fuzzyPass.t,
		fuzzyPass.v)
	if err != nil {
		t.Errorf("TestPushPass:logList.Push err = %v", err)
	}
	if logList.front == nil {
		t.Fail()
	}
	if logList.length != 2 {
		t.Fail()
	}
}

/*Testing the Push function for collision detection.

This test makes sure that the Push function adds the Node to the list and that
it updates the length accordingly.  In this case the Push should fail with a
specific error, and the length should not change.
*/
func TestPushFail(t *testing.T) {
	setup()
	if logList.front != nil {
		t.Fail()
	}
	err := logList.Push(exactPass.m, exactPass.r, exactPass.i, exactPass.t,
		exactPass.v)
	if err != nil {
		t.Errorf("TestPushFail:logList.Push err = %v", err)
	}
	if logList.front == nil {
		t.Fail()
	}
	if logList.length != 1 {
		t.Fail()
	}
	// Cause a collision!
	err = logList.Push(exactPass.m, exactPass.r, exactPass.i, exactPass.t,
		exactPass.v)
	if err != LogListCollisionError {
		t.Errorf("TestPushFail:logList.Push err = %v", err)
	}
	if logList.front == nil {
		t.Fail()
	}
	if logList.length != 1 {
		t.Fail()
	}
}

// Makes sure that the Size function works.
func TestSize(t *testing.T) {
	setup()
	if logList.Size() != 0 {
		t.Fail()
	}
	err := logList.Push(exactPass.m, exactPass.r, exactPass.i, exactPass.t,
		exactPass.v)
	if err != nil {
		t.Errorf("TestSize:logList.Push err = %v", err)
	}
	if logList.Size() != 1 {
		t.Fail()
	}
}

/*Testing the Pop function.

Trying to Pop on an empty list should generate a specific error and the length
should remain unchanged.
*/
func TestPopEmpty(t *testing.T) {
	setup()
	if !logList.Empty() {
		t.Fail()
	}
	if logList.Size() != 0 {
		t.Fail()
	}
	_, err := logList.Pop()
	if err != LogListEmptyError {
		t.Errorf("TestPopEmpty:logList.Pop err = %v", err)
	}
	if logList.Size() != 0 {
		t.Fail()
	}
}

/*Testing the Pop function.

Trying to Pop on an NOT empty list should NOT generate an error.  The value
returned should match the one we passed in and the length should be one less.
*/
func TestPopFull(t *testing.T) {
	setup()
	err := logList.Push(exactPass.m, exactPass.r, exactPass.i,
		exactPass.t, exactPass.v)
	if err != nil {
		t.Errorf("TestPopFull:logList.Push err = %v", err)
	}
	testValue, err := logList.Pop()
	if err != nil {
		t.Errorf("TestPopFull:logList.Pop err = %v", err)
	}
	if testValue != exactPass.v {
		t.Errorf("TestPopFull:logList.Pop %v != %v", testValue,
			exactPass.v)
	}
	if logList.Size() != 0 {
		t.Fail()
	}
	if !logList.Empty() {
		t.Fail()
	}
}

/*Testing the Match function.

Trying to Match on an empty list should generate a specific error and the length
should remain unchanged.
*/
func TestMatchEmpty(t *testing.T) {
	setup()
	if !logList.Empty() {
		t.Fail()
	}
	if logList.Size() != 0 {
		t.Fail()
	}
	testValue, err := logList.Match(exactPass.m, exactPass.r, exactPass.i,
		exactPass.t)
	if err != LogListMatchNotFoundError {
		t.Fail()
	}
	if testValue != nil {
		t.Fail()
	}
	if logList.Size() != 0 {
		t.Fail()
	}
}

/*Testing the Match function.

Trying to Match on an NOT empty list with an item that is in the list should NOT
generate an error and the length should be one less.  The value returned should
match the value that was stored.
*/
func TestMatchExactPass(t *testing.T) {
	setup()
	err := loadLogList()
	if err != nil {
		t.Errorf("TestMatchExactPass:loadLogList err = %v", err)
	}
	if logList.Size() != int64(len(nodeList)) {
		t.Fail()
	}
	testValue, err := logList.Match(exactPass.m, exactPass.r, exactPass.i,
		exactPass.t)
	if err != nil {
		t.Errorf("TestMatchExactPass:logList.Match err = %v", err)
	}
	if testValue != exactPass.v {
		t.Errorf("TestMatchExactPass:logList.Match %v != %v", testValue,
			exactPass.v)
	}
	if logList.Size() != int64(len(nodeList))-1 {
		t.Fail()
	}
}

/*Testing the Match function.

Trying to Match on an NOT empty list with an item that is in the list should NOT
generate an error and the length should be one less.  Since we are testing the
fuzzy match, the value returned should match the exact value that was stored.
*/
func TestMatchFuzzyPass(t *testing.T) {
	setup()
	err := loadLogList()
	// Relax the tolerance to allow for a fuzzy match
	logList.SetTolerance(tolerance)
	if err != nil {
		t.Errorf("TestMatchFuzzyPass:loadLogList err = %v", err)
	}
	if logList.Size() != int64(len(nodeList)) {
		t.Fail()
	}
	testValue, err := logList.Match(fuzzyPass.m, fuzzyPass.r, fuzzyPass.i,
		fuzzyPass.t)
	if err != nil {
		t.Errorf("TestMatchFuzzyPass:logList.Match err = %v", err)
	}
	if testValue != exactPass.v {
		t.Errorf("TestMatchFuzzyPass:logList.Match %v != %v", testValue,
			exactPass.v)
	}
	if logList.Size() != int64(len(nodeList))-1 {
		t.Fail()
	}
}

/*Testing the Match function.

Trying to Match on an NOT empty list with an item that is NOT in the list should
generate an error and the length should be unchanged.  Since we are testing the
fuzzy match, we are using a search value that will fail by falling outside the
allowed tolerance.
*/
func TestMatchFail(t *testing.T) {
	setup()
	err := loadLogList()
	// Relax the tolerance to allow for a fuzzy match
	logList.SetTolerance(tolerance)
	if err != nil {
		t.Errorf("TestMatchFuzzyFail:loadLogList err = %v", err)
	}
	if logList.Size() != int64(len(nodeList)) {
		t.Fail()
	}
	testValue, err := logList.Match(fuzzyFail.m, fuzzyFail.r, fuzzyFail.i,
		fuzzyFail.t)
	if err != LogListMatchNotFoundError {
		t.Fail()
	}
	if testValue != nil {
		t.Fail()
	}
	if logList.Size() != int64(len(nodeList)) {
		t.Fail()
	}
}

/*Benchmarking tests:

These tests are not run by "go test" by default.  To run them you need to use
the -test.bench flag (http://golang.org/cmd/go/#Description_of_testing_flags).

go test -test.bench="." logList_test.go logList.go
*/

// Setup for benchmarking methods.
func benchmarkSetup(nodeCount int) []*Node {
	nodeList := make([]*Node, nodeCount)

	methodList := []string{"GET", "POST", "GET", "GET", "GET"}

	rand.Seed(1)
	for i := 0; i < nodeCount; i++ {
		im := rand.Int63n(int64(len(methodList)))
		rv := rand.Int63n(10)
		iv := rand.Int63n(255)
		tv := rand.Int63n(1337794942)
		vv := rand.Int63n(1000)
		nodeList[i] = &Node{
			m: methodList[im],
			r: fmt.Sprintf("/feather/%d", rv),
			i: fmt.Sprintf("192.168.1.%d", iv),
			t: tv,
			v: vv,
		}
	}
	return nodeList
}

// Benchmarking test for pushing items onto a LogList.
func BenchmarkPush(b *testing.B) {
	b.StopTimer()
	logList = &LogList{}
	nodeCount := b.N
	nodeList := benchmarkSetup(nodeCount)
	logList.SetTolerance(int64(tolerance))
	b.StartTimer()
	for i := 0; i < nodeCount; i++ {
		err := logList.Push(nodeList[i].m, nodeList[i].r,
			nodeList[i].i, nodeList[i].t, nodeList[i].v)
		if err == LogListCollisionError {
			continue
		}
		if err != nil {
			b.Fail()
		}
	}
}

// Benchmarking test for matching results.
func BenchmarkMatch(b *testing.B) {
	b.StopTimer()
	logList = &LogList{}
	nodeCount := b.N
	nodeList := benchmarkSetup(nodeCount)
	logList.SetTolerance(int64(tolerance))
	for i := 0; i < nodeCount; i++ {
		err := logList.Push(nodeList[i].m, nodeList[i].r,
			nodeList[i].i, nodeList[i].t, nodeList[i].v)
		if err == LogListCollisionError {
			continue
		}
		if err != nil {
			b.Fail()
		}
	}
	b.StartTimer()
	for i := 0; i < nodeCount; i++ {
		_, err := logList.Match(nodeList[i].m, nodeList[i].r,
			nodeList[i].i, nodeList[i].t)
		if err == LogListMatchNotFoundError {
			continue
		}
		if err != nil {
			b.Fail()
		}
	}
}

// Benchmarking test for poping results.
func BenchmarkPop(b *testing.B) {
	b.StopTimer()
	logList = &LogList{}
	nodeCount := b.N
	nodeList := benchmarkSetup(nodeCount)
	logList.SetTolerance(int64(tolerance))
	for i := 0; i < nodeCount; i++ {
		err := logList.Push(nodeList[i].m, nodeList[i].r,
			nodeList[i].i, nodeList[i].t, nodeList[i].v)
		if err == LogListCollisionError {
			continue
		}
		if err != nil {
			b.Fail()
		}
	}
	b.StartTimer()
	for i := 0; i < nodeCount; i++ {
		_, err := logList.Pop()
		if err == LogListMatchNotFoundError {
			continue
		}
		if err == LogListEmptyError {
			break
		}
		if err != nil {
			b.Fail()
		}
	}
}
