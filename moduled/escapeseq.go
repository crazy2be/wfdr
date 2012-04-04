package moduled

type Seq byte

var (
	SEQ_NONE  = Seq(0)
	SEQ_UP    = Seq('A')
	SEQ_DOWN  = Seq('B')
	SEQ_RIGHT = Seq('C')
	SEQ_LEFT  = Seq('D')
)
