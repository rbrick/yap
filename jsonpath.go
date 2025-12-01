package yap

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var IndexedPattern = regexp.MustCompile(`\[([0-9]+)\]`)

type Resolver func(data any) (any, error)

func ArrayIndexResolver(index int) Resolver {
	return func(data any) (any, error) {
		if data == nil {
			return nil, fmt.Errorf("null data")
		}

		switch v := data.(type) {
		case []any:
			if index < 0 || index >= len(v) {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}
			return v[index], nil
		default:
			return nil, fmt.Errorf("data is not an array")
		}
	}
}

func KeyResolver(key string) Resolver {
	return func(data any) (any, error) {
		if data == nil {
			return nil, fmt.Errorf("null data")
		}

		switch v := data.(type) {
		case map[string]any:
			val, exists := v[key]
			if !exists {
				return nil, fmt.Errorf("key %s does not exist", key)
			}
			return val, nil
		default:
			return nil, fmt.Errorf("data is not an object")
		}
	}
}

func RootResolver() Resolver {
	return func(data any) (any, error) {
		return data, nil
	}
}

// Path consists of segments to follow in a JSON structure
// A path semantically looks like: $.store.book[0].title
// with $ = root, store = segment, book = segment, [0] = index, title = segment
// multiple indices are allowed per segment, e.g. book[0][1]
// and segments can be nested, e.g. store.book.title
type Path struct {
	Segments []*Segment
}

func (p *Path) Resolve(data any) (any, error) {
	var current any = data
	var err error

	for _, segment := range p.Segments {
		for _, resolver := range segment.Resolvers {
			current, err = resolver(current)
			if err != nil {
				return nil, err
			}
		}
	}

	return current, nil
}

type Segment struct {
	Name      string
	Resolvers []Resolver
}

func ParseSegment(str string) (*Segment, error) {
	segment := &Segment{
		Name:      str,
		Resolvers: []Resolver{},
	}

	key := IndexedPattern.ReplaceAllString(str, "")

	switch key {
	case "$":
		segment.Resolvers = append(segment.Resolvers, RootResolver())
	default:
		segment.Resolvers = append(segment.Resolvers, KeyResolver(key))
	}

	indices := IndexedPattern.FindAllStringSubmatch(str, -1)

	for _, match := range indices {
		// match[0] is the full match, match[1] is the first capturing group
		index, _ := strconv.Atoi(match[1])

		// array index resolver
		resolver := ArrayIndexResolver(index)

		segment.Resolvers = append(segment.Resolvers, resolver)
	}

	return segment, nil
}

func ParsePath(str string) (*Path, error) {
	segments := strings.Split(str, ".")

	path := &Path{
		Segments: []*Segment{},
	}

	for _, seg := range segments {
		segment := strings.TrimSpace(seg)

		parsedSegment, err := ParseSegment(segment)
		if err != nil {
			return nil, err
		}

		path.Segments = append(path.Segments, parsedSegment)
	}

	return path, nil
}

func NewPath(segments []*Segment) *Path {
	return &Path{
		Segments: segments,
	}
}
