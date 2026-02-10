package schema

import (
	"encoding/json"

	"go.yaml.in/yaml/v3"
)

type Extensions map[string]any

type RefOr[T any] struct {
	Value *T
	Ref   string
}

func (r *RefOr[T]) UnmarshalYAML(node *yaml.Node) error {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			if k != nil && k.Value == "$ref" {
				var ref string
				if err := v.Decode(&ref); err != nil {
					return err
				}
				r.Ref = ref
				r.Value = nil
				return nil
			}
		}
	}
	var val T
	if err := node.Decode(&val); err != nil {
		return err
	}
	r.Ref = ""
	r.Value = &val
	return nil
}

func (r *RefOr[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err == nil {
		if raw, ok := obj["$ref"]; ok {
			var ref string
			if err := json.Unmarshal(raw, &ref); err != nil {
				return err
			}
			r.Ref = ref
			r.Value = nil
			return nil
		}
	}
	var val T
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	r.Ref = ""
	r.Value = &val
	return nil
}

type Document struct {
	Info              Info                       `yaml:"info" json:"info"`
	Paths             Paths                      `yaml:"paths,omitempty" json:"paths,omitempty"`
	Webhooks          map[string]RefOr[PathItem] `yaml:"webhooks,omitempty" json:"webhooks,omitempty"`
	Components        *Components                `yaml:"components,omitempty" json:"components,omitempty"`
	ExternalDocs      *ExternalDocumentation     `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Extensions        Extensions                 `yaml:",inline" json:"-"`
	OpenAPI           string                     `yaml:"openapi" json:"openapi"`
	JSONSchemaDialect string                     `yaml:"jsonSchemaDialect,omitempty" json:"jsonSchemaDialect,omitempty"`
	Servers           []Server                   `yaml:"servers,omitempty" json:"servers,omitempty"`
	Security          []SecurityRequirement      `yaml:"security,omitempty" json:"security,omitempty"`
	Tags              []Tag                      `yaml:"tags,omitempty" json:"tags,omitempty"`
}

type Info struct {
	Contact        *Contact   `yaml:"contact,omitempty" json:"contact,omitempty"`
	License        *License   `yaml:"license,omitempty" json:"license,omitempty"`
	Extensions     Extensions `yaml:",inline" json:"-"`
	Title          string     `yaml:"title" json:"title"`
	Summary        string     `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description    string     `yaml:"description,omitempty" json:"description,omitempty"`
	TermsOfService string     `yaml:"termsOfService,omitempty" json:"termsOfService,omitempty"`
	Version        string     `yaml:"version" json:"version"`
}

type Contact struct {
	Extensions Extensions `yaml:",inline" json:"-"`
	Name       string     `yaml:"name,omitempty" json:"name,omitempty"`
	URL        string     `yaml:"url,omitempty" json:"url,omitempty"`
	Email      string     `yaml:"email,omitempty" json:"email,omitempty"`
}

type License struct {
	Extensions Extensions `yaml:",inline" json:"-"`
	Name       string     `yaml:"name" json:"name"`
	Identifier string     `yaml:"identifier,omitempty" json:"identifier,omitempty"`
	URL        string     `yaml:"url,omitempty" json:"url,omitempty"`
}

type Server struct {
	Variables   map[string]ServerVariable `yaml:"variables,omitempty" json:"variables,omitempty"`
	Extensions  Extensions                `yaml:",inline" json:"-"`
	URL         string                    `yaml:"url" json:"url"`
	Description string                    `yaml:"description,omitempty" json:"description,omitempty"`
}

type ServerVariable struct {
	Extensions  Extensions `yaml:",inline" json:"-"`
	Default     string     `yaml:"default" json:"default"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Enum        []string   `yaml:"enum,omitempty" json:"enum,omitempty"`
}

type Paths map[string]RefOr[PathItem]

type PathItem struct {
	Delete      *Operation         `yaml:"delete,omitempty" json:"delete,omitempty"`
	Trace       *Operation         `yaml:"trace,omitempty" json:"trace,omitempty"`
	Extensions  Extensions         `yaml:",inline" json:"-"`
	Get         *Operation         `yaml:"get,omitempty" json:"get,omitempty"`
	Put         *Operation         `yaml:"put,omitempty" json:"put,omitempty"`
	Post        *Operation         `yaml:"post,omitempty" json:"post,omitempty"`
	Head        *Operation         `yaml:"head,omitempty" json:"head,omitempty"`
	Options     *Operation         `yaml:"options,omitempty" json:"options,omitempty"`
	Patch       *Operation         `yaml:"patch,omitempty" json:"patch,omitempty"`
	Summary     string             `yaml:"summary,omitempty" json:"summary,omitempty"`
	Ref         string             `yaml:"$ref,omitempty" json:"$ref,omitempty"`
	Description string             `yaml:"description,omitempty" json:"description,omitempty"`
	Servers     []Server           `yaml:"servers,omitempty" json:"servers,omitempty"`
	Parameters  []RefOr[Parameter] `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

type Operation struct {
	RequestBody  *RefOr[RequestBody]        `yaml:"requestBody,omitempty" json:"requestBody,omitempty"`
	ExternalDocs *ExternalDocumentation     `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Responses    Responses                  `yaml:"responses" json:"responses"`
	Callbacks    map[string]RefOr[Callback] `yaml:"callbacks,omitempty" json:"callbacks,omitempty"`
	Extensions   Extensions                 `yaml:",inline" json:"-"`
	Summary      string                     `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description  string                     `yaml:"description,omitempty" json:"description,omitempty"`
	OperationID  string                     `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Parameters   []RefOr[Parameter]         `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Tags         []string                   `yaml:"tags,omitempty" json:"tags,omitempty"`
	Security     []SecurityRequirement      `yaml:"security,omitempty" json:"security,omitempty"`
	Servers      []Server                   `yaml:"servers,omitempty" json:"servers,omitempty"`
	Deprecated   bool                       `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
}

type ExternalDocumentation struct {
	Extensions  Extensions `yaml:",inline" json:"-"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	URL         string     `yaml:"url" json:"url"`
}

type Parameter struct {
	Example         any                       `yaml:"example,omitempty" json:"example,omitempty"`
	Explode         *bool                     `yaml:"explode,omitempty" json:"explode,omitempty"`
	Extensions      Extensions                `yaml:",inline" json:"-"`
	Content         map[string]MediaType      `yaml:"content,omitempty" json:"content,omitempty"`
	Examples        map[string]RefOr[Example] `yaml:"examples,omitempty" json:"examples,omitempty"`
	Schema          *RefOr[Schema]            `yaml:"schema,omitempty" json:"schema,omitempty"`
	Style           string                    `yaml:"style,omitempty" json:"style,omitempty"`
	Name            string                    `yaml:"name,omitempty" json:"name,omitempty"`
	Description     string                    `yaml:"description,omitempty" json:"description,omitempty"`
	In              string                    `yaml:"in,omitempty" json:"in,omitempty"`
	AllowReserved   bool                      `yaml:"allowReserved,omitempty" json:"allowReserved,omitempty"`
	AllowEmptyValue bool                      `yaml:"allowEmptyValue,omitempty" json:"allowEmptyValue,omitempty"`
	Deprecated      bool                      `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
	Required        bool                      `yaml:"required,omitempty" json:"required,omitempty"`
}

type RequestBody struct {
	Content     map[string]MediaType `yaml:"content,omitempty" json:"content,omitempty"`
	Extensions  Extensions           `yaml:",inline" json:"-"`
	Description string               `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool                 `yaml:"required,omitempty" json:"required,omitempty"`
}

type MediaType struct {
	Schema     *RefOr[Schema]            `yaml:"schema,omitempty" json:"schema,omitempty"`
	Example    any                       `yaml:"example,omitempty" json:"example,omitempty"`
	Examples   map[string]RefOr[Example] `yaml:"examples,omitempty" json:"examples,omitempty"`
	Encoding   map[string]Encoding       `yaml:"encoding,omitempty" json:"encoding,omitempty"`
	Extensions Extensions                `yaml:",inline" json:"-"`
}

type Encoding struct {
	Headers       map[string]RefOr[Header] `yaml:"headers,omitempty" json:"headers,omitempty"`
	Explode       *bool                    `yaml:"explode,omitempty" json:"explode,omitempty"`
	Extensions    Extensions               `yaml:",inline" json:"-"`
	ContentType   string                   `yaml:"contentType,omitempty" json:"contentType,omitempty"`
	Style         string                   `yaml:"style,omitempty" json:"style,omitempty"`
	AllowReserved bool                     `yaml:"allowReserved,omitempty" json:"allowReserved,omitempty"`
}

type Responses map[string]RefOr[Response]

type Response struct {
	Headers     map[string]RefOr[Header] `yaml:"headers,omitempty" json:"headers,omitempty"`
	Content     map[string]MediaType     `yaml:"content,omitempty" json:"content,omitempty"`
	Links       map[string]RefOr[Link]   `yaml:"links,omitempty" json:"links,omitempty"`
	Extensions  Extensions               `yaml:",inline" json:"-"`
	Description string                   `yaml:"description,omitempty" json:"description,omitempty"`
}

type Header struct {
	Example     any                       `yaml:"example,omitempty" json:"example,omitempty"`
	Explode     *bool                     `yaml:"explode,omitempty" json:"explode,omitempty"`
	Schema      *RefOr[Schema]            `yaml:"schema,omitempty" json:"schema,omitempty"`
	Examples    map[string]RefOr[Example] `yaml:"examples,omitempty" json:"examples,omitempty"`
	Content     map[string]MediaType      `yaml:"content,omitempty" json:"content,omitempty"`
	Extensions  Extensions                `yaml:",inline" json:"-"`
	Description string                    `yaml:"description,omitempty" json:"description,omitempty"`
	Style       string                    `yaml:"style,omitempty" json:"style,omitempty"`
	Required    bool                      `yaml:"required,omitempty" json:"required,omitempty"`
	Deprecated  bool                      `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
}

type Example struct {
	Value         any        `yaml:"value,omitempty" json:"value,omitempty"`
	Extensions    Extensions `yaml:",inline" json:"-"`
	Summary       string     `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description   string     `yaml:"description,omitempty" json:"description,omitempty"`
	ExternalValue string     `yaml:"externalValue,omitempty" json:"externalValue,omitempty"`
}

type Link struct {
	RequestBody  any            `yaml:"requestBody,omitempty" json:"requestBody,omitempty"`
	Parameters   map[string]any `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Server       *Server        `yaml:"server,omitempty" json:"server,omitempty"`
	Extensions   Extensions     `yaml:",inline" json:"-"`
	OperationRef string         `yaml:"operationRef,omitempty" json:"operationRef,omitempty"`
	OperationID  string         `yaml:"operationId,omitempty" json:"operationId,omitempty"`
	Description  string         `yaml:"description,omitempty" json:"description,omitempty"`
}

type Callback map[string]PathItem

type Components struct {
	Schemas         map[string]Schema                `yaml:"schemas,omitempty" json:"schemas,omitempty"`
	Responses       map[string]RefOr[Response]       `yaml:"responses,omitempty" json:"responses,omitempty"`
	Parameters      map[string]RefOr[Parameter]      `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Examples        map[string]RefOr[Example]        `yaml:"examples,omitempty" json:"examples,omitempty"`
	RequestBodies   map[string]RefOr[RequestBody]    `yaml:"requestBodies,omitempty" json:"requestBodies,omitempty"`
	Headers         map[string]RefOr[Header]         `yaml:"headers,omitempty" json:"headers,omitempty"`
	SecuritySchemes map[string]RefOr[SecurityScheme] `yaml:"securitySchemes,omitempty" json:"securitySchemes,omitempty"`
	Links           map[string]RefOr[Link]           `yaml:"links,omitempty" json:"links,omitempty"`
	Callbacks       map[string]RefOr[Callback]       `yaml:"callbacks,omitempty" json:"callbacks,omitempty"`
	PathItems       map[string]RefOr[PathItem]       `yaml:"pathItems,omitempty" json:"pathItems,omitempty"`
	Extensions      Extensions                       `yaml:",inline" json:"-"`
}

type SecurityRequirement map[string][]string

type SecurityScheme struct {
	Flows            *OAuthFlows `yaml:"flows,omitempty" json:"flows,omitempty"`
	Extensions       Extensions  `yaml:",inline" json:"-"`
	Type             string      `yaml:"type,omitempty" json:"type,omitempty"`
	Description      string      `yaml:"description,omitempty" json:"description,omitempty"`
	Name             string      `yaml:"name,omitempty" json:"name,omitempty"`
	In               string      `yaml:"in,omitempty" json:"in,omitempty"`
	Scheme           string      `yaml:"scheme,omitempty" json:"scheme,omitempty"`
	BearerFormat     string      `yaml:"bearerFormat,omitempty" json:"bearerFormat,omitempty"`
	OpenIDConnectURL string      `yaml:"openIdConnectUrl,omitempty" json:"openIdConnectUrl,omitempty"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `yaml:"implicit,omitempty" json:"implicit,omitempty"`
	Password          *OAuthFlow `yaml:"password,omitempty" json:"password,omitempty"`
	ClientCredentials *OAuthFlow `yaml:"clientCredentials,omitempty" json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `yaml:"authorizationCode,omitempty" json:"authorizationCode,omitempty"`
	Extensions        Extensions `yaml:",inline" json:"-"`
}

type OAuthFlow struct {
	Scopes           map[string]string `yaml:"scopes,omitempty" json:"scopes,omitempty"`
	Extensions       Extensions        `yaml:",inline" json:"-"`
	AuthorizationURL string            `yaml:"authorizationUrl,omitempty" json:"authorizationUrl,omitempty"`
	TokenURL         string            `yaml:"tokenUrl,omitempty" json:"tokenUrl,omitempty"`
	RefreshURL       string            `yaml:"refreshUrl,omitempty" json:"refreshUrl,omitempty"`
}

type Tag struct {
	ExternalDocs *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Extensions   Extensions             `yaml:",inline" json:"-"`
	Name         string                 `yaml:"name" json:"name"`
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
}

type Schema struct {
	Example       any                    `yaml:"example,omitempty" json:"example,omitempty"`
	Discriminator *Discriminator         `yaml:"discriminator,omitempty" json:"discriminator,omitempty"`
	XML           *XML                   `yaml:"xml,omitempty" json:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `yaml:"externalDocs,omitempty" json:"externalDocs,omitempty"`
	Other         map[string]any         `yaml:",inline" json:"-"`
	Deprecated    bool                   `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
}

func (s *Schema) UnmarshalJSON(data []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	type schemaAlias struct {
		Example       any                    `json:"example,omitempty"`
		Discriminator *Discriminator         `json:"discriminator,omitempty"`
		XML           *XML                   `json:"xml,omitempty"`
		ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty"`
		Deprecated    bool                   `json:"deprecated,omitempty"`
	}
	var alias schemaAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	s.Example = alias.Example
	s.Discriminator = alias.Discriminator
	s.XML = alias.XML
	s.ExternalDocs = alias.ExternalDocs
	s.Deprecated = alias.Deprecated
	s.Other = obj
	return nil
}

type Discriminator struct {
	Mapping      map[string]string `yaml:"mapping,omitempty" json:"mapping,omitempty"`
	Extensions   Extensions        `yaml:",inline" json:"-"`
	PropertyName string            `yaml:"propertyName,omitempty" json:"propertyName,omitempty"`
}

type XML struct {
	Extensions Extensions `yaml:",inline" json:"-"`
	Name       string     `yaml:"name,omitempty" json:"name,omitempty"`
	Namespace  string     `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Prefix     string     `yaml:"prefix,omitempty" json:"prefix,omitempty"`
	Attribute  bool       `yaml:"attribute,omitempty" json:"attribute,omitempty"`
	Wrapped    bool       `yaml:"wrapped,omitempty" json:"wrapped,omitempty"`
}
