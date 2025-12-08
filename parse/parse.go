// package parse is a small parser combinator library.
package parse

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

// ParseResult the result from running a parser function.
type ParseResult[A any] struct {
	result    A
	remainder []byte
}

// Parser is a function that extracts a value and moves the input along.
// TODO: might need to be a struct for better error messages.
type Parser[A any] func([]byte) (ParseResult[A], error)

// Run runs the parser on an input, returning the result. It is not an error if
// the parser does not consume the entire input, but anything left over is not
// returned.
func Run[A any](p Parser[A], input []byte) (A, error) {
	result, err := p(input)
	if err != nil {
		var a A
		return a, err
	}
	return result.result, nil
}

// Apply returns a new paresr that runs the first parser, then applies the
// provided mapping function to its results (if it succeeds).
func Apply[A, B any](p Parser[A], f func(A) B) Parser[B] {
	return func(input []byte) (ParseResult[B], error) {
		r, err := p(input)
		if err != nil {
			return ParseResult[B]{}, err
		}
		return ParseResult[B]{
			result:    f(r.result),
			remainder: r.remainder,
		}, nil
	}
}

// Many returns a parser that applies the given parser repeatedly until it
// fails. This may be zero times, so the returned parser itself never fails.
func Many[A any](p Parser[A]) Parser[[]A] {
	return func(input []byte) (ParseResult[[]A], error) {
		var results []A
		for {
			r, err := p(input)
			if err != nil {
				break
			}
			results = append(results, r.result)
			input = r.remainder
		}
		return ParseResult[[]A]{
			result:    results,
			remainder: input,
		}, nil
	}
}

// Some returns a parser that applies the given parser as many times as it can,
// and errors if it can not apply it at least once.
func Some[A any](p Parser[A]) Parser[[]A] {
	return func(input []byte) (ParseResult[[]A], error) {
		r, err := Many(p)(input)
		if err != nil {
			// Not actual possible.
			return ParseResult[[]A]{}, err
		}
		if len(r.result) == 0 {
			return ParseResult[[]A]{}, errors.New("Some: could not parse at least once")
		}
		return r, nil
	}
}

// Between returns a parser that runs the three provided parsers in turn,
// returning the middle one.
func Between[A, B, C any](a Parser[A], b Parser[B], c Parser[C]) Parser[B] {
	return func(input []byte) (ParseResult[B], error) {
		aResult, err := a(input)
		if err != nil {
			return ParseResult[B]{}, err
		}
		result, err := b(aResult.remainder)
		if err != nil {
			return ParseResult[B]{}, err
		}
		cResult, err := c(result.remainder)
		if err != nil {
			return ParseResult[B]{}, err
		}
		return ParseResult[B]{
			result:    result.result,
			remainder: cResult.remainder,
		}, nil
	}
}

// SepBy returns a parser that runs the first provided parser as many times as
// it can, as long as each successful invocation is separated by a successful
// invocation of the second parser.
func SepBy[A, B any](a Parser[A], b Parser[B]) Parser[[]A] {
	return func(input []byte) (ParseResult[[]A], error) {
		var results []A
		for {
			aResult, err := a(input)
			if err != nil {
				return ParseResult[[]A]{}, err
			}
			input = aResult.remainder

			bResult, err := b(input)
			if err != nil {
				break
			}
			input = bResult.remainder
		}
		return ParseResult[[]A]{
			result:    results,
			remainder: input,
		}, nil
	}
}

// Pair is a pair of values.
type Pair[A, B any] struct {
	First  A
	Second B
}

// Seq returns a parser that runs the two provided parsers in sequence and
// returns both results.
func Seq[A, B any](a Parser[A], b Parser[B]) Parser[Pair[A, B]] {
	return func(input []byte) (ParseResult[Pair[A, B]], error) {
		aResult, err := a(input)
		if err != nil {
			return ParseResult[Pair[A, B]]{}, err
		}
		bResult, err := b(aResult.remainder)
		if err != nil {
			return ParseResult[Pair[A, B]]{}, err
		}

		return ParseResult[Pair[A, B]]{
			result: Pair[A, B]{
				First:  aResult.result,
				Second: bResult.result,
			},
			remainder: bResult.remainder,
		}, nil
	}
}

// SeqL returns a parser that runs a then b, and yields the result of parser
// a.
func SeqL[A, B any](a Parser[A], b Parser[B]) Parser[A] {
	return Apply(Seq(a, b), func(p Pair[A, B]) A { return p.First })
}

// SeqR returns a parser that runs a then b, and yields the result of parser
// b.
func SeqR[A, B any](a Parser[A], b Parser[B]) Parser[B] {
	return Apply(Seq(a, b), func(p Pair[A, B]) B { return p.Second })
}

// Literal returns a parser that parses an exact string.
func Literal(s string) Parser[string] {
	return func(input []byte) (ParseResult[string], error) {
		if !bytes.HasPrefix(input, []byte(s)) {
			return unexpectedError(s, input)
		}
		return ParseResult[string]{
			result:    s,
			remainder: input[len(input):],
		}, nil
	}
}

// Byte returns a parser that parses a single exact byte.
func Byte(b byte) Parser[byte] {
	return func(input []byte) (ParseResult[byte], error) {
		if len(input) == 0 || input[0] != b {
			return unexpectedError(b, input)
		}
		return ParseResult[byte]{
			result:    b,
			remainder: input[1:],
		}, nil
	}
}

// Uint is a parser that parses a single unsigned integer.
func Uint[U constraints.Unsigned](input []byte) (ParseResult[U], error) {
	end := 0
	for {
		r, w := utf8.DecodeRune(input[end:])
		if !unicode.IsDigit(r) {
			break
		}
		end += w
	}
	// TODO: not this temporary string it's pretty silly
	result, err := strconv.ParseUint(string(input[:end]), 10, 64)
	if err != nil {
		return ParseResult[U]{}, err
	}
	return ParseResult[U]{
		result:    U(result),
		remainder: input[end:],
	}, nil
}

func unexpectedError[A any](want A, input []byte) (ParseResult[A], error) {
	// TODO: %q if A is stringy enough?
	if len(input) > 25 {
		input = input[:25]
	}
	return ParseResult[A]{}, fmt.Errorf("expected %v, found: %q", want, input)
}
