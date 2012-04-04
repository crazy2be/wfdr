package moduled

import (
	"errors"
	"bufio"
	"fmt"
	"io"
)

type Shell struct {
	hist    [][]byte
	histpos int
	
	// Buffer of characters on the current line
	linebuf []byte
	// Current position of the cursor
	pos     int
	
	// Where to read commands from
	rd     *bufio.Reader
	wr      io.Writer
}

func NewShell(rd io.Reader, wr io.Writer) *Shell {
	s := new(Shell)
	s.rd = bufio.NewReader(rd)
	s.wr = wr
	return s
}

func (s *Shell) parseTokens() ([]string, error) {
	tokens := make([]string, 0)
	buf := make([]byte, 0)
	
	for _, b := range s.linebuf {
		switch b {
		case ' ', '\t':
			if buf != nil {
				tokens = append(tokens, string(buf))
				buf = nil
			}
		default:
			buf = append(buf, b)
		}
	}
	tokens = append(tokens, string(buf))
	return tokens, nil
}

// ReadCommand presents the user with an interactive prompt where they can enter a command, and includes facilities for backspace, and history. Returns the parsed command string (command and arguments), and an error, if any.
func (s *Shell) ReadCommand() ([]string, error) {
	s.linebuf = nil
	s.histpos = -1
	for {
		b, err := s.rd.ReadByte()
		if err != nil {
			return nil, err
		}
		
		switch b {
		case '\n':
			if s.histpos == -1 {
				s.hist = append(s.hist, s.linebuf)
			} else {
				s.hist[len(s.hist) - 1] = s.linebuf
			}
			return s.parseTokens()
		case 127: // Backspace (Actually DEL)
			s.bksp(3)
			if len(s.linebuf) == 0 || s.linebuf == nil {
				break
			}
			s.linebuf = s.linebuf[:len(s.linebuf)-1]
		case 27: // ESC
			seq, err := s.readEscSeq()
			if err != nil {
				return nil, err
			}
			err = s.handleEscSeq(seq)
			if err != nil {
				return nil, err
			}
		default:
			s.linebuf = append(s.linebuf, b)
		}
	}
	panic("not reached!")
}

func (s *Shell) readEscSeq() (Seq, error) {
	b, err := s.rd.ReadByte()
	if err != nil {
		return SEQ_NONE, err
	}
	if b != '[' {
		return SEQ_NONE, errors.New("Expected '[' after ESC")
	}
	
	pmc := make([]byte, 0)
	for {
		b, err = s.rd.ReadByte()
		if err != nil {
			return SEQ_NONE, err
		}
		if (b < 22) {
			return SEQ_NONE, fmt.Errorf("Unexpected character %d ('%c')", b, b)
		} else if (b <= 47) {
			pmc = append(pmc, b)
		} else if (b <= 57) {
			// Number
			return SEQ_NONE, fmt.Errorf("Numbers in escape sequences not yet supported!")
		} else if (b <= 63) {
			return SEQ_NONE, fmt.Errorf("Unexpected character %d ('%c')", b, b)
		} else if (b <= 126) {
			return Seq(b), nil
		} else {
			return SEQ_NONE, fmt.Errorf("Unexpected character %d ('%c')", b, b)
		}
	}
	panic("Not reached!")
}

func (s *Shell) bksp(num int) {
	for i := 0; i < num; i++ {
		fmt.Fprintf(s.wr, "%c %c", 8, 8)
	}
}

func (s *Shell) handleEscSeq(seq Seq) error {
	if seq == SEQ_UP || seq == SEQ_DOWN {
		if s.histpos == -1 {
			s.hist = append(s.hist, s.linebuf)
			s.histpos = len(s.hist)-1
		}
	}
	switch (seq) {
	case SEQ_UP:
		s.histpos--
		if s.histpos < 0 {
			s.histpos = 0
			s.bksp(4)
			return nil
		}
	case SEQ_DOWN:
		s.histpos++
		if s.histpos >= len(s.hist) - 1 {
			s.histpos = len(s.hist) - 1
			s.bksp(4)
			return nil
		}
	default:
		fmt.Printf("Read escape sequence of %d", seq)
		return fmt.Errorf("Read unsupported escape sequence of %d ('%c')", seq, seq)
	}
	if seq == SEQ_UP || seq == SEQ_DOWN {
		s.bksp(len(s.linebuf) + 4)
		s.linebuf = s.hist[s.histpos]
		fmt.Fprintf(s.wr, "%s", s.linebuf)
	}
	return nil
}

func (c *Conn) InterpretCommand(args []string) error {
	if args == nil || len(args) < 1 {
		return errors.New("No command provided!")
	}
	cmd := args[0]
	args = args[1:]
	
	switch (cmd) {
		case "start":
			for _, module := range args {
				err := c.Start(module)
				if err != nil {
					return err
				}
			}
		case "stop":
			for _, module := range args {
				if len(args) < 1 {
					return errors.New("Module name required!")
				}
				return c.Stop(module)
			}
		case "restart":
			for _, module := range args {
				if len(args) < 1 {
					return errors.New("Module name required!")
				}
				return c.Restart(module)
			}
	}
	
	return nil
}