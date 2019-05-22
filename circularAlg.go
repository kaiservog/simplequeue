package main

/* circular Q is the circular queue struct */
type circularAlg struct {
	rep repository
}

func newCircularAlg(rep repository) *circularAlg {
	q := circularAlg{rep: rep}
	return &q
}

func incrementOrReset(idx, max int) int {
	idx++

	if idx >= max {
		return 0
	}

	return idx
}

func (ca *circularAlg) put(m, id string) error {
	q, err := ca.rep.queue(id)
	if err != nil {
		return err
	}

	ca.rep.putMessage(q.idxI, m, id)

	nidxO := q.idxO
	if q.idxI == q.idxO {
		_, err := ca.rep.getMessage(q.idxO+1, id)
		if err != nil && err.Error() != "empty" {
			nidxO = incrementOrReset(q.idxO, q.depth)
		}
	}

	nidxI := incrementOrReset(q.idxI, q.depth)
	err = ca.rep.updateIdx(nidxI, nidxO, id)

	if err != nil {
		return err
	}

	return nil
}

func (ca *circularAlg) get(id string) (string, error) {
	q, err := ca.rep.queue(id)
	if err != nil {
		return "", err
	}

	m, err := ca.rep.getMessage(q.idxO, id)

	if err != nil {
		return "", err
	}

	err = ca.rep.deleteMessage(q.idxO, id)
	if err != nil {
		return "", err
	}

	nidxO := incrementOrReset(q.idxO, q.depth)
	ca.rep.updateIdx(-1, nidxO, id)

	return m, nil
}
