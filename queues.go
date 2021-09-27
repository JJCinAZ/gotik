package gotik

import (
	"fmt"
	"strconv"
)

// Get queue tree entry by Name, including its children
func (c *Client) GetQueueTreeByName(name string) (QueueTree, error) {
	detail, err := c.Run("/queue/tree/print", "?name="+name)
	if err != nil {
		return QueueTree{}, err
	}
	tree, err := c.parseQueueTreePrint(detail)
	if err != nil || len(tree) == 0 {
		return QueueTree{}, err
	}
	return tree[0], err
}

// Given an Interface name or a Parent name, retrieves any queue and it's children attached
// to the interface or parent queue.
func (c *Client) GetQueueTree(parent string) ([]QueueTree, error) {
	detail, err := c.Run("/queue/tree/print", "?parent="+parent)
	if err != nil {
		return nil, err
	}
	return c.parseQueueTreePrint(detail)
}

func (c *Client) parseQueueTreePrint(detail *Reply) ([]QueueTree, error) {
	var queues []QueueTree
	queues = make([]QueueTree, 0)
	for _, re := range detail.Re {
		if parentQueue, err := parseQueueTreeEntry(re.Map); err == nil {
			child, err := c.Run("/queue/tree/print", "?parent="+parentQueue.Name)
			if err != nil {
				return nil, err
			}
			for _, reChild := range child.Re {
				if childQueue, err := parseQueueTreeEntry(reChild.Map); err == nil {
					parentQueue.Children = append(parentQueue.Children, childQueue)
				} else {
					return nil, err
				}
			}
			queues = append(queues, parentQueue)
		} else {
			return nil, err
		}
	}
	return queues, nil
}

func (c *Client) GetQueueTreeAll() ([]QueueTree, error) {
	var (
		queues []QueueTree
		err    error
		detail *Reply
	)
	queues = make([]QueueTree, 0)
	detail, err = c.Run("/queue/tree/print", "")
	if err != nil {
		return nil, err
	}
	for _, re := range detail.Re {
		if q, err := parseQueueTreeEntry(re.Map); err != nil {
			return nil, err
		} else {
			if parentQueue := findName(queues, q.Parent); parentQueue != nil {
				parentQueue.Children = append(parentQueue.Children, q)
			} else {
				if isRosId(q.Parent) {
					// We didn't find a queue Name matching this parent, so if the parent
					// is an ROS API ID like "*100014" then the parent is missing so change
					// parent to "unknown"
					q.Parent = "unknown"
				}
				queues = append(queues, q)
			}
		}
	}
	return queues, nil
}

func findName(tree []QueueTree, name string) *QueueTree {
	for i := 0; i < len(tree); i++ {
		if tree[i].Name == name {
			return &tree[i]
		}
		if q := findName(tree[i].Children, name); q != nil {
			return q
		}
	}
	return nil
}

func parseQueueTreeEntry(data map[string]string) (QueueTree, error) {
	/*
	 (string) (len=8) "limit-at": (string) (len=1) "0",
	 (string) (len=9) "max-limit": (string) (len=9) "100000000",
	 (string) (len=11) "burst-limit": (string) (len=1) "0",
	 (string) (len=7) "dropped": (string) (len=1) "0",
	 (string) (len=3) ".id": (string) (len=8) "*1000015",
	 (string) (len=8) "priority": (string) (len=1) "8",
	 (string) (len=15) "burst-threshold": (string) (len=1) "0",
	 (string) (len=11) "packet-rate": (string) (len=1) "0",
	 (string) (len=14) "queued-packets": (string) (len=1) "0",
	 (string) (len=4) "name": (string) (len=11) "Nextrio-Out",
	 (string) (len=6) "parent": (string) (len=12) "e1-NextrioFW",
	 (string) (len=10) "burst-time": (string) (len=2) "0s",
	 (string) (len=12) "queued-bytes": (string) (len=1) "0",
	 (string) (len=7) "invalid": (string) (len=4) "true",
	 (string) (len=8) "disabled": (string) (len=4) "true",
	 (string) (len=11) "packet-mark": (string) "",
	 (string) (len=5) "queue": (string) (len=7) "default",
	 (string) (len=11) "bucket-size": (string) (len=3) "0.1",
	 (string) (len=5) "bytes": (string) (len=1) "0",
	 (string) (len=7) "packets": (string) (len=1) "0",
	 (string) (len=4) "rate": (string) (len=1) "0"
	*/
	q := QueueTree{
		ID:             data[".id"],
		BucketSize:     parseFloat32(data["bucket-size"]),
		BurstLimit:     parseInt(data["burst-limit"]),
		BurstTime:      data["burst-time"],
		BurstThreshold: parseInt(data["burst-threshold"]),
		Disabled:       parseBool(data["disabled"]),
		Dynamic:        parseBool(data["dynamic"]),
		Invalid:        parseBool(data["invalid"]),
		Name:           data["name"],
		LimitAt:        parseInt(data["limit-at"]),
		MaxLimit:       parseInt(data["max-limit"]),
		PacketMark:     data["packet-mark"],
		Parent:         data["parent"],
		Priority:       parseInt(data["priority"]),
		Queue:          data["queue"],
		Comment:        data["comment"],
		Children:       nil,
	}
	return q, nil
}

// remove queue tree by name (removes the named queue tree entry and any of its children)
func (c *Client) RemoveQueueTreeByName(name string) error {
	tree, err := c.GetQueueTreeByName(name)
	if err != nil {
		return err
	}
	return c.RemoveQueueTree(tree, true)
}

// remove queue tree with or without children
func (c *Client) RemoveQueueTree(queue QueueTree, removeChildren bool) error {
	var removeIDS []string

	// children first
	if removeChildren {
		removeIDS = append(removeIDS, queue.getChildrenIDS()...)
	}
	// then parent
	removeIDS = append(removeIDS, queue.ID)

	for i := range removeIDS {
		_, err := c.Run("/queue/tree/remove", "=.id="+removeIDS[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// iterate through queue children for their id
func (q *QueueTree) getChildrenIDS() (childrenIDS []string) {
	for i := range q.Children {
		childrenIDS = append(childrenIDS, q.Children[i].ID)
	}
	return
}

func (c *Client) parseSimpleQueueEntry(data map[string]string) (SimpleQueue, error) {
	/* /queue/simple/print
	   (map[string]string) (len=30) {
	    (string) (len=6) "target": (string) (len=16) "e1-v195-Nextrio2",
	    (string) (len=11) "total-bytes": (string) (len=1) "0",
	    (string) (len=7) "dropped": (string) (len=3) "0/0",
	    (string) (len=12) "queued-bytes": (string) (len=3) "0/0",
	    (string) (len=4) "rate": (string) (len=3) "0/0",
	    (string) (len=7) "dynamic": (string) (len=5) "false",
	    (string) (len=12) "packet-marks": (string) "",
	    (string) (len=11) "burst-limit": (string) (len=3) "0/0",
	    (string) (len=8) "disabled": (string) (len=4) "true",
	    (string) (len=15) "burst-threshold": (string) (len=3) "0/0",
	    (string) (len=17) "total-packet-rate": (string) (len=1) "0",
	    (string) (len=18) "total-queued-bytes": (string) (len=1) "0",
	    (string) (len=20) "total-queued-packets": (string) (len=1) "0",
	    (string) (len=8) "limit-at": (string) (len=3) "0/0",
	    (string) (len=9) "max-limit": (string) (len=10) "20000000/0",
	    (string) (len=11) "bucket-size": (string) (len=7) "0.1/0.1",
	    (string) (len=7) "packets": (string) (len=3) "0/0",
	    (string) (len=10) "total-rate": (string) (len=1) "0",
	    (string) (len=3) ".id": (string) (len=2) "*4",
	    (string) (len=4) "name": (string) (len=16) "q-Nextrio2Upload",
	    (string) (len=6) "parent": (string) (len=4) "none",
	    (string) (len=13) "total-packets": (string) (len=1) "0",
	    (string) (len=13) "total-dropped": (string) (len=1) "0",
	    (string) (len=8) "priority": (string) (len=3) "8/8",
	    (string) (len=5) "queue": (string) (len=21) "default/default-small",
	    (string) (len=14) "queued-packets": (string) (len=3) "0/0",
	    (string) (len=7) "invalid": (string) (len=4) "true",
	    (string) (len=10) "burst-time": (string) (len=5) "0s/0s",
	    (string) (len=5) "bytes": (string) (len=3) "0/0",
	    (string) (len=11) "packet-rate": (string) (len=3) "0/0"
	*/
	q := SimpleQueue{
		ID:          data[".id"],
		Comment:     data["comment"],
		Dst:         data["dst"],
		Name:        data["name"],
		PacketMarks: data["packet-marks"],
		Parent:      data["parent"],
		Target:      data["target"],
		Time:        data["time"],
	}
	q.Disabled, _ = strconv.ParseBool(data["disabled"])
	q.Dynamic, _ = strconv.ParseBool(data["dynamic"])
	q.Invalid, _ = strconv.ParseBool(data["invalid"])
	if i, err := strconv.ParseInt(data["priority"], 10, 32); err == nil {
		q.Priority = int(i)
	}
	q.BucketSize = splitFloat32(data["bucket-size"])
	q.BurstLimit = splitInt(data["burst-limit"])
	q.BurstThreshold = splitInt(data["burst-threshold"])
	q.BurstTime = splitString2(data["burst-time"])
	q.LimitAt = splitInt(data["limit-at"])
	q.MaxLimit = splitInt(data["max-limit"])
	q.Queue = splitString2(data["queue"])
	return q, nil
}

// get all parent queue trees
func (c *Client) GetSimpleQueues(target string) ([]SimpleQueue, error) {
	var (
		queues []SimpleQueue
		err    error
		detail *Reply
	)
	if len(target) > 0 {
		detail, err = c.Run("/queue/simple/print", "?target="+target)
	} else {
		detail, err = c.Run("/queue/simple/print", "")
	}
	if err != nil {
		return nil, err
	}
	queues = make([]SimpleQueue, 0, 32)
	for _, re := range detail.Re {
		if q, err := c.parseSimpleQueueEntry(re.Map); err == nil {
			queues = append(queues, q)
		}
	}
	return queues, nil
}

// remove simple queue
func (c *Client) RemoveSimpleQueue(ID string) error {
	if isRosId(ID) {
		_, err := c.Run("/queue/simple/remove", "=.id="+ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add a single, new Queue Tree entry.
// Upon successful return, queue entries will have their ID's filled in
func (c *Client) AddQueueTree(queue *QueueTree) error {
	id, err := c.AddQueueTreeSingle(*queue)
	if err == nil {
		queue.ID = id
		for i := 0; i < len(queue.Children); i++ {
			if err = c.AddQueueTree(&queue.Children[i]); err != nil {
				return err
			}
		}
	}
	return err
}

// Add a single, new Queue Tree entry.
// Returns the ID of the new entry or error
func (c *Client) AddQueueTreeSingle(queue QueueTree) (string, error) {
	/*
		Children:       nil,
	*/
	parts := make([]string, 0, 10)
	parts = append(parts, "/queue/tree/add")
	parts = append(parts, fmt.Sprintf("=name=%s", queue.Name))
	parts = append(parts, fmt.Sprintf("=disabled=%t", queue.Disabled))
	parts = append(parts, fmt.Sprintf("=parent=%s", queue.Parent))
	parts = append(parts, fmt.Sprintf("=max-limit=%d", queue.MaxLimit))
	if queue.LimitAt > 0 {
		parts = append(parts, fmt.Sprintf("=limit-at=%d", queue.LimitAt))
	}
	if len(queue.PacketMark) > 0 {
		parts = append(parts, fmt.Sprintf("=packet-mark=%s", queue.PacketMark))
	}
	if queue.Priority > 0 && queue.Priority < 9 {
		parts = append(parts, fmt.Sprintf("=priority=%d", queue.Priority))
	}
	if len(queue.Queue) > 0 {
		parts = append(parts, fmt.Sprintf("=queue=%s", queue.Queue))
	}
	if len(queue.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", queue.Comment))
	}
	if queue.BucketSize != 0 {
		parts = append(parts, fmt.Sprintf("=bucket-size=%f", queue.BucketSize))
	}
	if queue.BurstLimit != 0 && queue.BurstThreshold != 0 {
		parts = append(parts, fmt.Sprintf("=burst-limit=%d", queue.BurstLimit))
		parts = append(parts, fmt.Sprintf("=burst-threshold=%d", queue.BurstThreshold))
	}
	if len(queue.BurstTime) > 0 {
		parts = append(parts, fmt.Sprintf("=burst-time=%s", queue.BurstTime))
	}
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}
