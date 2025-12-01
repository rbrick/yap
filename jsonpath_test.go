package yap

import (
	"encoding/json"
	"testing"
)

func TestParsePath(t *testing.T) {
	pathStr := "$.store.book[0].title"
	path, err := ParsePath(pathStr)

	if err != nil {
		t.Fatalf("failed to parse path: %v", err)
	}

	if len(path.Segments) != 4 {
		t.Fatalf("expected 4 segments, got %d", len(path.Segments))
	}

	expectedSegments := []string{"$", "store", "book[0]", "title"}

	for i, seg := range path.Segments {
		if seg.Name != expectedSegments[i] {
			t.Errorf("expected segment %d to be %s, got %s", i, expectedSegments[i], seg.Name)
		}
	}
}

func TestSegmentResolvers(t *testing.T) {
	segmentStr := "book[0]"
	segment, err := ParseSegment(segmentStr)

	if err != nil {
		t.Fatalf("failed to parse segment: %v", err)
	}

	if segment.Name != "book[0]" {
		t.Errorf("expected segment name to be 'book[0]', got '%s'", segment.Name)
	}

	if len(segment.Resolvers) != 2 {
		t.Fatalf("expected 2 resolvers, got %d", len(segment.Resolvers))
	}
}

func TestRootSegment(t *testing.T) {
	segmentStr := "$"
	segment, err := ParseSegment(segmentStr)

	if err != nil {
		t.Fatalf("failed to parse root segment: %v", err)
	}

	if segment.Name != "$" {
		t.Errorf("expected segment name to be '$', got '%s'", segment.Name)
	}

	if len(segment.Resolvers) != 1 {
		t.Fatalf("expected 1 resolver for root, got %d", len(segment.Resolvers))
	}
}

func TestResolution(t *testing.T) {

	jsonData := `{
		"store": {
			"book": [
				{
					"title": "The Great Gatsby"
				},
				{
					"title": "1984"
				}
			]
		}
	}`

	var data any
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	pathStr := "$.store.book[0].title"
	path, err := ParsePath(pathStr)
	if err != nil {
		t.Fatalf("failed to parse path: %v", err)
	}

	resolved, err := path.Resolve(data)
	if err != nil {
		t.Fatalf("failed to resolve path: %v", err)
	}

	expectedTitle := "The Great Gatsby"
	if resolved != expectedTitle {
		t.Errorf("expected title to be '%s', got '%v'", expectedTitle, resolved)
	}
}
