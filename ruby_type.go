package robjspace

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
	Mask
	Match
	Module
	Nil
	Node
	None
	Object
	Rational
	Regexp
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
	case Mask:
		return "MATCH"
	case Match:
		return "MODULE"
	case Module:
		return "NIL"
	case Nil:
		return "NODE"
	case Node:
		return "NONE"
	case None:
		return "OBJECT"
	case Object:
		return "RATIONAL"
	case Rational:
		return "REGEXP"
	case Regexp:
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
		return Mask, nil
	case "MODULE":
		return Match, nil
	case "NIL":
		return Module, nil
	case "NODE":
		return Nil, nil
	case "NONE":
		return Node, nil
	case "OBJECT":
		return None, nil
	case "RATIONAL":
		return Object, nil
	case "REGEXP":
		return Rational, nil
	case "ROOT":
		return Regexp, nil
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
