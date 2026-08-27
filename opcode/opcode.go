package opcode

type ArgType byte

const (
	ArgNone     ArgType = iota
	ArgRegister
	ArgDirect
	ArgIndirect
)

type Op struct {
	Name     string
	Opcode   byte
	Cycles   int
	HasPcode bool
	HasIdx   bool
	NbParams int
	Types    [3][]ArgType
}

var Ops = [17]Op{
	{},
	{Name: "live", Opcode: 1, Cycles: 10, NbParams: 1, Types: [3][]ArgType{{ArgDirect}, nil, nil}},
	{Name: "ld", Opcode: 2, Cycles: 5, HasPcode: true, NbParams: 2, Types: [3][]ArgType{{ArgIndirect, ArgDirect}, {ArgRegister}, nil}},
	{Name: "st", Opcode: 3, Cycles: 5, HasPcode: true, NbParams: 2, Types: [3][]ArgType{{ArgRegister}, {ArgRegister, ArgIndirect}, nil}},
	{Name: "add", Opcode: 4, Cycles: 10, HasPcode: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister}, {ArgRegister}, {ArgRegister}}},
	{Name: "sub", Opcode: 5, Cycles: 10, HasPcode: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister}, {ArgRegister}, {ArgRegister}}},
	{Name: "and", Opcode: 6, Cycles: 6, HasPcode: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister}}},
	{Name: "or", Opcode: 7, Cycles: 6, HasPcode: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister}}},
	{Name: "xor", Opcode: 8, Cycles: 6, HasPcode: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister}}},
	{Name: "zjmp", Opcode: 9, Cycles: 20, HasIdx: true, NbParams: 1, Types: [3][]ArgType{{ArgDirect}, nil, nil}},
	{Name: "ldi", Opcode: 10, Cycles: 25, HasPcode: true, HasIdx: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister, ArgDirect}, {ArgRegister}}},
	{Name: "sti", Opcode: 11, Cycles: 25, HasPcode: true, HasIdx: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister}, {ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister, ArgDirect}}},
	{Name: "fork", Opcode: 12, Cycles: 800, HasIdx: true, NbParams: 1, Types: [3][]ArgType{{ArgDirect}, nil, nil}},
	{Name: "lld", Opcode: 13, Cycles: 10, HasPcode: true, NbParams: 2, Types: [3][]ArgType{{ArgIndirect, ArgDirect}, {ArgRegister}, nil}},
	{Name: "lldi", Opcode: 14, Cycles: 50, HasPcode: true, HasIdx: true, NbParams: 3, Types: [3][]ArgType{{ArgRegister, ArgIndirect, ArgDirect}, {ArgRegister, ArgDirect}, {ArgRegister}}},
	{Name: "lfork", Opcode: 15, Cycles: 1000, HasIdx: true, NbParams: 1, Types: [3][]ArgType{{ArgDirect}, nil, nil}},
	{Name: "nop", Opcode: 16, Cycles: 2, HasPcode: true, NbParams: 1, Types: [3][]ArgType{{ArgRegister}, nil, nil}},
}

func ArgToPcode(t ArgType) byte {
	switch t {
	case ArgRegister:
		return 1
	case ArgDirect:
		return 2
	case ArgIndirect:
		return 3
	}
	return 0
}

func PcodeToArg(p byte) ArgType {
	switch p {
	case 1:
		return ArgRegister
	case 2:
		return ArgDirect
	case 3:
		return ArgIndirect
	}
	return ArgNone
}

func ClassifyArg(s string) ArgType {
	if len(s) == 0 {
		return ArgNone
	}
	switch s[0] {
	case 'r':
		return ArgRegister
	case '%':
		return ArgDirect
	default:
		return ArgIndirect
	}
}

func FindOp(name string) *Op {
	for i := range Ops {
		if Ops[i].Name == name {
			return &Ops[i]
		}
	}
	return nil
}
