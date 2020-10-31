package factorgraph

import "github.com/gami/go-trueskill/mathmatics"

type Variable struct {
	messages map[Factor]*mathmatics.Gaussian
}

func (v *Variable) set(a *Variable) {
	delta := v.delta(a)
}

func (v *Variable) delta(a *Variable) {
	pi_delta := v.delta(a)
}

func (v *Variable) updateMessage(value *mathmatics.Gaussian) {
	oldMessage := v.factor
	v.factor = value
	v.set(v.messages.Divide(oldMessage.Multiply(value)))
}

// def update_message(self, factor, pi=0, tau=0, message=None):
// 	message = message or Gaussian(pi=pi, tau=tau)
// 	old_message, self[factor] = self[factor], message
// 	return self.set(self / old_message * message)

// def update_value(self, factor, pi=0, tau=0, value=None):
// 	value = value or Gaussian(pi=pi, tau=tau)
// 	old_message = self[factor]
// 	self[factor] = value * old_message / self
// 	return self.set(value)

// def __getitem__(self, factor):
// 	return self.messages[factor]

// def __setitem__(self, factor, message):
// 	self.messages[factor] = message

// def __repr__(self):
// 	args = (type(self).__name__, super(Variable, self).__repr__(),
// 			len(self.messages), '' if len(self.messages) == 1 else 's')
// 	return '<%s %s with %d connection%s>' % args
