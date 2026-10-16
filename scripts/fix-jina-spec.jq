# Rewrites the pydantic-style schemas in Jina AI's OpenAPI 3.1 document into a
# form oapi-codegen can generate.
#
# 1. `anyOf: [X, {type: null}]` becomes X plus `nullable: true`.
# 2. A property typed `anyOf: [<enum>, array of the same <enum>]` gets the enum
#    hoisted into a shared component named after the property, and both
#    branches point at it. Inline, oapi-codegen names the scalar enum and the
#    array item enum the same and stops with a duplicate type name.
# 3. An inline JSON request body with a title becomes a component named after
#    the title (`Embedding Request` becomes `EmbeddingRequest`). Inline, the
#    body type is a defined type over the oneOf union, which drops the
#    union's MarshalJSON method, so the client would send `{}`. Through a
#    component the body type is an alias and keeps the method.

def isnull: . == {"type": "null"};

def pascal: split("_") | map((.[:1] | ascii_upcase) + .[1:]) | join("");

def drop_null_branches:
  walk(
    if type == "object" and (.anyOf | type) == "array" and any(.anyOf[]; isnull) then
      (.anyOf | map(select(isnull | not))) as $rest
      | if ($rest | length) == 1 then del(.anyOf) + $rest[0] + {nullable: true}
        else (.anyOf = $rest) + {nullable: true} end
    else . end
  );

def scalar_or_list_enum:
  type == "object"
  and (.anyOf | type) == "array"
  and (.anyOf | length) == 2
  and (.anyOf[0].enum != null)
  and .anyOf[1].type == "array"
  and .anyOf[1].items == .anyOf[0];

def hoist_enums:
  . as $doc
  | [
      $doc.components.schemas | to_entries[]
      | .key as $schema
      | (.value.properties // {}) | to_entries[]
      | select(.value | scalar_or_list_enum)
      | {schema: $schema, prop: .key, name: (.key | pascal), enum: .value.anyOf[0]}
    ] as $hits
  | if ($hits | group_by(.name) | any(map(.enum) | unique | length > 1)) then
      error("hoisted enums with one name but different values")
    else . end
  | reduce $hits[] as $h (.;
      .components.schemas[$h.name] = $h.enum
      | .components.schemas[$h.schema].properties[$h.prop].anyOf =
        [{"$ref": "#/components/schemas/\($h.name)"},
         {type: "array", items: {"$ref": "#/components/schemas/\($h.name)"}}]
    );

def hoist_request_bodies:
  reduce (
    paths(objects and has("requestBody"))
    | select(length == 3 and .[0] == "paths")
  ) as $op (.;
    getpath($op + ["requestBody", "content", "application/json", "schema"]) as $schema
    | if ($schema | type) == "object" and ($schema.title | type) == "string"
        and ($schema | has("$ref") | not) then
        ($schema.title | gsub("[^A-Za-z0-9]"; "")) as $name
        | if .components.schemas[$name] != null then
            error("request body component \($name) already exists")
          else
            .components.schemas[$name] = $schema
            | setpath($op + ["requestBody", "content", "application/json", "schema"];
                {"$ref": "#/components/schemas/\($name)"})
          end
      else . end
  );

drop_null_branches | hoist_enums | hoist_request_bodies
