package factorgraph

type Variable struct {
	messages *Gaussian
}

func (v *Variable) set(a *Variable) {
	delta :=  v.delta(a)
}

func (v *Variable) delta(a *Variable) {
	pi_delta :=  v.delta(a)
}

def set(self, val):
	delta = self.delta(val)
	self.pi, self.tau = val.pi, val.tau
	return delta

def delta(self, other):
	pi_delta = abs(self.pi - other.pi)
	if pi_delta == inf:
		return 0.
	return max(abs(self.tau - other.tau), math.sqrt(pi_delta))

def update_message(self, factor, pi=0, tau=0, message=None):
	message = message or Gaussian(pi=pi, tau=tau)
	old_message, self[factor] = self[factor], message
	return self.set(self / old_message * message)

def update_value(self, factor, pi=0, tau=0, value=None):
	value = value or Gaussian(pi=pi, tau=tau)
	old_message = self[factor]
	self[factor] = value * old_message / self
	return self.set(value)

def __getitem__(self, factor):
	return self.messages[factor]

def __setitem__(self, factor, message):
	self.messages[factor] = message

def __repr__(self):
	args = (type(self).__name__, super(Variable, self).__repr__(),
			len(self.messages), '' if len(self.messages) == 1 else 's')
	return '<%s %s with %d connection%s>' % args
