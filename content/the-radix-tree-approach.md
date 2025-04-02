---
{
    "title": "Building a Fast Router in Go: The Radix Tree Approach",
    "date": "2025-04-02T02:58:57Z"
}
---

Ever wondered why some web frameworks are lightning fast while others feel like they're running through molasses? Let me tell you about one of the coolest parts of building Zinc - our router's radix tree implementation.

## The Problem

When you're building a web framework, routing is everything. You need to match URLs to handlers quickly. While Go 1.22's `net/http` package brought some nice improvements to pattern matching, it still doesn't handle complex routing patterns with the flexibility and performance we need. It's like trying to find a needle in a haystack - but the haystack is your URL patterns and the needle is the right handler.

## Enter the Radix Tree

A radix tree (or compressed trie) is perfect for this. Think of it like a family tree, but for URL patterns. Each node represents a part of a URL, and we compress common prefixes to save space. It's like having a smart address book where "Matt" and "Matthew" share the same "Mat" prefix.

Here's a quick example of how our routes look in the tree:

```
/           -> HomeHandler
├── about   -> AboutHandler
├── post/   -> PostHandler
│   └── :id -> PostDetailHandler
└── api/    -> APIHandler
    └── v1/ -> APIv1Handler
```

## Why It's Fast

The beauty of this approach is O(k) lookup time, where k is the length of the URL. No more scanning through a list of routes - we just follow the tree branches. It's like having a GPS for your URLs.

## The Cool Part

What makes our implementation special is how we handle dynamic segments (like `:id` in `/post/:id`). Instead of doing regex matches or string splitting, we use the tree structure itself to identify parameters. It's elegant, and it's fast.

## The Code

Here's a sneak peek at how it works (simplified, of course):

```go
// Node represents a single node in our radix tree
type Node struct {
    path     string
    handler  http.HandlerFunc
    children map[string]*Node
    params   []string
}

// addRoute adds a new route to the tree
func (n *Node) addRoute(path string, handler http.HandlerFunc) {
    // Split path into segments
    segments := strings.Split(path, "/")
    current := n
    
    for _, segment := range segments {
        if segment == "" {
            continue
        }
        
        // Check if it's a parameter
        if strings.HasPrefix(segment, ":") {
            current.params = append(current.params, segment[1:])
        }
        
        // Add to tree
        if _, exists := current.children[segment]; !exists {
            current.children[segment] = &Node{
                path:     segment,
                children: make(map[string]*Node),
            }
        }
        current = current.children[segment]
    }
    
    current.handler = handler
}
```

## The Payoff

The result? Lightning-fast route matching, even with thousands of routes. We're talking microseconds here, not milliseconds. And the memory footprint? Tiny compared to regex-based solutions.

## Wrapping Up

Building a fast router isn't just about matching URLs - it's about doing it efficiently at scale. The radix tree approach gives us that sweet spot of performance and simplicity. Plus, it's just cool to see how a data structure from the 1960s is still kicking ass in modern web frameworks.

Next time you're building a web framework (or just curious about how they work), give the radix tree a look. It might just be the performance boost you're looking for.

