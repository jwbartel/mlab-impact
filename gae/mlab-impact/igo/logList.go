/* The ns package provides the net-score.org functionality.  */
package igo

/*This file defines the LogList data structure.

A LogList is a data structure that is used for matching up data to the server
logs associated to the requests that generated the logs.  

Items are stored in the list according to the http method (GET, POST, etc.), the
requested resource, the remote IP address, and finally the timestamp.  It is
assumed that the timestamp that was recorded by the handler will not exactly
match those that were recorded by the log API, and this is why we have a
tolerance parameter that will allow for fuzzy matching.

EXAMPLE LogList:

LogList -- GET ---------------------------------- POST
            |                                      |
            / ----------------------- /info/       /upload/
	    |                           |             |
	    192.168.1.1 -- 192.168.1.2  192.168.11.1  192.168.1.1
	         |              |            |             |
		 12:00 -- 1:00  12:01        12:05         12:01 -- 1:01
		   |       |      |            |             |       |
		   v0      v1     v2           v3            v4      v5

USAGE:
	see logList_test.go for examples.

TESTING:
	$go test logList_test.go logList.go
*/

import (
	"errors"
)

// Error values
var (
	LogListMatchNotFoundError = errors.New("LogList match not found")
	LogListCollisionError     = errors.New("LogList time collision")
	LogListEmptyError         = errors.New("LogList is empty")
	LogListUnknownError       = errors.New("LogList is broken")
)

// This struct is the basic unit of storage in a LogList
type Element struct {
	Value     interface{}
	prev      *Element
	next      *Element
	nextLevel *Element
}

// This struct forms the head of a LogList.
type LogList struct {
	front     *Element
	length    int64
	tolerance int64
}

// This method resets or clears a LogList.
func (l *LogList) Init() {
	l.front = nil
	l.length = 0
	l.tolerance = 0
}

// This method changes the tolerance value to allow fuzzy matching.
func (l *LogList) SetTolerance(t int64) {
	l.tolerance = t
}

/*This method returns the size of the LogList.

The size of a LogList is not the total number of elements in the list, but the
number of leaf elements.  These are the elements that contain the value passed
to the Push call.
*/
func (l *LogList) Size() int64 {
	return l.length
}

// This method tells us if the list is empty.
func (l *LogList) Empty() bool {
	return l.front == nil
}

/*This method adds an element to the LogList.

This works by searching for where the element is likely to be found and then
inserting the element in that place.
*/
func (l *LogList) Push(method, resource, ip string, time int64,
	value interface{}) error {

	// Setup the branch (unused sections will be discarded).
	valueLevel := &Element{
		Value: value,
	}
	timeLevel := &Element{
		Value:     time,
		nextLevel: valueLevel,
	}
	ipLevel := &Element{
		Value:     ip,
		nextLevel: timeLevel,
	}
	resourceLevel := &Element{
		Value:     resource,
		nextLevel: ipLevel,
	}
	methodLevel := &Element{
		Value:     method,
		nextLevel: resourceLevel,
	}

	// Assume the addition will be OK and take it back if it is not.
	l.length++

	mLevel := l.front
	var rLevel *Element
	var iLevel *Element
	var tLevel *Element

	// Empty lists are easy.
	if mLevel == nil {
		l.front = methodLevel
		return nil
	}

	// Find or set the method.
	for {
		if mLevel.Value.(string) == method {
			rLevel = mLevel.nextLevel
			break
		}
		if mLevel.next == nil {
			methodLevel.prev = mLevel
			mLevel.next = methodLevel
			return nil
		}
		mLevel = mLevel.next
	}
	// Find or set the resource.
	for {
		if rLevel.Value.(string) == resource {
			iLevel = rLevel.nextLevel
			break
		}
		if rLevel.next == nil {
			resourceLevel.prev = rLevel
			rLevel.next = resourceLevel
			return nil
		}
		rLevel = rLevel.next
	}
	// Find or set the ip.
	for {
		if iLevel.Value.(string) == ip {
			tLevel = iLevel.nextLevel
			break
		}
		if iLevel.next == nil {
			ipLevel.prev = iLevel
			iLevel.next = ipLevel
			return nil
		}
		iLevel = iLevel.next
	}
	// Find or set the time.  Results are in a sorted list of times.
	for {
		if tLevel.Value.(int64) == time {
			l.length--
			return LogListCollisionError
		}
		if tLevel.Value.(int64) > time {
			if tLevel.prev == nil {
				iLevel.nextLevel = timeLevel
				timeLevel.next = tLevel
				tLevel.prev = timeLevel
				return nil
			}
			timeLevel.next = tLevel
			timeLevel.prev = tLevel.prev
			tLevel.prev.next = timeLevel
			tLevel.prev = timeLevel
			return nil

		}
		if tLevel.Value.(int64) < time {
			if tLevel.next == nil {
				timeLevel.prev = tLevel
				tLevel.next = timeLevel
				return nil
			}
			tLevel = tLevel.next
		}
	}
	l.length--
	return LogListUnknownError
}

/*This method will find and return a value from a LogList.

Starting with the method, resource, and IP address there must be an exact match
for each level.  When we reach the time it is possible that we could find an
exact match, but if not we can look for a fuzzy match based on tolerance.

If the item is found, then it is removed from the LogList starting from the 
element node with the value, working its way back up to the head of the LogList.
Any empty branches are removed along the way.
*/
func (l *LogList) Match(method, resource, ip string,
	time int64) (interface{}, error) {

	var resourceLevel *Element
	var ipLevel *Element
	var timeLevel *Element
	var valueLevel *Element
	var value interface{}

	// Take care of the empty LogList case.
	methodLevel := l.front
	if methodLevel == nil {
		return nil, LogListMatchNotFoundError
	}

	// Find or set the method.
	for {
		if methodLevel.Value.(string) == method {
			resourceLevel = methodLevel.nextLevel
			break
		}
		if methodLevel.next == nil {
			return nil, LogListMatchNotFoundError
		}
		methodLevel = methodLevel.next
	}
	// Find or set the resource.
	for {
		if resourceLevel.Value.(string) == resource {
			ipLevel = resourceLevel.nextLevel
			break
		}
		if resourceLevel.next == nil {
			return nil, LogListMatchNotFoundError
		}
		resourceLevel = resourceLevel.next
	}
	// Find or set the ip.
	for {
		if ipLevel.Value.(string) == ip {
			timeLevel = ipLevel.nextLevel
			break
		}
		if ipLevel.next == nil {
			return nil, LogListMatchNotFoundError
		}
		ipLevel = ipLevel.next
	}
	// Find or set the time (exact match in sorted list).
	for {
		if timeLevel.Value.(int64) == time {
			valueLevel = timeLevel.nextLevel
			value = valueLevel.Value
			break
		}
		if timeLevel.next == nil ||
			timeLevel.next.Value.(int64) > time {
			// Start over with a fuzzy match.
			timeLevel = ipLevel.nextLevel
			break
		}
		timeLevel = timeLevel.next
	}

	if value == nil {
		// Find or set the time (fuzzy match in sorted list).
		for {
			if timeLevel.Value.(int64) > time {
				if timeLevel.prev == nil {
					da := timeLevel.Value.(int64) - time
					if da <= l.tolerance {
						valueLevel = timeLevel.nextLevel
						value = valueLevel.Value
						break
					}
					return nil, LogListMatchNotFoundError
				}
				da := timeLevel.Value.(int64) - time
				db := time - timeLevel.prev.Value.(int64)
				if da < db {
					if da < l.tolerance {
						valueLevel = timeLevel.nextLevel
						value = valueLevel.Value
						break
					}
					return nil, LogListMatchNotFoundError
				}
				if db < l.tolerance {
					valueLevel = timeLevel.nextLevel
					value = valueLevel.Value
					break
				}
				return nil, LogListMatchNotFoundError
			}
			if timeLevel.next == nil {
				da := time - timeLevel.Value.(int64)
				if da <= l.tolerance {
					valueLevel = timeLevel.nextLevel
					value = valueLevel.Value
					break
				}
				return nil, LogListMatchNotFoundError
			}
			timeLevel = timeLevel.next
		}
	}

	l.length--

	/*Splicing out the matched node.

	Apply these rules starting from the bottom up and continue until you
	reach a break or the head.

	Rules for removing nodes:
	 0.  valueLevel and timeLevel are always removed; continue
	 1.  If node.nextLevel == nil then delete(node); continue
	 2.a If node.prev == nil then parent.nextLevel = node.next;
	 2.b If node.next != nil then node.next.prev = node.prev; continue
	 3.a If node.prev != nil then node.prev.next = node.next;
	 3.b If node.next != nil then node.next.prev = node.prev; break
	*/

	// Rule 0
	if timeLevel.prev != nil {
		// Rule 3.a
		timeLevel.prev.next = timeLevel.next
		// Rule 3.b
		if timeLevel.next != nil {
			timeLevel.next.prev = timeLevel.prev
		}
		return value, nil
	}
	// Rule 2.a
	ipLevel.nextLevel = timeLevel.next
	// Rule 2.b
	if timeLevel.next != nil {
		timeLevel.next.prev = nil
	}

	if ipLevel.nextLevel != nil {
		return value, nil
	}
	// Rule 1
	if ipLevel.prev != nil {
		ipLevel.prev.next = ipLevel.next
		if ipLevel.next != nil {
			ipLevel.next.prev = ipLevel.prev
		}
		return value, nil
	}
	resourceLevel.nextLevel = ipLevel.next
	if ipLevel.next != nil {
		ipLevel.next.prev = nil
	}

	if resourceLevel.nextLevel != nil {
		return value, nil
	}
	if resourceLevel.prev != nil {
		resourceLevel.prev.next = resourceLevel.next
		if resourceLevel.next != nil {
			resourceLevel.next.prev = resourceLevel.prev
		}
		return value, nil
	}
	methodLevel.nextLevel = resourceLevel.next
	if resourceLevel.next != nil {
		resourceLevel.next.prev = nil
	}

	if methodLevel.nextLevel != nil {
		return value, nil
	}
	if methodLevel.prev != nil {
		methodLevel.prev.next = methodLevel.next
		if methodLevel.next != nil {
			methodLevel.next.prev = methodLevel.prev
		}
		return value, nil
	}
	l.front = methodLevel.next
	if methodLevel.next != nil {
		methodLevel.next.prev = nil
	}
	return value, nil
}

/*This method returns (and removes) the last element in a LogList.

It finds the last leaf in the LogList, and returns it.  Before completing it
works its way back up through the levels removing any empty branches along the
way.
*/
func (l *LogList) Pop() (interface{}, error) {
	var resourceLevel *Element
	var ipLevel *Element
	var timeLevel *Element
	var valueLevel *Element
	var value interface{}

	// Take care of the empty LogList case.
	methodLevel := l.front
	if methodLevel == nil {
		return nil, LogListEmptyError
	}

	// Find or set the method.
	for {
		if methodLevel.next == nil {
			resourceLevel = methodLevel.nextLevel
			break
		}
		methodLevel = methodLevel.next
	}
	// Find or set the resource.
	for {
		if resourceLevel.next == nil {
			ipLevel = resourceLevel.nextLevel
			break
		}
		resourceLevel = resourceLevel.next
	}
	// Find or set the ip.
	for {
		if ipLevel.next == nil {
			timeLevel = ipLevel.nextLevel
			break
		}
		ipLevel = ipLevel.next
	}
	// Find or set the time.
	for {
		if timeLevel.next == nil {
			valueLevel = timeLevel.nextLevel
			value = valueLevel.Value
			break
		}
		timeLevel = timeLevel.next
	}

	l.length--

	/*Splicing out the tail node.

	Apply these rules starting from the bottom up and continue until you
	reach a break or the head.

	Rules for removing nodes:
	 0.  valueLevel and timeLevel are always removed; continue
	 1.  If node.nextLevel == nil then delete(node); continue
	 2.  If node.prev != nil then node.prev.next = nil; break
	*/

	//Rule 0
	if timeLevel.prev != nil {
		//Rule 2
		timeLevel.prev.next = nil
		return value, nil
	}
	//Rule 1
	if ipLevel.prev != nil {
		ipLevel.prev.next = nil
		return value, nil
	}
	if resourceLevel.prev != nil {
		resourceLevel.prev.next = nil
		return value, nil
	}
	if methodLevel.prev != nil {
		methodLevel.prev.next = nil
	} else {
		// The above does not for for the head of the list.  The prev
		// for the first method is always nil, but would point to the
		// LogList.  This means that the list is empty.
		l.front = nil
	}
	return value, nil
}
