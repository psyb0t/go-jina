# Rewrites Jina AI's Reader OpenAPI document into a form oapi-codegen can
# generate.
#
# 1. The path parameter `{url}` becomes `{targetUrl}`. A path parameter is
#    not sent by name, so the request is unchanged, but a Go parameter named
#    `url` shadows the net/url package in the generated client and it does
#    not compile.
# 2. `additionalProperties` schemas for "any value" carry `nullable: true`
#    without a type. oapi-codegen then declares the map as
#    map[string]*interface{} but writes interface{} values into it, which does
#    not compile. An untyped schema already allows null, so the flag goes.
# 3. Query and header parameters with an untyped "any value" schema (such as
#    `jsonSchema`) become `type: string`. Untyped, oapi-codegen makes the
#    field a plain interface{} and always serializes it, so a request that
#    leaves it unset panics on the nil value. On the wire a query or header
#    value is a string anyway; a JSON value goes in as its encoded text.

def untyped:
  (has("type") or has("$ref") or has("oneOf") or has("anyOf") or has("allOf")) | not;

def rename_url_param:
  .paths |= with_entries(
    if .key | contains("{url}") then
      .key |= sub("\\{url\\}"; "{targetUrl}")
      | .value |= map_values(
          if type == "object" and (.parameters | type) == "array" then
            .parameters |= map(
              if .in == "path" and .name == "url" then .name = "targetUrl" else . end
            )
          else . end
        )
    else . end
  );

def drop_untyped_nullable_maps:
  walk(
    if type == "object" and (.additionalProperties | type) == "object"
      and (.additionalProperties | untyped) then
      .additionalProperties |= del(.nullable)
    else . end
  );

def type_untyped_params:
  .paths |= map_values(map_values(
    if type == "object" and (.parameters | type) == "array" then
      .parameters |= map(
        if (.in == "query" or .in == "header") and (.schema | type) == "object"
          and (.schema | untyped) then
          .schema |= (del(.nullable) + {type: "string"})
        else . end
      )
    else . end
  ));

rename_url_param | drop_untyped_nullable_maps | type_untyped_params
