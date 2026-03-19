package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Encoder provides options for rendering JSON.
type Encoder struct {
	Indent      string // e.g., "  " (two spaces). Empty => compact mode.
	SortMapKeys bool   // sort map keys for deterministic output
}

// fieldInfo caches metadata for struct fields.
type fieldInfo struct {
	Name      string // JSON field name
	Index     int    // struct field index
	Omitempty bool   // omit when zero
}

var structCache sync.Map // map[reflect.Type][]fieldInfo

// Marshal renders any Go value as JSON-like text (without encoding/json).
func (e *Encoder) Marshal(v any) string {
	var sb strings.Builder
	e.encode(reflect.ValueOf(v), 0, &sb)
	return sb.String()
}

// --------------------------- Encoding Core ---------------------------

func (e *Encoder) encode(rv reflect.Value, depth int, out *strings.Builder) {
	if !rv.IsValid() {
		out.WriteString("null")
		return
	}

	// Follow pointers
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			out.WriteString("null")
			return
		}
		e.encode(rv.Elem(), depth, out)
		return
	}

	switch rv.Kind() {
	case reflect.Bool:
		out.WriteString(strconv.FormatBool(rv.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out.WriteString(strconv.FormatInt(rv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		out.WriteString(strconv.FormatUint(rv.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		out.WriteString(strconv.FormatFloat(rv.Float(), 'f', -1, 64))
	case reflect.String:
		out.WriteString(strconv.Quote(rv.String()))
	case reflect.Slice, reflect.Array:
		e.encodeSlice(rv, depth, out)
	case reflect.Map:
		e.encodeMap(rv, depth, out)
	case reflect.Struct:
		e.encodeStruct(rv, depth, out)
	case reflect.Interface:
		if rv.IsNil() {
			out.WriteString("null")
		} else {
			e.encode(rv.Elem(), depth, out)
		}
	default:
		// Unsupported kinds (chan, func, complex, unsafe.Pointer) => null
		out.WriteString("null")
	}
}

func (e *Encoder) encodeSlice(rv reflect.Value, depth int, out *strings.Builder) {
	n := rv.Len()
	if n == 0 {
		out.WriteString("[]")
		return
	}
	if e.Indent == "" {
		out.WriteByte('[')
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(',')
			}
			e.encode(rv.Index(i), depth, out)
		}
		out.WriteByte(']')
		return
	}
	out.WriteString("[\n")
	for i := 0; i < n; i++ {
		e.indent(depth+1, out)
		e.encode(rv.Index(i), depth+1, out)
		if i < n-1 {
			out.WriteByte(',')
		}
		out.WriteByte('\n')
	}
	e.indent(depth, out)
	out.WriteByte(']')
}

func (e *Encoder) encodeMap(rv reflect.Value, depth int, out *strings.Builder) {
	if rv.Len() == 0 {
		out.WriteString("{}")
		return
	}

	// Only string keys are valid JSON; others will be stringified.
	keys := rv.MapKeys()
	if e.SortMapKeys {
		sort.Slice(keys, func(i, j int) bool { return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface()) })
	}

	if e.Indent == "" {
		out.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				out.WriteByte(',')
			}
			out.WriteString(strconv.Quote(fmt.Sprint(k.Interface())))
			out.WriteByte(':')
			e.encode(rv.MapIndex(k), depth, out)
		}
		out.WriteByte('}')
		return
	}

	out.WriteString("{\n")
	for i, k := range keys {
		e.indent(depth+1, out)
		out.WriteString(strconv.Quote(fmt.Sprint(k.Interface())))
		out.WriteString(": ")
		e.encode(rv.MapIndex(k), depth+1, out)
		if i < rv.Len()-1 {
			out.WriteByte(',')
		}
		out.WriteByte('\n')
	}
	e.indent(depth, out)
	out.WriteByte('}')
}

func (e *Encoder) encodeStruct(rv reflect.Value, depth int, out *strings.Builder) {
	fields := cachedFields(rv.Type())
	// Build visible fields (respecting omitempty)
	// parts := make([]string, 0, len(fields))

	if e.Indent == "" {
		out.WriteByte('{')
		first := true
		for _, f := range fields {
			fv := rv.Field(f.Index)
			if f.Omitempty && isZero(fv) {
				continue
			}
			if !first {
				out.WriteByte(',')
			}
			first = false
			out.WriteString(strconv.Quote(f.Name))
			out.WriteByte(':')
			e.encode(fv, depth, out)
		}
		if first { // no fields emitted
			out.WriteByte('}')
			return
		}
		out.WriteByte('}')
		return
	}

	// Pretty (indented) path
	out.WriteString("{\n")
	emitted := 0
	for i, f := range fields {
		fv := rv.Field(f.Index)
		if f.Omitempty && isZero(fv) {
			continue
		}
		if emitted > 0 {
			out.WriteByte('\n') // ensure each new field goes on its own line
		}
		e.indent(depth+1, out)
		out.WriteString(strconv.Quote(f.Name))
		out.WriteString(": ")
		e.encode(fv, depth+1, out)

		// Find if there are more fields to emit after this one
		if hasMoreNonEmpty(fields[i+1:], rv) {
			out.WriteByte(',')
		}
		emitted++
	}
	if emitted > 0 {
		out.WriteByte('\n')
	}
	e.indent(depth, out)
	out.WriteByte('}')
}

func (e *Encoder) indent(depth int, out *strings.Builder) {
	if e.Indent == "" {
		return
	}
	for i := 0; i < depth; i++ {
		out.WriteString(e.Indent)
	}
}

// --------------------------- Struct Field Caching ---------------------------
func cachedFields(t reflect.Type) []fieldInfo {
	if v, ok := structCache.Load(t); ok {
		return v.([]fieldInfo)
	}
	// Build field list (only exported, direct fields; embedded support omitted for brevity)
	var fields []fieldInfo
	n := t.NumField()
	for i := 0; i < n; i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" { // unexported
			continue
		}
		name := sf.Name
		omitempty := false
		if tag, ok := sf.Tag.Lookup("json"); ok {
			if tag == "-" {
				continue
			}
			// Support `json:"name,omitempty"` and `json:"name"`
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
			for _, opt := range parts[1:] {
				if opt == "omitempty" {
					omitempty = true
				}
			}
		}
		fields = append(fields, fieldInfo{
			Name:      name,
			Index:     i,
			Omitempty: omitempty,
		})
	}
	structCache.Store(t, fields)
	return fields
}

func hasMoreNonEmpty(rest []fieldInfo, rv reflect.Value) bool {
	for _, f := range rest {
		fv := rv.Field(f.Index)
		if f.Omitempty && isZero(fv) {
			continue
		}
		return true
	}
	return false
}

func isZero(v reflect.Value) bool {
	// Use reflect.Value.IsZero when available; here we implement a small fallback.
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		// Consider struct zero if all fields are zero (shallow check)
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return !v.IsValid()
	}
}

// --------------------------- Demo Types & main ---------------------------

type TLSConfig struct {
	Enabled bool   `json:"enabled,omitempty"`
	CertPEM string `json:"cert_pem,omitempty"`
	KeyPEM  string `json:"key_pem,omitempty"`
}

type ServerConfig struct {
	Name     string            `json:"name"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags,omitempty"`
	Limits   map[string]uint32 `json:"limits,omitempty"`
	TLS      *TLSConfig        `json:"tls,omitempty"`
	Metadata map[string]any    `json:"metadata,omitempty"`
}

func main() {
	cfg := ServerConfig{
		Name: "api",
		Port: 8080,
		Tags: []string{"go", "prod"},
		Limits: map[string]uint32{
			"requests_per_min": 1200,
			"conn":             100,
		},
		TLS: &TLSConfig{
			Enabled: true,
			CertPEM: "-----BEGIN CERTIFICATE-----...",
			// KeyPEM omitted => omitempty
		},
		Metadata: map[string]any{
			"region":   "eu-central-1",
			"revision": 42,
			"debug":    false,
		},
	}

	pretty := (&Encoder{Indent: "  ", SortMapKeys: true}).Marshal(cfg)
	compact := (&Encoder{}).Marshal(cfg)

	fmt.Println("=== Pretty ===")
	fmt.Println(pretty)
	fmt.Println("\n=== Compact ===")
	fmt.Println(compact)
}
