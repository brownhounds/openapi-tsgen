package schema

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrNilDoc                         = errors.New("nil doc")
	ErrUnsupportedRef                 = errors.New("unsupported $ref")
	ErrMissingComponentPathItem       = errors.New("missing components.pathItems")
	ErrNestedPathItemRef              = errors.New("nested pathItem $ref")
	ErrMissingComponentRequestBody    = errors.New("missing components.requestBodies")
	ErrNestedRequestBodyRef           = errors.New("nested requestBody $ref")
	ErrMissingComponentResponse       = errors.New("missing components.responses")
	ErrNestedResponseRef              = errors.New("nested response $ref")
	ErrMissingComponentParameter      = errors.New("missing components.parameters")
	ErrNestedParameterRef             = errors.New("nested parameter $ref")
	ErrMissingComponentHeader         = errors.New("missing components.headers")
	ErrNestedHeaderRef                = errors.New("nested header $ref")
	ErrMissingComponentSecurityScheme = errors.New("missing components.securitySchemes")
	ErrNestedSecuritySchemeRef        = errors.New("nested securityScheme $ref")
)

type IR struct {
	Paths                     map[string]IRPathItem
	Webhooks                  map[string]IRPathItem
	ComponentsSchemas         map[string]string
	ComponentsResponses       map[string]string
	ComponentsRequestBody     map[string]string
	ComponentsParameters      map[string]string
	ComponentsHeaders         map[string]string
	ComponentsSecuritySchemes map[string]string
	Enums                     map[string]string
	Servers                   []Server
}

type schemaMode int

type IRPathItem struct {
	Ops map[string]IROperation
}

type IROperation struct {
	PathParams   map[string]paramResolved
	QueryParams  map[string]paramResolved
	HeaderParams map[string]paramResolved
	CookieParams map[string]paramResolved
	Responses    map[string]string
	RequestBody  string
	Security     []SecurityRequirement
	Servers      []Server
}

func ToIR(doc *Document) (*IR, error) {
	if doc == nil {
		return nil, ErrNilDoc
	}

	out := &IR{
		Paths:                     map[string]IRPathItem{},
		Webhooks:                  map[string]IRPathItem{},
		ComponentsSchemas:         map[string]string{},
		ComponentsResponses:       map[string]string{},
		ComponentsRequestBody:     map[string]string{},
		ComponentsParameters:      map[string]string{},
		ComponentsHeaders:         map[string]string{},
		ComponentsSecuritySchemes: map[string]string{},
		Enums:                     map[string]string{},
		Servers:                   doc.Servers,
	}

	ctx := newEnumContext(out.Enums)
	if err := populateComponents(out, doc, ctx); err != nil {
		return nil, err
	}
	if err := populatePaths(out, doc, ctx); err != nil {
		return nil, err
	}
	if err := populateWebhooks(out, doc, ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func populateComponents(out *IR, doc *Document, ctx *enumContext) error {
	if doc.Components == nil {
		return nil
	}
	if err := populateComponentSchemas(out, doc, ctx); err != nil {
		return err
	}
	if err := populateComponentResponses(out, doc, ctx); err != nil {
		return err
	}
	if err := populateComponentRequestBodies(out, doc, ctx); err != nil {
		return err
	}
	if err := populateComponentParameters(out, doc, ctx); err != nil {
		return err
	}
	if err := populateComponentHeaders(out, doc, ctx); err != nil {
		return err
	}
	if err := populateComponentSecuritySchemes(out, doc); err != nil {
		return err
	}
	return nil
}

func populateComponentSchemas(out *IR, doc *Document, ctx *enumContext) error {
	if len(doc.Components.Schemas) == 0 {
		return nil
	}
	keys := make([]string, 0, len(doc.Components.Schemas))
	for k := range doc.Components.Schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sch := doc.Components.Schemas[k]
		out.ComponentsSchemas[k] = schemaToTS(doc, &RefOr[Schema]{Value: &sch}, 0, ctx, k, modeDefault)
	}
	return nil
}

func populateComponentResponses(out *IR, doc *Document, ctx *enumContext) error {
	if len(doc.Components.Responses) == 0 {
		return nil
	}
	keys := make([]string, 0, len(doc.Components.Responses))
	for k := range doc.Components.Responses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		respRef := doc.Components.Responses[k]
		resp, err := resolveResponse(doc, respRef)
		if err != nil {
			return fmt.Errorf("components.responses.%s: %w", k, err)
		}
		out.ComponentsResponses[k] = responseToTS(doc, resp, ctx, k, modeOutput)
	}
	return nil
}

func populateComponentRequestBodies(out *IR, doc *Document, ctx *enumContext) error {
	if len(doc.Components.RequestBodies) == 0 {
		return nil
	}
	keys := make([]string, 0, len(doc.Components.RequestBodies))
	for k := range doc.Components.RequestBodies {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		rb, err := resolveRequestBody(doc, doc.Components.RequestBodies[k])
		if err != nil {
			return fmt.Errorf("components.requestBodies.%s: %w", k, err)
		}
		out.ComponentsRequestBody[k] = requestBodyToTS(doc, rb, ctx, k, modeInput)
	}
	return nil
}

func populateComponentParameters(out *IR, doc *Document, ctx *enumContext) error {
	if len(doc.Components.Parameters) == 0 {
		return nil
	}
	keys := make([]string, 0, len(doc.Components.Parameters))
	for k := range doc.Components.Parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		p, err := resolveParameter(doc, doc.Components.Parameters[k])
		if err != nil {
			return fmt.Errorf("components.parameters.%s: %w", k, err)
		}
		out.ComponentsParameters[k] = parameterToTS(doc, p, ctx, k, modeInput)
	}
	return nil
}

func populateComponentHeaders(out *IR, doc *Document, ctx *enumContext) error {
	if len(doc.Components.Headers) == 0 {
		return nil
	}
	keys := make([]string, 0, len(doc.Components.Headers))
	for k := range doc.Components.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h, err := resolveHeader(doc, doc.Components.Headers[k])
		if err != nil {
			return fmt.Errorf("components.headers.%s: %w", k, err)
		}
		out.ComponentsHeaders[k] = headerToTS(doc, h, ctx, k, modeOutput)
	}
	return nil
}

func populateComponentSecuritySchemes(out *IR, doc *Document) error {
	if len(doc.Components.SecuritySchemes) == 0 {
		return nil
	}
	keys := make([]string, 0, len(doc.Components.SecuritySchemes))
	for k := range doc.Components.SecuritySchemes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ss, err := resolveSecurityScheme(doc, doc.Components.SecuritySchemes[k])
		if err != nil {
			return fmt.Errorf("components.securitySchemes.%s: %w", k, err)
		}
		out.ComponentsSecuritySchemes[k] = securitySchemeToTS(ss)
	}
	return nil
}

func populatePaths(out *IR, doc *Document, ctx *enumContext) error {
	pathKeys := make([]string, 0, len(doc.Paths))
	for k := range doc.Paths {
		pathKeys = append(pathKeys, k)
	}
	sort.Strings(pathKeys)

	for _, path := range pathKeys {
		pi, err := resolvePathItem(doc, doc.Paths[path])
		if err != nil {
			return fmt.Errorf("path %q: %w", path, err)
		}
		if pi == nil {
			continue
		}

		ops, err := pathItemToOps(doc, pi, ctx)
		if err != nil {
			return fmt.Errorf("path %q: %w", path, err)
		}

		if len(ops) == 0 {
			continue
		}

		out.Paths[path] = IRPathItem{Ops: ops}
	}

	return nil
}

func populateWebhooks(out *IR, doc *Document, ctx *enumContext) error {
	if len(doc.Webhooks) == 0 {
		return nil
	}

	keys := make([]string, 0, len(doc.Webhooks))
	for k := range doc.Webhooks {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		pi, err := resolvePathItem(doc, doc.Webhooks[name])
		if err != nil {
			return fmt.Errorf("webhook %q: %w", name, err)
		}
		if pi == nil {
			continue
		}

		ops, err := pathItemToOps(doc, pi, ctx)
		if err != nil {
			return fmt.Errorf("webhook %q: %w", name, err)
		}
		if len(ops) == 0 {
			continue
		}
		out.Webhooks[name] = IRPathItem{Ops: ops}
	}
	return nil
}

func pathItemToOps(doc *Document, pi *PathItem, ctx *enumContext) (map[string]IROperation, error) {
	pathParams, err := collectParams(doc, pi.Parameters, ctx)
	if err != nil {
		return nil, err
	}

	ops := map[string]IROperation{}
	methods := []struct {
		op   *Operation
		name string
	}{
		{op: pi.Get, name: "get"},
		{op: pi.Post, name: "post"},
		{op: pi.Put, name: "put"},
		{op: pi.Patch, name: "patch"},
		{op: pi.Delete, name: "delete"},
		{op: pi.Options, name: "options"},
		{op: pi.Head, name: "head"},
		{op: pi.Trace, name: "trace"},
	}

	addOp := func(method string, op *Operation) error {
		if op == nil {
			return nil
		}

		opParams, err := collectParams(doc, op.Parameters, ctx)
		if err != nil {
			return fmt.Errorf("%s params: %w", method, err)
		}

		mergedParams := mergeParamMaps(pathParams, opParams)
		pathOnly := filterParamsByIn(mergedParams, "path")
		queryOnly := filterParamsByIn(mergedParams, "query")
		headerOnly := filterParamsByIn(mergedParams, "header")
		cookieOnly := filterParamsByIn(mergedParams, "cookie")

		reqTS, err := opRequestBodyTS(doc, op, ctx, method)
		if err != nil {
			return fmt.Errorf("%s requestBody: %w", method, err)
		}

		respTS, err := opResponsesTS(doc, op, ctx, method)
		if err != nil {
			return fmt.Errorf("%s responses: %w", method, err)
		}

		security := op.Security
		if len(security) == 0 {
			security = doc.Security
		}

		servers := op.Servers
		if len(servers) == 0 {
			servers = pi.Servers
		}
		if len(servers) == 0 {
			servers = doc.Servers
		}

		ops[method] = IROperation{
			PathParams:   pathOnly,
			QueryParams:  queryOnly,
			HeaderParams: headerOnly,
			CookieParams: cookieOnly,
			RequestBody:  reqTS,
			Responses:    respTS,
			Security:     security,
			Servers:      servers,
		}
		return nil
	}

	for _, method := range methods {
		if err := addOp(method.name, method.op); err != nil {
			return nil, err
		}
	}

	return ops, nil
}

type paramKey struct {
	Name string
	In   string
}

type paramResolved struct {
	In       string
	TS       string
	Required bool
}

func collectParams(doc *Document, params []RefOr[Parameter], ctx *enumContext) (map[paramKey]paramResolved, error) {
	out := map[paramKey]paramResolved{}
	for i := range params {
		refTS := ""
		if params[i].Ref != "" {
			if name, ok := refComponentName(params[i].Ref, "parameters"); ok {
				refTS = componentParameterRef(name)
			}
		}
		p, err := resolveParameter(doc, params[i])
		if err != nil {
			return nil, err
		}
		if p == nil {
			continue
		}
		ts := refTS
		if ts == "" {
			ts = parameterToTS(doc, p, ctx, p.Name, modeInput)
		}
		required := p.Required
		if p.In == "path" {
			required = true
		}
		out[paramKey{Name: p.Name, In: p.In}] = paramResolved{In: p.In, TS: ts, Required: required}
	}
	return out, nil
}

func mergeParamMaps(a, b map[paramKey]paramResolved) map[paramKey]paramResolved {
	out := map[paramKey]paramResolved{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

func filterParamsByIn(m map[paramKey]paramResolved, in string) map[string]paramResolved {
	out := map[string]paramResolved{}
	for k, v := range m {
		if k.In == in {
			out[k.Name] = v
		}
	}
	return out
}

func opRequestBodyTS(doc *Document, op *Operation, ctx *enumContext, opName string) (string, error) {
	if op.RequestBody == nil {
		return tsNever, nil
	}
	if op.RequestBody.Ref != "" {
		if name, ok := refComponentName(op.RequestBody.Ref, "requestBodies"); ok {
			return componentRequestBodyRef(name), nil
		}
	}
	rb, err := resolveRequestBody(doc, *op.RequestBody)
	if err != nil {
		return "", err
	}
	return requestBodyToTS(doc, rb, ctx, joinEnumHint(opName, "RequestBody"), modeInput), nil
}

func opResponsesTS(doc *Document, op *Operation, ctx *enumContext, opName string) (map[string]string, error) {
	out := map[string]string{}

	if len(op.Responses) == 0 {
		return out, nil
	}

	keys := make([]string, 0, len(op.Responses))
	for k := range op.Responses {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		ai, aok := parseStatus(a)
		bi, bok := parseStatus(b)
		if aok && bok {
			return ai < bi
		}
		if aok != bok {
			return aok
		}
		return a < b
	})

	for _, code := range keys {
		r := op.Responses[code]
		if r.Ref != "" {
			if name, ok := refComponentName(r.Ref, "responses"); ok {
				out[code] = componentResponseRef(name)
				continue
			}
		}
		resp, err := resolveResponse(doc, r)
		if err != nil {
			return nil, err
		}
		out[code] = responseToTS(doc, resp, ctx, joinEnumHint(opName, "Response_"+code), modeOutput)
	}

	return out, nil
}

func parseStatus(s string) (int, bool) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return i, true
}

func contentToTS(doc *Document, content map[string]MediaType, empty string, ctx *enumContext, nameHint string, mode schemaMode) string {
	if len(content) == 0 {
		return empty
	}
	keys := make([]string, 0, len(content))
	for k := range content {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	seen := map[string]bool{}
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		mt := content[k]
		ts := tsUnknown
		if mt.Schema != nil {
			ts = schemaToTS(doc, mt.Schema, 0, ctx, joinEnumHint(nameHint, mediaTypeSuffix(k)), mode)
		}
		if !seen[ts] {
			seen[ts] = true
			parts = append(parts, ts)
		}
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return "(" + strings.Join(parts, " | ") + ")"
}

func requestBodyToTS(doc *Document, rb *RequestBody, ctx *enumContext, nameHint string, mode schemaMode) string {
	if rb == nil {
		return tsUnknown
	}
	return contentToTS(doc, rb.Content, tsUnknown, ctx, nameHint, mode)
}

func responseToTS(doc *Document, resp *Response, ctx *enumContext, nameHint string, mode schemaMode) string {
	if resp == nil {
		return tsNever
	}

	bodyTS := contentToTS(doc, resp.Content, tsNever, ctx, nameHint, mode)

	if len(resp.Headers) == 0 {
		return bodyTS
	}

	keys := make([]string, 0, len(resp.Headers))
	for k := range resp.Headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString("  headers: {\n")
	for _, k := range keys {
		h := resp.Headers[k]
		refTS := ""
		if h.Ref != "" {
			if name, ok := refComponentName(h.Ref, "headers"); ok {
				refTS = componentHeaderRef(name)
			}
		}
		hv, err := resolveHeader(doc, h)
		if err != nil {
			refTS = tsUnknown
			hv = nil
		}
		fieldTS := refTS
		if fieldTS == "" {
			fieldTS = headerToTS(doc, hv, ctx, k, modeOutput)
		}
		if hv != nil && hv.Required {
			b.WriteString("    " + safeProp(k) + ": " + fieldTS + ";\n")
		} else {
			b.WriteString("    " + safeProp(k) + "?: " + fieldTS + ";\n")
		}
	}
	b.WriteString("  };\n")
	b.WriteString("  body: " + bodyTS + ";\n")
	b.WriteString("}")
	return b.String()
}

func parameterToTS(doc *Document, p *Parameter, ctx *enumContext, nameHint string, mode schemaMode) string {
	if p == nil {
		return tsUnknown
	}
	if p.Schema != nil {
		return schemaToTS(doc, p.Schema, 0, ctx, nameHint, mode)
	}
	return contentToTS(doc, p.Content, tsUnknown, ctx, nameHint, mode)
}

func headerToTS(doc *Document, h *Header, ctx *enumContext, nameHint string, mode schemaMode) string {
	if h == nil {
		return tsUnknown
	}
	if h.Schema != nil {
		return schemaToTS(doc, h.Schema, 0, ctx, nameHint, mode)
	}
	return contentToTS(doc, h.Content, tsUnknown, ctx, nameHint, mode)
}

func resolvePathItem(doc *Document, v RefOr[PathItem]) (*PathItem, error) {
	if v.Ref == "" {
		return v.Value, nil
	}
	name, ok := refComponentName(v.Ref, "pathItems")
	if !ok || doc.Components == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedRef, v.Ref)
	}
	pi, ok := doc.Components.PathItems[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingComponentPathItem, name)
	}
	if pi.Ref != "" {
		return nil, fmt.Errorf("%w: %q", ErrNestedPathItemRef, pi.Ref)
	}
	return pi.Value, nil
}

func resolveRequestBody(doc *Document, v RefOr[RequestBody]) (*RequestBody, error) {
	if v.Ref == "" {
		return v.Value, nil
	}
	name, ok := refComponentName(v.Ref, "requestBodies")
	if !ok || doc.Components == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedRef, v.Ref)
	}
	rb, ok := doc.Components.RequestBodies[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingComponentRequestBody, name)
	}
	if rb.Ref != "" {
		return nil, fmt.Errorf("%w: %q", ErrNestedRequestBodyRef, rb.Ref)
	}
	return rb.Value, nil
}

func resolveResponse(doc *Document, v RefOr[Response]) (*Response, error) {
	if v.Ref == "" {
		return v.Value, nil
	}
	name, ok := refComponentName(v.Ref, "responses")
	if !ok || doc.Components == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedRef, v.Ref)
	}
	r, ok := doc.Components.Responses[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingComponentResponse, name)
	}
	if r.Ref != "" {
		return nil, fmt.Errorf("%w: %q", ErrNestedResponseRef, r.Ref)
	}
	return r.Value, nil
}

func resolveParameter(doc *Document, v RefOr[Parameter]) (*Parameter, error) {
	if v.Ref == "" {
		return v.Value, nil
	}
	name, ok := refComponentName(v.Ref, "parameters")
	if !ok || doc.Components == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedRef, v.Ref)
	}
	p, ok := doc.Components.Parameters[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingComponentParameter, name)
	}
	if p.Ref != "" {
		return nil, fmt.Errorf("%w: %q", ErrNestedParameterRef, p.Ref)
	}
	return p.Value, nil
}

func resolveHeader(doc *Document, v RefOr[Header]) (*Header, error) {
	if v.Ref == "" {
		return v.Value, nil
	}
	name, ok := refComponentName(v.Ref, "headers")
	if !ok || doc.Components == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedRef, v.Ref)
	}
	h, ok := doc.Components.Headers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingComponentHeader, name)
	}
	if h.Ref != "" {
		return nil, fmt.Errorf("%w: %q", ErrNestedHeaderRef, h.Ref)
	}
	return h.Value, nil
}

func resolveSecurityScheme(doc *Document, v RefOr[SecurityScheme]) (*SecurityScheme, error) {
	if v.Ref == "" {
		return v.Value, nil
	}
	name, ok := refComponentName(v.Ref, "securitySchemes")
	if !ok || doc.Components == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedRef, v.Ref)
	}
	ss, ok := doc.Components.SecuritySchemes[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrMissingComponentSecurityScheme, name)
	}
	if ss.Ref != "" {
		return nil, fmt.Errorf("%w: %q", ErrNestedSecuritySchemeRef, ss.Ref)
	}
	return ss.Value, nil
}

func refComponentName(ref, section string) (string, bool) {
	prefix := "#/components/" + section + "/"
	if strings.HasPrefix(ref, prefix) {
		return strings.TrimPrefix(ref, prefix), true
	}
	return "", false
}

func schemaToTS(doc *Document, s *RefOr[Schema], depth int, ctx *enumContext, nameHint string, mode schemaMode) string {
	if s == nil {
		return tsUnknown
	}
	if s.Ref != "" {
		return schemaRefToTS(doc, s.Ref, depth, ctx, mode)
	}
	if s.Value == nil {
		return tsUnknown
	}
	if depth > 30 {
		return tsUnknown
	}

	return schemaValueToTS(doc, s.Value.Other, depth, ctx, nameHint, mode)
}

func schemaRefToTS(doc *Document, ref string, depth int, ctx *enumContext, mode schemaMode) string {
	if name, ok := refComponentName(ref, "schemas"); ok {
		if mode == modeDefault {
			return componentSchemaRef(name)
		}
		if doc != nil && doc.Components != nil {
			if sch, ok := doc.Components.Schemas[name]; ok {
				return schemaToTS(doc, &RefOr[Schema]{Value: &sch}, depth+1, ctx, name, mode)
			}
		}
		return componentSchemaRef(name)
	}
	return tsUnknown
}

func schemaValueToTS(doc *Document, o map[string]any, depth int, ctx *enumContext, nameHint string, mode schemaMode) string {
	if ts, ok := schemaEnumOrConstToTS(o, ctx, nameHint); ok {
		return ts
	}
	if ts, ok := schemaCombinatorToTS(doc, o, depth, ctx, nameHint, mode); ok {
		return ts
	}
	return schemaTypeToTS(doc, o, depth, ctx, nameHint, mode)
}

func schemaEnumOrConstToTS(o map[string]any, ctx *enumContext, nameHint string) (string, bool) {
	if cv, ok := o["const"]; ok {
		if ctx != nil {
			if ts := ctx.emitEnum(nameHint, []any{cv}, o); ts != "" {
				return ts, true
			}
		}
		if ts := literalToTS(cv); ts != "" {
			return applyNullable(ts, o), true
		}
	}
	if ev := anySlice(o["enum"]); len(ev) > 0 {
		if ctx != nil {
			if ts := ctx.emitEnum(nameHint, ev, o); ts != "" {
				return ts, true
			}
		}
		parts := make([]string, 0, len(ev))
		for _, it := range ev {
			if ts := literalToTS(it); ts != "" {
				parts = append(parts, ts)
			} else {
				return applyNullable(tsUnknown, o), true
			}
		}
		return applyNullable("("+strings.Join(parts, " | ")+")", o), true
	}
	return "", false
}

func schemaCombinatorToTS(doc *Document, o map[string]any, depth int, ctx *enumContext, nameHint string, mode schemaMode) (string, bool) {
	if oneOf := anySlice(o["oneOf"]); len(oneOf) > 0 {
		return applyNullable(schemaListToTS(doc, oneOf, depth, ctx, nameHint, "OneOf", " | ", mode), o), true
	}
	if anyOf := anySlice(o["anyOf"]); len(anyOf) > 0 {
		return applyNullable(schemaListToTS(doc, anyOf, depth, ctx, nameHint, "AnyOf", " | ", mode), o), true
	}
	if allOf := anySlice(o["allOf"]); len(allOf) > 0 {
		return applyNullable(schemaListToTS(doc, allOf, depth, ctx, nameHint, "AllOf", " & ", mode), o), true
	}
	return "", false
}

func schemaListToTS(doc *Document, items []any, depth int, ctx *enumContext, nameHint, hintPrefix, joiner string, mode schemaMode) string {
	parts := make([]string, 0, len(items))
	for i, it := range items {
		parts = append(parts, schemaAnyToTS(doc, it, depth+1, ctx, joinEnumHint(nameHint, hintPrefix+strconv.Itoa(i+1)), mode))
	}
	return "(" + strings.Join(parts, joiner) + ")"
}

func schemaTypeToTS(doc *Document, o map[string]any, depth int, ctx *enumContext, nameHint string, mode schemaMode) string {
	t, _ := o["type"].(string)
	switch t {
	case schemaTypeString:
		return applyNullable(schemaTypeString, o)
	case "number", "integer":
		return applyNullable("number", o)
	case "boolean":
		return applyNullable("boolean", o)
	case schemaTypeNull:
		return schemaTypeNull
	case "array":
		items := o["items"]
		if items == nil {
			return applyNullable(tsUnknown+"[]", o)
		}
		return applyNullable(schemaAnyToTS(doc, items, depth+1, ctx, joinEnumHint(nameHint, "Item"), mode)+"[]", o)
	case "object":
		return applyNullable(objectToTS(doc, o, depth+1, ctx, nameHint, mode), o)
	case "":
		if props, ok := o["properties"].(map[string]any); ok && len(props) > 0 {
			return applyNullable(objectToTS(doc, o, depth+1, ctx, nameHint, mode), o)
		}
		if req := anySlice(o["required"]); len(req) > 0 {
			return applyNullable(objectToTS(doc, o, depth+1, ctx, nameHint, mode), o)
		}
		if o["items"] != nil {
			return applyNullable(schemaAnyToTS(doc, o["items"], depth+1, ctx, joinEnumHint(nameHint, "Item"), mode)+"[]", o)
		}
		return applyNullable(tsUnknown, o)
	default:
		return applyNullable(tsUnknown, o)
	}
}

func securitySchemeToTS(s *SecurityScheme) string {
	if s == nil {
		return tsUnknown
	}
	switch s.Type {
	case "apiKey":
		return securitySchemeObject(map[string]string{
			"type": "apiKey",
			"name": s.Name,
			"in":   s.In,
		}, s.Description)
	case "http":
		fields := map[string]string{
			"type":   "http",
			"scheme": s.Scheme,
		}
		if s.BearerFormat != "" {
			fields["bearerFormat"] = s.BearerFormat
		}
		return securitySchemeObject(fields, s.Description)
	case "oauth2":
		fields := map[string]string{
			"type":  tsStringLiteralOrString("oauth2"),
			"flows": oauthFlowsToTS(s.Flows),
		}
		return securitySchemeObjectTS(fields, s.Description)
	case "openIdConnect":
		fields := map[string]string{
			"type":             "openIdConnect",
			"openIdConnectUrl": s.OpenIDConnectURL,
		}
		return securitySchemeObject(fields, s.Description)
	default:
		fields := map[string]string{}
		if s.Type != "" {
			fields["type"] = s.Type
		}
		if s.Name != "" {
			fields["name"] = s.Name
		}
		if s.In != "" {
			fields["in"] = s.In
		}
		if s.Scheme != "" {
			fields["scheme"] = s.Scheme
		}
		if s.BearerFormat != "" {
			fields["bearerFormat"] = s.BearerFormat
		}
		if s.OpenIDConnectURL != "" {
			fields["openIdConnectUrl"] = s.OpenIDConnectURL
		}
		if s.Flows != nil {
			fields["flows"] = oauthFlowsToTS(s.Flows)
		}
		if len(fields) == 0 {
			return tsRecordUnknown
		}
		return securitySchemeObjectTS(fields, s.Description)
	}
}

func securitySchemeObject(fields map[string]string, description string) string {
	out := map[string]string{}
	for k, v := range fields {
		out[k] = tsStringLiteralOrString(v)
	}
	return securitySchemeObjectTS(out, description)
}

func securitySchemeObjectTS(fields map[string]string, description string) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("{\n")
	if description != "" {
		b.WriteString("  description?: string;\n")
	}
	for _, k := range keys {
		writeTSObjectField(&b, "  ", safeProp(k), fields[k])
	}
	b.WriteString("}")
	return b.String()
}

func oauthFlowsToTS(f *OAuthFlows) string {
	if f == nil {
		return tsUnknown
	}
	fields := map[string]string{}
	if f.Implicit != nil {
		fields["implicit"] = oauthFlowToTS(f.Implicit)
	}
	if f.Password != nil {
		fields["password"] = oauthFlowToTS(f.Password)
	}
	if f.ClientCredentials != nil {
		fields["clientCredentials"] = oauthFlowToTS(f.ClientCredentials)
	}
	if f.AuthorizationCode != nil {
		fields["authorizationCode"] = oauthFlowToTS(f.AuthorizationCode)
	}
	if len(fields) == 0 {
		return "{}"
	}
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString("{\n")
	for _, k := range keys {
		writeTSObjectField(&b, "  ", safeProp(k), fields[k])
	}
	b.WriteString("}")
	return b.String()
}

func oauthFlowToTS(f *OAuthFlow) string {
	if f == nil {
		return tsUnknown
	}
	var b strings.Builder
	b.WriteString("{\n")
	if f.AuthorizationURL != "" {
		b.WriteString("  authorizationUrl: string;\n")
	}
	if f.TokenURL != "" {
		b.WriteString("  tokenUrl: string;\n")
	}
	if f.RefreshURL != "" {
		b.WriteString("  refreshUrl: string;\n")
	}
	if len(f.Scopes) > 0 {
		keys := make([]string, 0, len(f.Scopes))
		for k := range f.Scopes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		b.WriteString("  scopes: {\n")
		for _, k := range keys {
			b.WriteString("    " + safeProp(k) + ": string;\n")
		}
		b.WriteString("  };\n")
	} else {
		b.WriteString("  scopes: Record<string, string>;\n")
	}
	b.WriteString("}")
	return b.String()
}

func tsStringLiteralOrString(v string) string {
	if v == "" {
		return schemaTypeString
	}
	return strconv.Quote(v)
}

func literalToTS(v any) string {
	switch val := v.(type) {
	case string:
		return strconv.Quote(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case nil:
		return schemaTypeNull
	default:
		return ""
	}
}

func applyNullable(ts string, o map[string]any) string {
	if o == nil {
		return ts
	}
	if v, ok := o["nullable"]; ok {
		if b, ok := v.(bool); ok && b && ts != schemaTypeNull {
			return "(" + ts + " | " + schemaTypeNull + ")"
		}
	}
	return ts
}

type enumContext struct {
	enums map[string]string
	used  map[string]bool
}

func newEnumContext(enums map[string]string) *enumContext {
	if enums == nil {
		enums = map[string]string{}
	}
	used := map[string]bool{
		"Components": true,
		"Routes":     true,
	}
	for name := range enums {
		used[name] = true
	}
	return &enumContext{enums: enums, used: used}
}

func (c *enumContext) emitEnum(nameHint string, values []any, o map[string]any) string {
	if c == nil || len(values) == 0 {
		return ""
	}
	_, nullable, _ := schemaEnumValues(o)
	base := enumBaseNameFromHint(nameHint, values)
	enumName := base
	if _, exists := c.enums[enumName]; !exists {
		i := 2
		for c.used[enumName] {
			enumName = base + strconv.Itoa(i)
			i++
		}
		c.used[enumName] = true
	}
	kind, members, ok := enumMembers(values)
	if !ok || len(members) == 0 {
		return ""
	}
	if _, exists := c.enums[enumName]; !exists {
		c.enums[enumName] = emitEnumTS(enumName, members, kind)
	}
	if nullable {
		return "(" + enumName + " | null)"
	}
	return enumName
}

func joinEnumHint(base, part string) string {
	if part == "" {
		return base
	}
	if base == "" {
		return part
	}
	return base + "_" + part
}

func mediaTypeSuffix(mt string) string {
	if mt == "" {
		return ""
	}
	return "Media_" + sanitizeIdent(mt)
}

func enumBaseNameFromHint(hint string, values []any) string {
	base := camelCaseFromHint(hint)
	if base == "" {
		base = "Enum"
	}
	if suffix := enumValueSuffix(values); suffix != "" {
		base += suffix
	}
	if !strings.HasSuffix(base, "Enum") {
		base += "Enum"
	}
	return base
}

func schemaEnumValues(o map[string]any) (values []any, nullable, ok bool) {
	if o == nil {
		return nil, false, false
	}
	if v, okConst := o["const"]; okConst {
		values = []any{v}
		ok = true
	} else if ev := anySlice(o["enum"]); len(ev) > 0 {
		values = ev
		ok = true
	}
	if !ok {
		return nil, false, false
	}
	if v, ok := o["nullable"]; ok {
		if b, ok := v.(bool); ok && b {
			nullable = true
		}
	}
	return values, nullable, true
}

type enumKind int

type enumMember struct {
	Name  string
	Value string
}

func enumMembers(values []any) (enumKind, []enumMember, bool) {
	if len(values) == 0 {
		return enumInvalid, nil, false
	}
	kind := enumInvalid
	seenNames := map[string]int{}
	out := make([]enumMember, 0, len(values))
	for _, v := range values {
		var (
			memberName string
			memberVal  string
			valKind    enumKind
		)
		switch val := v.(type) {
		case string:
			valKind = enumString
			memberVal = strconv.Quote(val)
			memberName = enumMemberNameFromString(val)
		case int:
			valKind = enumNumber
			memberVal = strconv.Itoa(val)
			memberName = enumMemberNameFromNumber(memberVal)
		case int64:
			valKind = enumNumber
			memberVal = strconv.FormatInt(val, 10)
			memberName = enumMemberNameFromNumber(memberVal)
		case int32:
			valKind = enumNumber
			memberVal = strconv.FormatInt(int64(val), 10)
			memberName = enumMemberNameFromNumber(memberVal)
		case float64:
			valKind = enumNumber
			memberVal = strconv.FormatFloat(val, 'f', -1, 64)
			memberName = enumMemberNameFromNumber(memberVal)
		case float32:
			valKind = enumNumber
			memberVal = strconv.FormatFloat(float64(val), 'f', -1, 32)
			memberName = enumMemberNameFromNumber(memberVal)
		default:
			return enumInvalid, nil, false
		}
		if kind == enumInvalid {
			kind = valKind
		} else if kind != valKind {
			return enumInvalid, nil, false
		}
		if memberName == "" || !isIdent(memberName) {
			memberName = enumValuePrefix
		}
		if n := seenNames[memberName]; n > 0 {
			seenNames[memberName] = n + 1
			memberName = memberName + "_" + strconv.Itoa(n+1)
		} else {
			seenNames[memberName] = 1
		}
		out = append(out, enumMember{Name: memberName, Value: memberVal})
	}
	return kind, out, true
}

func emitEnumTS(name string, members []enumMember, _ enumKind) string {
	var b strings.Builder
	b.WriteString("export const enum " + name + " {\n")
	for _, m := range members {
		b.WriteString("  " + m.Name + " = " + m.Value + ",\n")
	}
	b.WriteString("}\n\n")
	return b.String()
}

func sanitizeIdent(s string) string {
	var b strings.Builder
	for i, r := range s {
		switch {
		case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || r == '$' || (i > 0 && r >= '0' && r <= '9'):
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			if i == 0 {
				b.WriteRune('_')
			}
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

func enumMemberNameFromString(s string) string {
	if s == "" {
		return "Empty"
	}
	var b strings.Builder
	for i, r := range s {
		switch {
		case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (i > 0 && r >= '0' && r <= '9'):
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			if i == 0 {
				b.WriteRune('_')
			}
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	name := strings.ToUpper(b.String())
	if !isIdent(name) {
		return enumValuePrefix
	}
	return name
}

func enumMemberNameFromNumber(s string) string {
	s = strings.ReplaceAll(s, "-", "NEG_")
	s = strings.ReplaceAll(s, ".", "_")
	name := enumNumberPrefix + s
	if !isIdent(name) {
		return enumValuePrefix
	}
	return name
}

func enumValueSuffix(values []any) string {
	if len(values) != 1 {
		return ""
	}
	switch v := values[0].(type) {
	case string:
		return camelCaseFromHint(v)
	case int:
		return enumValuePrefix + strconv.Itoa(v)
	case int64:
		return enumValuePrefix + strconv.FormatInt(v, 10)
	case int32:
		return enumValuePrefix + strconv.FormatInt(int64(v), 10)
	case float64:
		return enumValuePrefix + strings.ReplaceAll(strconv.FormatFloat(v, 'f', -1, 64), ".", "_")
	case float32:
		return enumValuePrefix + strings.ReplaceAll(strconv.FormatFloat(float64(v), 'f', -1, 32), ".", "_")
	default:
		return ""
	}
}

func camelCaseFromHint(hint string) string {
	parts := splitWords(hint)
	if len(parts) == 0 {
		return ""
	}
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(strings.ToLower(p[1:]))
		}
	}
	return b.String()
}

func splitWords(s string) []string {
	raw := strings.FieldsFunc(s, func(r rune) bool {
		return (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9')
	})
	out := make([]string, 0, len(raw))
	for _, token := range raw {
		if isNoiseToken(token) {
			continue
		}
		for _, w := range splitCamelToken(token) {
			if isNoiseToken(w) {
				continue
			}
			out = append(out, w)
		}
	}
	return out
}

func splitCamelToken(s string) []string {
	if s == "" {
		return nil
	}
	out := []string{}
	start := 0
	for i := 1; i < len(s); i++ {
		prev := s[i-1]
		cur := s[i]
		if prev >= 'a' && prev <= 'z' && cur >= 'A' && cur <= 'Z' {
			out = append(out, s[start:i])
			start = i
			continue
		}
		if prev >= '0' && prev <= '9' && (cur < '0' || cur > '9') {
			out = append(out, s[start:i])
			start = i
			continue
		}
		if (prev < '0' || prev > '9') && cur >= '0' && cur <= '9' {
			out = append(out, s[start:i])
			start = i
			continue
		}
	}
	out = append(out, s[start:])
	return out
}

func isNoiseToken(s string) bool {
	if s == "" {
		return true
	}
	l := strings.ToLower(s)
	switch l {
	case "allof", "anyof", "oneof", "of", "item", "media", "additional", "properties", "response", "requestbody":
		return true
	}
	if strings.HasPrefix(l, "allof") || strings.HasPrefix(l, "anyof") || strings.HasPrefix(l, "oneof") {
		if len(l) > 5 {
			allDigits := true
			for _, r := range l[5:] {
				if r < '0' || r > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				return true
			}
		}
	}
	if strings.HasPrefix(l, "of") && len(l) > 2 {
		allDigits := true
		for _, r := range l[2:] {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			return true
		}
	}
	return false
}

func includeProperty(o map[string]any, mode schemaMode) bool {
	if mode == modeDefault || o == nil {
		return true
	}
	if mode == modeInput {
		if v, ok := o["readOnly"]; ok {
			if b, ok := v.(bool); ok && b {
				return false
			}
		}
	}
	if mode == modeOutput {
		if v, ok := o["writeOnly"]; ok {
			if b, ok := v.(bool); ok && b {
				return false
			}
		}
	}
	return true
}

func writeTSObjectField(b *strings.Builder, indent, key, value string) {
	lines := strings.Split(value, "\n")
	if len(lines) == 1 {
		b.WriteString(indent + key + ": " + value + ";\n")
		return
	}
	b.WriteString(indent + key + ": " + lines[0] + "\n")
	for i := 1; i < len(lines); i++ {
		if i == len(lines)-1 {
			b.WriteString(indent + lines[i] + ";\n")
			continue
		}
		b.WriteString(indent + lines[i] + "\n")
	}
}

func schemaAnyToTS(doc *Document, v any, depth int, ctx *enumContext, nameHint string, mode schemaMode) string {
	if v == nil {
		return tsUnknown
	}
	if depth > 30 {
		return tsUnknown
	}
	if m, ok := v.(map[string]any); ok {
		if ref, ok := m["$ref"].(string); ok && ref != "" {
			return schemaToTS(doc, &RefOr[Schema]{Ref: ref}, depth+1, ctx, nameHint, mode)
		}
		r := &RefOr[Schema]{Value: &Schema{Other: m}}
		return schemaToTS(doc, r, depth+1, ctx, nameHint, mode)
	}
	return tsUnknown
}

func objectToTS(doc *Document, o map[string]any, depth int, ctx *enumContext, nameHint string, mode schemaMode) string {
	props, _ := o["properties"].(map[string]any)
	req := stringSet(anySlice(o["required"]))
	depConstraints := dependentRequiredConstraints(doc, o, props, req, depth, ctx, nameHint, mode)

	extraInfo, hasExtra := extraPropsToTS(doc, o, depth, ctx, nameHint, mode)
	extraTS := renderExtraProps(extraInfo)
	if len(props) == 0 {
		if len(req) > 0 {
			keys := make([]string, 0, len(req))
			for k := range req {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			var b strings.Builder
			b.WriteString("{\n")
			for _, k := range keys {
				b.WriteString("  " + safeProp(k) + ": " + tsUnknown + ";\n")
			}
			b.WriteString("}")
			base := b.String()
			if hasExtra && extraTS != "" {
				return "(" + base + " & " + extraTS + ")"
			}
			return base
		}
		if hasExtra && extraTS != "" {
			return extraTS
		}
		if extraInfo.additionalFalse {
			return "Record<string, never>"
		}
		return tsRecordUnknown
	}

	keys := make([]string, 0, len(props))
	for k := range props {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("{\n")
	for _, k := range keys {
		propMap, _ := props[k].(map[string]any)
		if !includeProperty(propMap, mode) {
			continue
		}
		ts := schemaAnyToTS(doc, props[k], depth+1, ctx, joinEnumHint(nameHint, k), mode)
		if req[k] {
			b.WriteString("  " + safeProp(k) + ": " + ts + ";\n")
		} else {
			b.WriteString("  " + safeProp(k) + "?: " + ts + ";\n")
		}
	}
	b.WriteString("}")
	base := b.String()
	if hasExtra && extraTS != "" {
		base = "(" + base + " & " + extraTS + ")"
	}

	if len(depConstraints) > 0 {
		base = "(" + base + " & " + strings.Join(depConstraints, " & ") + ")"
	}

	ifSchema, hasIf := o["if"]
	thenSchema, hasThen := o["then"]
	elseSchema, hasElse := o["else"]
	if hasIf || hasThen || hasElse {
		if ts := ifThenElseToTS(doc, props, base, ifSchema, thenSchema, elseSchema, depth, ctx, nameHint, mode, hasIf, hasThen, hasElse); ts != "" {
			return ts
		}
	}
	return base
}

type fieldSpec struct {
	Name     string
	TS       string
	Optional bool
}

func objectTypeFromFields(fields []fieldSpec) string {
	if len(fields) == 0 {
		return "{}"
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	var b strings.Builder
	b.WriteString("{\n")
	for _, f := range fields {
		if f.Optional {
			b.WriteString("  " + safeProp(f.Name) + "?: " + f.TS + ";\n")
		} else {
			b.WriteString("  " + safeProp(f.Name) + ": " + f.TS + ";\n")
		}
	}
	b.WriteString("}")
	return b.String()
}

func propertySchemaToTS(doc *Document, props map[string]any, name string, depth int, ctx *enumContext, nameHint string, mode schemaMode) string {
	if props == nil {
		return tsUnknown
	}
	if v, ok := props[name]; ok {
		return schemaAnyToTS(doc, v, depth+1, ctx, joinEnumHint(nameHint, name), mode)
	}
	return tsUnknown
}

func dependentRequiredConstraints(doc *Document, o, props map[string]any, req map[string]bool, depth int, ctx *enumContext, nameHint string, mode schemaMode) []string {
	dr, ok := o["dependentRequired"].(map[string]any)
	if !ok || len(dr) == 0 {
		return nil
	}

	out := make([]string, 0, len(dr))
	for key, v := range dr {
		names, ok := v.([]any)
		if !ok || len(names) == 0 {
			continue
		}
		depFields := make([]fieldSpec, 0, len(names))
		for _, it := range names {
			s, ok := it.(string)
			if !ok {
				continue
			}
			depFields = append(depFields, fieldSpec{
				Name:     s,
				TS:       propertySchemaToTS(doc, props, s, depth, ctx, nameHint, mode),
				Optional: false,
			})
		}
		if len(depFields) == 0 {
			continue
		}

		keyTS := propertySchemaToTS(doc, props, key, depth, ctx, nameHint, mode)
		if req[key] {
			fields := append([]fieldSpec{{Name: key, TS: keyTS, Optional: false}}, depFields...)
			out = append(out, objectTypeFromFields(fields))
			continue
		}

		absentFields := []fieldSpec{{Name: key, TS: tsNever, Optional: true}}
		requiredFields := append([]fieldSpec{{Name: key, TS: keyTS, Optional: false}}, depFields...)
		out = append(out, "("+objectTypeFromFields(absentFields)+" | "+objectTypeFromFields(requiredFields)+")")
	}

	return out
}

func ifThenElseToTS(doc *Document, props map[string]any, base string, ifSchema, thenSchema, elseSchema any, depth int, ctx *enumContext, nameHint string, mode schemaMode, hasIf, hasThen, hasElse bool) string {
	ifProp, ifVals, ok := extractIfPropertyValues(ifSchema)
	if ok {
		baseEnum := propertyEnumLiterals(props, ifProp)
		elseVals := subtractLiterals(baseEnum, ifVals)

		parts := []string{}
		if hasThen {
			thenPart := base
			if len(ifVals) > 0 {
				thenPart = "(" + thenPart + " & " + objectTypeFromFields([]fieldSpec{{Name: ifProp, TS: literalUnion(ifVals), Optional: false}}) + ")"
			}
			thenTS := schemaAnyToTS(doc, thenSchema, depth+1, ctx, joinEnumHint(nameHint, "Then"), mode)
			if thenTS != "" {
				thenPart = "(" + thenPart + " & " + thenTS + ")"
			}
			parts = append(parts, thenPart)
		}
		if hasElse || len(elseVals) > 0 {
			elsePart := base
			if len(elseVals) > 0 {
				elsePart = "(" + elsePart + " & " + objectTypeFromFields([]fieldSpec{{Name: ifProp, TS: literalUnion(elseVals), Optional: false}}) + ")"
			}
			if hasElse {
				elseTS := schemaAnyToTS(doc, elseSchema, depth+1, ctx, joinEnumHint(nameHint, "Else"), mode)
				if elseTS != "" {
					elsePart = "(" + elsePart + " & " + elseTS + ")"
				}
			}
			parts = append(parts, elsePart)
		}
		if len(parts) == 1 {
			return parts[0]
		}
		if len(parts) > 1 {
			return "(" + strings.Join(parts, " | ") + ")"
		}
	}

	parts := []string{}
	if hasThen {
		parts = append(parts, schemaAnyToTS(doc, thenSchema, depth+1, ctx, joinEnumHint(nameHint, "Then"), mode))
	}
	if hasElse {
		parts = append(parts, schemaAnyToTS(doc, elseSchema, depth+1, ctx, joinEnumHint(nameHint, "Else"), mode))
	}
	if len(parts) == 0 && hasIf {
		parts = append(parts, schemaAnyToTS(doc, ifSchema, depth+1, ctx, joinEnumHint(nameHint, "If"), mode))
	}
	if len(parts) == 1 {
		return "(" + base + " & " + parts[0] + ")"
	}
	if len(parts) > 1 {
		return "(" + base + " & (" + strings.Join(parts, " | ") + "))"
	}
	return ""
}

func extractIfPropertyValues(v any) (propName string, literals []string, ok bool) {
	m, ok := v.(map[string]any)
	if !ok {
		return "", nil, false
	}
	props, ok := m["properties"].(map[string]any)
	if !ok || len(props) != 1 {
		return "", nil, false
	}
	for name, raw := range props {
		prop, ok := raw.(map[string]any)
		if !ok {
			return "", nil, false
		}
		if cv, ok := prop["const"]; ok {
			if ts := literalToTS(cv); ts != "" {
				return name, []string{ts}, true
			}
			return "", nil, false
		}
		if ev := anySlice(prop["enum"]); len(ev) > 0 {
			lits := make([]string, 0, len(ev))
			for _, it := range ev {
				if ts := literalToTS(it); ts != "" {
					lits = append(lits, ts)
				}
			}
			if len(lits) > 0 {
				return name, lits, true
			}
			return "", nil, false
		}
	}
	return "", nil, false
}

func propertyEnumLiterals(props map[string]any, name string) []string {
	if props == nil {
		return nil
	}
	raw, ok := props[name]
	if !ok {
		return nil
	}
	prop, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	if cv, ok := prop["const"]; ok {
		if ts := literalToTS(cv); ts != "" {
			return []string{ts}
		}
		return nil
	}
	if ev := anySlice(prop["enum"]); len(ev) > 0 {
		out := make([]string, 0, len(ev))
		for _, it := range ev {
			if ts := literalToTS(it); ts != "" {
				out = append(out, ts)
			}
		}
		return out
	}
	return nil
}

func subtractLiterals(base, remove []string) []string {
	if len(base) == 0 {
		return nil
	}
	rm := map[string]bool{}
	for _, v := range remove {
		rm[v] = true
	}
	out := make([]string, 0, len(base))
	for _, v := range base {
		if !rm[v] {
			out = append(out, v)
		}
	}
	return out
}

func literalUnion(values []string) string {
	if len(values) == 0 {
		return tsUnknown
	}
	if len(values) == 1 {
		return values[0]
	}
	return "(" + strings.Join(values, " | ") + ")"
}

type extraPropsInfo struct {
	additionalValue   string
	patternKeyTypes   []string
	patternValueTypes []string
	additionalEnabled bool
	additionalFalse   bool
}

func extraPropsToTS(doc *Document, o map[string]any, depth int, ctx *enumContext, nameHint string, mode schemaMode) (extraPropsInfo, bool) {
	info := extraPropsInfo{}

	if pp, ok := o["patternProperties"].(map[string]any); ok && len(pp) > 0 {
		keys := make([]string, 0, len(pp))
		for k := range pp {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if m, ok := pp[k].(map[string]any); ok {
				ts := schemaAnyToTS(doc, m, depth+1, ctx, joinEnumHint(nameHint, "Pattern"+strconv.Itoa(i+1)), mode)
				if ts != "" {
					info.patternValueTypes = append(info.patternValueTypes, ts)
					info.patternKeyTypes = append(info.patternKeyTypes, patternKeyType(k))
				}
			}
		}
	}

	if ap, ok := o["additionalProperties"]; ok {
		switch v := ap.(type) {
		case bool:
			if v {
				info.additionalEnabled = true
				info.additionalValue = tsUnknown
			} else {
				info.additionalFalse = true
			}
		case map[string]any:
			ts := schemaAnyToTS(doc, v, depth+1, ctx, joinEnumHint(nameHint, "AdditionalProperties"), mode)
			if ts != "" {
				info.additionalEnabled = true
				info.additionalValue = ts
			}
		}
	}

	hasExtra := len(info.patternValueTypes) > 0 || info.additionalEnabled
	return info, hasExtra
}

func renderExtraProps(info extraPropsInfo) string {
	hasPattern := len(info.patternValueTypes) > 0
	hasAdditional := info.additionalEnabled
	if !hasPattern && !hasAdditional {
		return ""
	}

	var patternTS string
	var patternValue string
	if hasPattern {
		patternValue = unionTypes(info.patternValueTypes)
		keyUnion := unionKeyTypes(info.patternKeyTypes)
		if keyUnion == "" {
			keyUnion = schemaTypeString
		}
		if keyUnion == schemaTypeString {
			patternTS = "Record<" + schemaTypeString + ", " + patternValue + ">"
		} else {
			patternTS = "{ [K in " + keyUnion + "]?: " + patternValue + " }"
		}
	}

	if hasAdditional {
		union := info.additionalValue
		if hasPattern {
			union = unionTypes([]string{union, patternValue})
		}
		base := "Record<" + schemaTypeString + ", " + union + ">"
		if hasPattern {
			return "(" + patternTS + " & " + base + ")"
		}
		return base
	}

	return patternTS
}

func unionTypes(items []string) string {
	seen := map[string]bool{}
	out := make([]string, 0, len(items))
	for _, it := range items {
		if it == "" {
			continue
		}
		if it == tsUnknown {
			return tsUnknown
		}
		if !seen[it] {
			seen[it] = true
			out = append(out, it)
		}
	}
	if len(out) == 0 {
		return tsUnknown
	}
	if len(out) == 1 {
		return out[0]
	}
	return "(" + strings.Join(out, " | ") + ")"
}

func unionKeyTypes(items []string) string {
	seen := map[string]bool{}
	out := make([]string, 0, len(items))
	for _, it := range items {
		if it == "" {
			continue
		}
		if it == schemaTypeString {
			return schemaTypeString
		}
		if !seen[it] {
			seen[it] = true
			out = append(out, it)
		}
	}
	if len(out) == 0 {
		return ""
	}
	if len(out) == 1 {
		return out[0]
	}
	return "(" + strings.Join(out, " | ") + ")"
}

func patternKeyType(pattern string) string {
	if pattern == "^[0-9]+$" || pattern == "^\\d+$" {
		return "`${number}`"
	}
	if !strings.HasPrefix(pattern, "^") {
		return schemaTypeString
	}
	p := strings.TrimPrefix(pattern, "^")
	exact := false
	if strings.HasSuffix(p, "$") {
		exact = true
		p = strings.TrimSuffix(p, "$")
	}
	if strings.HasSuffix(p, ".*") {
		p = strings.TrimSuffix(p, ".*")
	} else if strings.HasSuffix(p, ".+") {
		p = strings.TrimSuffix(p, ".+")
	}
	if p == "" || strings.Contains(p, "`") {
		return schemaTypeString
	}
	if strings.ContainsAny(p, "[]()|+*?.\\") {
		return schemaTypeString
	}
	if exact {
		return strconv.Quote(p)
	}
	return "`" + p + "${string}`"
}

func anySlice(v any) []any {
	a, _ := v.([]any)
	return a
}

func stringSet(v []any) map[string]bool {
	out := map[string]bool{}
	for _, it := range v {
		if s, ok := it.(string); ok {
			out[s] = true
		}
	}
	return out
}

func safeProp(k string) string {
	if isIdent(k) {
		return k
	}
	return strconv.Quote(k)
}

func componentSchemaRef(name string) string {
	return "Components[\"schemas\"][" + strconv.Quote(name) + "]"
}

func componentResponseRef(name string) string {
	return "Components[\"responses\"][" + strconv.Quote(name) + "]"
}

func componentRequestBodyRef(name string) string {
	return "Components[\"requestBodies\"][" + strconv.Quote(name) + "]"
}

func componentParameterRef(name string) string {
	return "Components[\"parameters\"][" + strconv.Quote(name) + "]"
}

func componentHeaderRef(name string) string {
	return "Components[\"headers\"][" + strconv.Quote(name) + "]"
}

func isIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !isIdentStart(r) {
				return false
			}
		} else {
			if !isIdentPart(r) {
				return false
			}
		}
	}
	return true
}

func isIdentStart(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || r == '$'
}

func isIdentPart(r rune) bool {
	return isIdentStart(r) || (r >= '0' && r <= '9')
}
