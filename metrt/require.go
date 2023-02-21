package metrt

type Require struct {
	assert Assert `chaingen:"pre(*)=r.assert.filter.T.t.Helper(),wrap(*)=require,-getValue"`
}

func (r Require) require(ok bool) {
	r.assert.filter.T.t.Helper()
	if ok {
		return
	}
	r.assert.filter.T.t.FailNow()
}

type RequireGroup[K comparable] struct {
	assert AssertGroup[K] `chaingen:"pre(*)=r.assert.group.T.t.Helper(),wrap(*)=require"`
}

func (r RequireGroup[K]) require(ok bool) {
	r.assert.group.T.t.Helper()
	if ok {
		return
	}
	r.assert.group.T.t.FailNow()
}
