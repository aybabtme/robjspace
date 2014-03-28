package rubyobj

import (
	"fmt"
)

type flagType uint64

const (
	frozen flagType = 1 << iota
	broken
	fstring
	gcMarked
	gcOld
	gcWbProtected
	shared
	embedded
)

func flagsFromSchema(schema *objectSchema) flagType {
	flag := flagType(0)
	if schema.Broken {
		flag |= broken
	}

	if schema.Frozen {
		flag |= frozen
	}

	if schema.Embedded {
		flag |= embedded
	}

	if schema.Broken {
		flag |= broken
	}

	if schema.Fstring {
		flag |= fstring
	}

	if schema.Shared {
		flag |= shared
	}

	if schema.Flags.Marked {
		flag |= gcMarked
	}

	if schema.Flags.Old {
		flag |= gcOld
	}

	if schema.Flags.WbProtected {
		flag |= gcWbProtected
	}

	return flag
}

type RubyType uint8

const (
	Array RubyType = iota
	Bignum
	Class
	Complex
	Data
	False
	File
	Fixnum
	Float
	Hash
	Iclass
	Match
	Module
	Nil
	Node
	None
	Object
	Rational
	Regexp
	Root
	String
	Struct
	Symbol
	True
	Undef
	Zombie
)

// Name is the string repesentation of this type in a ObjectSpace dump.
func (rt RubyType) Name() string {
	switch rt {
	case Array:
		return "ARRAY"
	case Bignum:
		return "BIGNUM"
	case Class:
		return "CLASS"
	case Complex:
		return "COMPLEX"
	case Data:
		return "DATA"
	case False:
		return "FALSE"
	case File:
		return "FILE"
	case Fixnum:
		return "FIXNUM"
	case Float:
		return "FLOAT"
	case Hash:
		return "HASH"
	case Iclass:
		return "ICLASS"
	case Match:
		return "MATCH"
	case Module:
		return "MODULE"
	case Nil:
		return "NIL"
	case Node:
		return "NODE"
	case None:
		return "NONE"
	case Object:
		return "OBJECT"
	case Rational:
		return "RATIONAL"
	case Regexp:
		return "REGEXP"
	case Root:
		return "ROOT"
	case String:
		return "STRING"
	case Struct:
		return "STRUCT"
	case Symbol:
		return "SYMBOL"
	case True:
		return "TRUE"
	case Undef:
		return "UNDEF"
	case Zombie:
		return "ZOMBIE"
	}
	panic(fmt.Sprintf("Missing RubyType '%T' in switch. This is a bug, please report it.", rt))
}

func typeFromName(typename string) (RubyType, error) {
	switch typename {
	case "ARRAY":
		return Array, nil
	case "BIGNUM":
		return Bignum, nil
	case "CLASS":
		return Class, nil
	case "COMPLEX":
		return Complex, nil
	case "DATA":
		return Data, nil
	case "FALSE":
		return False, nil
	case "FILE":
		return File, nil
	case "FIXNUM":
		return Fixnum, nil
	case "FLOAT":
		return Float, nil
	case "HASH":
		return Hash, nil
	case "ICLASS":
		return Iclass, nil
	case "MATCH":
		return Match, nil
	case "MODULE":
		return Module, nil
	case "NIL":
		return Nil, nil
	case "NODE":
		return Node, nil
	case "NONE":
		return None, nil
	case "OBJECT":
		return Object, nil
	case "RATIONAL":
		return Rational, nil
	case "REGEXP":
		return Regexp, nil
	case "ROOT":
		return Root, nil
	case "STRING":
		return String, nil
	case "STRUCT":
		return Struct, nil
	case "SYMBOL":
		return Symbol, nil
	case "TRUE":
		return True, nil
	case "UNDEF":
		return Undef, nil
	case "ZOMBIE":
		return Zombie, nil
	}
	return None, fmt.Errorf("not a Ruby type: %s", typename)
}
