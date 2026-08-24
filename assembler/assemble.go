package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"corwar/globalvar"
	"corwar/opcode"
)

type instr struct {
	op    *opcode.Op
	args  []string
	types []opcode.ArgType
	pc    int
}

func Assemble(path string) error {
	if !strings.HasSuffix(path, ".s") {
		return fmt.Errorf("%s: wrong file extension", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := cleanLines(string(data))

	name, comment, body, err := readHeader(lines)
	if err != nil {
		return fmt.Errorf("%s: %v", path, err)
	}

	instrs, labels, err := parseBody(body)
	if err != nil {
		return fmt.Errorf("%s: %v", path, err)
	}

	code, err := encodeAll(instrs, labels)
	if err != nil {
		return fmt.Errorf("%s: %v", path, err)
	}
	if len(code) > globalvar.ChampMaxSize {
		return fmt.Errorf("%s: program too big (%d > %d)", path, len(code), globalvar.ChampMaxSize)
	}

	out := strings.TrimSuffix(path, ".s") + ".cor"
	return writeCor(out, name, comment, code)
}

func cleanLines(text string) []string {
	raw := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	lines := make([]string, len(raw))
	for i, l := range raw {
		if j := strings.IndexByte(l, '#'); j >= 0 {
			l = l[:j]
		}
		lines[i] = l
	}
	return lines
}

func readHeader(lines []string) (name, comment string, body []string, err error) {
	i := 0
	for ; i < len(lines); i++ {
		s := strings.TrimSpace(lines[i])
		if s == "" {
			continue
		}
		if !strings.HasPrefix(s, ".") {
			break
		}
		key, value := headerValue(s)
		switch key {
		case ".name":
			name = value
		case ".description":
			comment = value
		default:
			return "", "", nil, fmt.Errorf("unknown header field %q", key)
		}
	}
	if name == "" || comment == "" {
		return "", "", nil, fmt.Errorf(`need .name "..." and .description "..."`)
	}
	if len(name) > globalvar.ProgNameLen || len(comment) > globalvar.CommentLen {
		return "", "", nil, fmt.Errorf("name or comment too long")
	}
	return name, comment, lines[i:], nil
}

func headerValue(s string) (string, string) {
	parts := strings.SplitN(s, `"`, 3)
	if len(parts) != 3 || strings.TrimSpace(parts[2]) != "" {
		return parts[0], ""
	}
	return strings.TrimSpace(parts[0]), parts[1]
}

func parseBody(lines []string) ([]instr, map[string]int, error) {
	var out []instr
	labels := map[string]int{}
	pc := 0

	for i, line := range lines {
		text := strings.TrimSpace(line)
		if text == "" {
			continue
		}
		if label, ok := takeLabel(text); ok {
			if _, dup := labels[label]; dup {
				return nil, nil, fmt.Errorf("line %d: duplicate label %q", i+1, label)
			}
			labels[label] = pc
			text = restAfterLabel(text)
			if text == "" {
				continue
			}
		}
		in, err := parseInstr(text, pc)
		if err != nil {
			return nil, nil, fmt.Errorf("line %d: %v", i+1, err)
		}
		out = append(out, in)
		pc += instrSize(in)
	}
	return out, labels, nil
}

func takeLabel(s string) (string, bool) {
	n := 0
	for n < len(s) && isLabelChar(s[n]) {
		n++
	}
	if n == 0 || n >= len(s) || s[n] != ':' {
		return "", false
	}
	if n+1 < len(s) && s[n+1] != ' ' && s[n+1] != '\t' {
		return "", false
	}
	return s[:n], true
}

func isLabelChar(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}

func restAfterLabel(s string) string {
	j := strings.IndexByte(s, ':')
	s = s[j+1:]
	return strings.TrimSpace(s)
}

func parseInstr(text string, pc int) (instr, error) {
	var in instr
	name := text
	argText := ""
	if j := strings.IndexAny(text, " \t"); j >= 0 {
		name, argText = text[:j], text[j:]
	}
	op := opcode.FindOp(name)
	if op == nil {
		return in, fmt.Errorf("unknown instruction %q", name)
	}

	args := []string{}
	if argText != "" {
		for _, a := range strings.Split(argText, ",") {
			a = strings.TrimSpace(a)
			if a == "" || strings.ContainsAny(a, " \t") {
				return in, fmt.Errorf("bad arguments for %s", name)
			}
			args = append(args, a)
		}
	}
	if len(args) != op.NbParams {
		return in, fmt.Errorf("%s needs %d args, got %d", name, op.NbParams, len(args))
	}

	types := make([]opcode.ArgType, len(args))
	for i, a := range args {
		t, err := checkArg(a, op.Types[i])
		if err != nil {
			return in, err
		}
		types[i] = t
	}
	return instr{op: op, args: args, types: types, pc: pc}, nil
}

func checkArg(arg string, allowed []opcode.ArgType) (opcode.ArgType, error) {
	t := opcode.ArgIndirect
	if len(arg) >= 2 && arg[0] == 'r' {
		t = opcode.ArgRegister
	} else if arg[0] == '%' {
		t = opcode.ArgDirect
	}
	ok := false
	for _, a := range allowed {
		if t == a {
			ok = true
		}
	}
	if !ok {
		return t, fmt.Errorf("arg %q has wrong type here", arg)
	}

	if t == opcode.ArgRegister {
		n, err := strconv.Atoi(arg[1:])
		if err != nil || n < 1 || n > globalvar.RegNumber {
			return t, fmt.Errorf("bad register %q", arg)
		}
		return t, nil
	}

	v := strings.TrimPrefix(arg, "%")
	if v == "" {
		return t, fmt.Errorf("empty arg %q", arg)
	}
	if v[0] == ':' {
		return t, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n < -2147483648 || n > 2147483647 {
		return t, fmt.Errorf("bad number %q", v)
	}
	return t, nil
}

func instrSize(in instr) int {
	size := 1
	if in.op.HasPcode {
		size++
	}
	for _, t := range in.types {
		switch t {
		case opcode.ArgRegister:
			size += 1
		case opcode.ArgDirect:
			if in.op.HasIdx {
				size += 2
			} else {
				size += 4
			}
		default:
			size += 2
		}
	}
	return size
}

func encodeAll(instrs []instr, labels map[string]int) ([]byte, error) {
	code := []byte{}
	for _, in := range instrs {
		buf := []byte{in.op.Opcode}
		if in.op.HasPcode {
			var p byte
			for i, t := range in.types {
				p |= opcode.ArgToPcode(t) << (6 - 2*uint(i))
			}
			buf = append(buf, p)
		}
		for i, a := range in.args {
			b, err := encodeArg(in.types[i], a, in.pc, labels, in.op.HasIdx)
			if err != nil {
				return nil, err
			}
			buf = append(buf, b...)
		}
		code = append(code, buf...)
	}
	return code, nil
}

func encodeArg(t opcode.ArgType, arg string, pc int, labels map[string]int, hasIdx bool) ([]byte, error) {
	switch t {
	case opcode.ArgRegister:
		n, _ := strconv.Atoi(arg[1:])
		return []byte{byte(n)}, nil
	case opcode.ArgDirect:
		size := 4
		if hasIdx {
			size = 2
		}
		return valueBytes(strings.TrimPrefix(arg, "%"), pc, labels, size)
	default:
		return valueBytes(arg, pc, labels, 2)
	}
}

func valueBytes(v string, pc int, labels map[string]int, size int) ([]byte, error) {
	n := int64(0)
	if v[0] == ':' {
		addr, ok := labels[v[1:]]
		if !ok {
			return nil, fmt.Errorf("unknown label %q", v[1:])
		}
		n = int64(addr - pc)
	} else {
		n, _ = strconv.ParseInt(v, 10, 64)
	}
	out := make([]byte, size)
	if size == 2 {
		binary.BigEndian.PutUint16(out, uint16(n))
	} else {
		binary.BigEndian.PutUint32(out, uint32(n))
	}
	return out, nil
}

func writeCor(path, name, comment string, code []byte) error {
	targetPath := filepath.Join("players", path)
	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	head := make([]byte, globalvar.HeaderSize)
	binary.BigEndian.PutUint32(head[0:], globalvar.Magic)
	copy(head[4:], name)
	binary.BigEndian.PutUint32(head[136:], uint32(len(code)))
	copy(head[140:], comment)
	if _, err := f.Write(head); err != nil {
		return err
	}
	_, err = f.Write(code)
	return err
}
