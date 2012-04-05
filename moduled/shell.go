package moduled

import (
	"errors"
	"bufio"
	"fmt"
	"io"
)

// type Shell provides a simple, interactive shell through which users can interact with your program. Includes unlimited history and line editing.
// Bug: Multiple line input doesn't work correctly with in-line editing.
type Shell struct {
	// Text to display as a prompt
	PS1     string
	
	// Full history, lines stored as raw input (not parsed)
	hist    [][]byte
	// Out current position in the history (used when "scrolling through" previous commands)
	histpos int
	
	// Buffer of characters on the current line
	linebuf []byte
	// Current position of the cursor
	linepos int
	
	rd     *bufio.Reader
	wr      io.Writer
}

// NewShell returns a Shell initialized with the given reader and writer (commonly attached to os.Stdin and os.Stdout respectively), as well as a simple default prompt.
func NewShell(rd io.Reader, wr io.Writer) *Shell {
	s := new(Shell)
	s.rd = bufio.NewReader(rd)
	s.wr = wr
	s.PS1 = "\x1B[32;1m >>> \x1B[m"
	return s
}

// parseTokens takes the current input line in the shell and parses in into distinct tokens. Returns any errors that occurred during parsing.
// TODO: parseTokens() should recognize commonly used syntax, like quoting tokens with spaces in them, and escaping spaces or newlines.
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

// Prompt presents the user with an interactive prompt where they can enter a command, and includes facilities for inline editing, such as removing characters with backspace and moving through the line with the left and right arrow keys.
// It also includes a endless history.
// Returns the parsed command string (command and arguments), and an error, if any.
func (s *Shell) Prompt() ([]string, error) {
	s.linebuf = nil
	s.histpos = -1
	s.linepos = -1
	fmt.Fprintf(s.wr, s.PS1)
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
		case 8, 127: // Backspace (Actually DEL)
			if s.linepos == 0 {
				s.overwrite(2)
				break
			}
			if len(s.linebuf) == 0 {
				s.overwrite(2)
				break
			}
			s.overwrite(3)
			if s.linepos == -1 {
				s.linebuf = s.linebuf[:len(s.linebuf)-1]
				break
			}
			
			s.overwritef(len(s.linebuf)-s.linepos+1)
			s.linebuf = append(s.linebuf[:s.linepos-1], s.linebuf[s.linepos:]...)
			s.linepos--
			if s.linepos < 0 {
				s.linepos = 0
			}
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
			if s.linepos == -1 {
				s.linebuf = append(s.linebuf, b)
			} else {
				s.overwritef(len(s.linebuf)-s.linepos)
				s.linebuf = append(s.linebuf[:s.linepos], append([]byte{b}, s.linebuf[s.linepos:]...)...)
				s.linepos++
			}
		}
	}
	panic("not reached!")
}

// readEscSeq reads an escape sequence (used for arrow keys, function keys, any many other functions not directly related to outputting characters) in standard format from the terminal, and returns the sequence read along with an error, if any. See http://en.wikipedia.org/wiki/ANSI_escape_code for more information on the structure of these escape codes
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

// bksp removes the current character from the terminal (overwrites it with a space), and moves the cursor back one space.
func (s *Shell) bksp(num int) {
	for i := 0; i < num; i++ {
		s.bkspb(' ')
	}
}

// bkspb overwrites the current character with one of your choosing, moving the cursor back one space.
func (s *Shell) bkspb(b byte) {
	fmt.Fprintf(s.wr, "%c%c%c", 8, b, 8)
}

// charat returns the character at the given offset from the current cursor position. It is used to replace portions of the line when their contents are overwritten with escape sequences (such as when the arrow keys are pressed), or to overwrite the remaining portion of the line when backspace is pressed.
func (s *Shell) charat(offset int) byte {
	if s.linepos < 0 || s.linepos + offset >= len(s.linebuf) {
		return ' '
	}
	return s.linebuf[s.linepos + offset]
}

// overwritef overwrites num characters in front of the cursor with values given by charat(). The cursor position is the same before and after the call (although it moves during the call)
func (s *Shell) overwritef(num int) {
	for i := 0; i < num; i++ {
		fmt.Fprintf(s.wr, "%c", s.charat(i))
	}
	for i := 0; i < num; i++ {
		fmt.Fprintf(s.wr, "%c", 8)
	}
}

// overwrite overwrites num characters behind the cursor with values given by charat(). The cursor position will move back by num characters after the call.
func (s *Shell) overwrite(num int) {
	for i := num - 1; i >= 0; i-- {
		s.bkspb(s.charat(i))
	}
}

// handleEscSeq implements rudementry handling of a few of the most useful escape sequences, including the arrow keys. It's current behaviour is:
// SEQ_UP: Moves to the previous history item, saving any currently typed text to the end of the history buffer
// SEQ_DOWN: Moves to the next history item
// SEQ_LEFT: Moves the cursor one position to the left
// SEQ_RIGHT: Moves the cursor one position to the right
func (s *Shell) handleEscSeq(seq Seq) error {
	if seq == SEQ_UP || seq == SEQ_DOWN {
		if s.histpos == -1 {
			s.hist = append(s.hist, s.linebuf)
			s.histpos = len(s.hist)-1
		}
	}
	if seq == SEQ_LEFT || seq == SEQ_RIGHT {
		if s.linepos == -1 {
			s.linepos = len(s.linebuf)
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
		if s.histpos > len(s.hist) - 1 {
			s.histpos = len(s.hist) - 1
			s.bksp(4)
			return nil
		}
	case SEQ_LEFT:
		s.overwrite(4)
		s.linepos--
		if s.linepos < 0 {
			s.linepos = 0
			break
		}
		fmt.Fprintf(s.wr, "%c[D", 27)
	case SEQ_RIGHT:
		s.overwrite(4)
		s.linepos++
		if s.linepos > len(s.linebuf) {
			s.linepos = len(s.linebuf)
			break
		}
		fmt.Fprintf(s.wr, "%c[C", 27)
	default:
		fmt.Printf("Read escape sequence of %d", seq)
		return fmt.Errorf("Read unsupported escape sequence of %d ('%c')", seq, seq)
	}
	if seq == SEQ_UP || seq == SEQ_DOWN {
		s.bksp(len(s.linebuf) + 4)
		s.linebuf = make([]byte, len(s.hist[s.histpos]))
		copy(s.linebuf, s.hist[s.histpos])
		fmt.Fprintf(s.wr, "%s", s.linebuf)
	}
	return nil
}