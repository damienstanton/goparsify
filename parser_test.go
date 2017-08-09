package goparsify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsify(t *testing.T) {

	t.Run("strings", func(t *testing.T) {
		require.Equal(t, "ff", Parsify("ff")(NewState("ffooo")).Token)
	})

	t.Run("parsers", func(t *testing.T) {
		require.Equal(t, "ff", Parsify(Chars("f"))(NewState("ffooo")).Token)
	})

	t.Run("parser funcs", func(t *testing.T) {
		node := Parsify(func(p *State) Result {
			return Result{Token: "hello"}
		})(NewState("ffooo"))

		require.Equal(t, "hello", node.Token)
	})

	t.Run("*parsers", func(t *testing.T) {
		var parser Parser
		parserfied := Parsify(&parser)
		parser = Chars("f")

		node := parserfied(NewState("ffooo"))
		require.Equal(t, "ff", node.Token)
	})

	require.Panics(t, func() {
		Parsify(1)
	})
}

func TestParsifyAll(t *testing.T) {
	parsers := ParsifyAll("ff", "gg")

	result := parsers[0](NewState("ffooo"))
	require.Equal(t, "ff", result.Token)

	result = parsers[1](NewState("ffooo"))
	require.Equal(t, "", result.Token)
}

func TestExact(t *testing.T) {
	t.Run("success string", func(t *testing.T) {
		node, ps := runParser("foobar", Exact("fo"))
		require.Equal(t, "fo", node.Token)
		require.Equal(t, "obar", ps.Get())
	})

	t.Run("success char", func(t *testing.T) {
		node, ps := runParser("foobar", Exact("f"))
		require.Equal(t, "f", node.Token)
		require.Equal(t, "oobar", ps.Get())
	})

	t.Run("error", func(t *testing.T) {
		_, ps := runParser("foobar", Exact("bar"))
		require.Equal(t, "bar", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("error char", func(t *testing.T) {
		_, ps := runParser("foobar", Exact("o"))
		require.Equal(t, "o", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("eof char", func(t *testing.T) {
		_, ps := runParser("", Exact("o"))
		require.Equal(t, "o", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})
}

func TestChars(t *testing.T) {
	t.Run("full match", func(t *testing.T) {
		node, ps := runParser("foobar", Chars("a-z"))
		require.Equal(t, "foobar", node.Token)
		require.Equal(t, "", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("partial match", func(t *testing.T) {
		node, ps := runParser("a1b2c3d4efg", Chars("1-4d-a"))
		require.Equal(t, "a1b2c3d4", node.Token)
		require.Equal(t, "efg", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("limited match", func(t *testing.T) {
		node, ps := runParser("a1b2c3d4efg", Chars("1-4d-a", 1, 2))
		require.Equal(t, "a1", node.Token)
		require.Equal(t, "b2c3d4efg", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("no match", func(t *testing.T) {
		_, ps := runParser("ffffff", Chars("0-9"))
		require.Equal(t, "offset 0: expected 0-9", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("no match with min", func(t *testing.T) {
		_, ps := runParser("ffffff", Chars("0-9", 4))
		require.Equal(t, "0-9", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("test exact matches", func(t *testing.T) {
		node, ps := runParser("aaff", Chars("abcd"))
		require.Equal(t, "aa", node.Token)
		require.Equal(t, 2, ps.Pos)
		require.False(t, ps.Errored())
	})

	t.Run("test not matches", func(t *testing.T) {
		node, ps := runParser("aaff", NotChars("ff"))
		require.Equal(t, "aa", node.Token)
		require.Equal(t, 2, ps.Pos)
		require.False(t, ps.Errored())
	})

	require.Panics(t, func() {
		Chars("a-b", 1, 2, 3)
	})
}

func TestParseString(t *testing.T) {
	Y := Map("hello", func(n Result) Result { return Result{Result: n.Token} })

	t.Run("full match", func(t *testing.T) {
		result, err := Run(Y, "hello")
		require.Equal(t, "hello", result)
		require.NoError(t, err)
	})

	t.Run("partial match", func(t *testing.T) {
		result, err := Run(Y, "hello world")
		require.Equal(t, "hello", result)
		require.Error(t, err)
		require.Equal(t, "left unparsed: world", err.Error())
	})

	t.Run("error", func(t *testing.T) {
		result, err := Run(Y, "world")
		require.Nil(t, result)
		require.Error(t, err)
		require.Equal(t, "offset 0: expected hello", err.Error())
	})
}

func runParser(input string, parser Parser) (Result, *State) {
	ps := NewState(input)
	result := parser(ps)
	return result, ps
}
