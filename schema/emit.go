package schema

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

func GeneratedHeader(generator, openAPIVersion string, generatedAt time.Time) string {
	var b strings.Builder

	b.WriteString(headerStart())
	b.WriteString(" * @Warning: THIS FILE IS AUTO-GENERATED - DO NOT EDIT\n")
	b.WriteString(" *\n")
	b.WriteString(" * Generator: " + generator + "\n")
	if openAPIVersion != "" {
		b.WriteString(" * OpenAPI version: " + openAPIVersion + "\n")
	}
	b.WriteString(" * Generated at: " + generatedAt.UTC().Format(time.RFC3339) + "\n")
	b.WriteString(" */\n\n")

	return b.String()
}

func headerStart() string {
	return "/*\n"
}

func EmitTypesFromIR(ir *IR) string {
	return EmitTypesFromIRAt(ir, time.Now(), "", "")
}

func EmitTypesFromIRAt(ir *IR, generatedAt time.Time, cliVersion, openAPIVersion string) string {
	var b strings.Builder

	generator := "openapi-tsgen"
	if cliVersion != "" {
		generator = generator + "@" + cliVersion
	}

	b.WriteString(GeneratedHeader(generator, openAPIVersion, generatedAt))

	writeEnums(&b, ir)
	writeServers(&b, ir)
	writeComponents(&b, ir)
	writeRoutes(&b, ir)
	writeWebhooks(&b, ir)
	return b.String()
}

func writeComponents(b *strings.Builder, ir *IR) {
	if len(ir.ComponentsSchemas) == 0 &&
		len(ir.ComponentsResponses) == 0 &&
		len(ir.ComponentsRequestBody) == 0 &&
		len(ir.ComponentsParameters) == 0 &&
		len(ir.ComponentsHeaders) == 0 &&
		len(ir.ComponentsSecuritySchemes) == 0 {
		return
	}

	b.WriteString("export type Components = {\n")
	writeComponentSection(b, "schemas", ir.ComponentsSchemas)
	writeComponentSection(b, "responses", ir.ComponentsResponses)
	writeComponentSection(b, "requestBodies", ir.ComponentsRequestBody)
	writeComponentSection(b, "parameters", ir.ComponentsParameters)
	writeComponentSection(b, "headers", ir.ComponentsHeaders)
	writeComponentSection(b, "securitySchemes", ir.ComponentsSecuritySchemes)
	b.WriteString("};\n\n")
}

func writeEnums(b *strings.Builder, ir *IR) {
	if len(ir.Enums) == 0 {
		return
	}
	keys := make([]string, 0, len(ir.Enums))
	for k := range ir.Enums {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		b.WriteString(ir.Enums[k])
	}
}

func writeComponentSection(b *strings.Builder, label string, values map[string]string) {
	if len(values) == 0 {
		return
	}
	b.WriteString("  " + label + ": {\n")
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		key := safeTSKey(k)
		writeTSField(b, "    ", key, values[k])
	}
	b.WriteString("  };\n")
}

func writeRoutes(b *strings.Builder, ir *IR) {
	if len(ir.Paths) == 0 {
		return
	}
	writePathItems(b, "Routes", ir.Paths)
}

func writeWebhooks(b *strings.Builder, ir *IR) {
	if len(ir.Webhooks) == 0 {
		return
	}
	writePathItems(b, "Webhooks", ir.Webhooks)
}

func writePathItems(b *strings.Builder, label string, items map[string]IRPathItem) {
	b.WriteString("export type " + label + " = {\n")

	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		item := items[key]
		b.WriteString("  " + strconv.Quote(key) + ": {\n")

		methods := make([]string, 0, len(item.Ops))
		for m := range item.Ops {
			methods = append(methods, m)
		}
		sort.Strings(methods)

		for _, method := range methods {
			op := item.Ops[method]
			b.WriteString("    " + method + ": {\n")

			writeParamsBlock(b, "params", op.PathParams)
			writeParamsBlock(b, "query", op.QueryParams)
			writeParamsBlock(b, "headers", op.HeaderParams)
			writeParamsBlock(b, "cookies", op.CookieParams)
			if len(op.Security) > 0 {
				writeTSField(b, "      ", "security", securityRequirementsToTS(op.Security))
			}
			if len(op.Servers) > 0 {
				writeTSField(b, "      ", "servers", serversToTS(op.Servers))
			}

			if op.RequestBody != tsNever {
				writeTSField(b, "      ", "requestBody", op.RequestBody)
			}
			b.WriteString("      responses: {\n")
			codes := make([]string, 0, len(op.Responses))
			for c := range op.Responses {
				codes = append(codes, c)
			}
			sort.Slice(codes, func(i, j int) bool {
				ai, aok := parseStatusCode(codes[i])
				bi, bok := parseStatusCode(codes[j])
				if aok && bok {
					return ai < bi
				}
				if aok != bok {
					return aok
				}
				return codes[i] < codes[j]
			})
			for _, c := range codes {
				key := c
				if _, ok := parseStatusCode(c); !ok && c != "default" {
					key = strconv.Quote(c)
				}
				writeTSField(b, "        ", key, op.Responses[c])
			}
			b.WriteString("      };\n")

			b.WriteString("    };\n")
		}

		b.WriteString("  };\n")
	}

	b.WriteString("};\n\n")
}

func writeServers(b *strings.Builder, ir *IR) {
	if len(ir.Servers) == 0 {
		return
	}
	b.WriteString("export type Servers = " + serversToTS(ir.Servers) + ";\n\n")
}

func serversToTS(servers []Server) string {
	if len(servers) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(servers))
	for _, s := range servers {
		parts = append(parts, serverToTS(s))
	}
	if len(parts) == 1 {
		return parts[0] + "[]"
	}
	return "(" + strings.Join(parts, " | ") + ")[]"
}

func serverToTS(s Server) string {
	fields := []fieldSpec{
		{Name: "url", TS: strconv.Quote(s.URL), Optional: false},
	}
	if s.Description != "" {
		fields = append(fields, fieldSpec{Name: "description", TS: strconv.Quote(s.Description), Optional: false})
	}
	if len(s.Variables) > 0 {
		fields = append(fields, fieldSpec{Name: "variables", TS: serverVariablesToTS(s.Variables), Optional: false})
	}
	return objectTypeFromFields(fields)
}

func serverVariablesToTS(vars map[string]ServerVariable) string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fields := make([]fieldSpec, 0, len(keys))
	for _, k := range keys {
		v := vars[k]
		valueFields := []fieldSpec{
			{Name: "default", TS: strconv.Quote(v.Default), Optional: false},
		}
		if v.Description != "" {
			valueFields = append(valueFields, fieldSpec{Name: "description", TS: strconv.Quote(v.Description), Optional: false})
		}
		if len(v.Enum) > 0 {
			enumVals := make([]string, 0, len(v.Enum))
			for _, ev := range v.Enum {
				enumVals = append(enumVals, strconv.Quote(ev))
			}
			valueFields = append(valueFields, fieldSpec{Name: "enum", TS: literalUnion(enumVals) + "[]", Optional: false})
		}
		fields = append(fields, fieldSpec{Name: k, TS: objectTypeFromFields(valueFields), Optional: false})
	}
	return objectTypeFromFields(fields)
}

func securityRequirementsToTS(reqs []SecurityRequirement) string {
	if len(reqs) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(reqs))
	for _, req := range reqs {
		parts = append(parts, securityRequirementToTS(req))
	}
	if len(parts) == 1 {
		return parts[0] + "[]"
	}
	return "(" + strings.Join(parts, " | ") + ")[]"
}

func securityRequirementToTS(req SecurityRequirement) string {
	if len(req) == 0 {
		return tsEmptyObject
	}
	keys := make([]string, 0, len(req))
	for k := range req {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fields := make([]fieldSpec, 0, len(keys))
	for _, k := range keys {
		scopes := req[k]
		fields = append(fields, fieldSpec{Name: k, TS: scopesToTS(scopes), Optional: false})
	}
	return objectTypeFromFields(fields)
}

func scopesToTS(scopes []string) string {
	if len(scopes) == 0 {
		return "string[]"
	}
	vals := make([]string, 0, len(scopes))
	for _, s := range scopes {
		vals = append(vals, strconv.Quote(s))
	}
	return literalUnion(vals) + "[]"
}

func writeTSField(b *strings.Builder, indent, key, value string) {
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

func safeTSKey(k string) string {
	if isTSIdent(k) {
		return k
	}
	return strconv.Quote(k)
}

func writeParamsBlock(b *strings.Builder, label string, params map[string]paramResolved) {
	if len(params) == 0 {
		return
	}
	b.WriteString("      " + label + ": {\n")
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		prop := safeTSKey(k)
		if params[k].Required {
			b.WriteString("        " + prop + ": " + params[k].TS + ";\n")
		} else {
			b.WriteString("        " + prop + "?: " + params[k].TS + ";\n")
		}
	}
	b.WriteString("      };\n")
}

func parseStatusCode(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func isTSIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !isTSIdentStart(r) {
				return false
			}
		} else {
			if !isTSIdentPart(r) {
				return false
			}
		}
	}
	return true
}

func isTSIdentStart(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || r == '$'
}

func isTSIdentPart(r rune) bool {
	return isTSIdentStart(r) || (r >= '0' && r <= '9')
}
