package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func propertiesFromType(tp apiSpec.Type) (spec.SchemaProperties, []string) {
	var (
		properties     = map[string]spec.Schema{}
		requiredFields []string
	)
	switch val := tp.(type) {
	case apiSpec.PointerType:
		return propertiesFromType(val.Type)
	case apiSpec.ArrayType:
		return propertiesFromType(val.Value)
	case apiSpec.DefineStruct, apiSpec.NestedStruct:
		rangeMemberAndDo(val, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
			var (
				jsonTagString                      = member.Name
				minimum, maximum                   *float64
				exclusiveMinimum, exclusiveMaximum bool
				example, defaultValue              any
				enum                               []any
			)
			jsonTag, _ := tag.Get(tagJson)
			if jsonTag != nil {
				jsonTagString = jsonTag.Name
				minimum, maximum, exclusiveMinimum, exclusiveMaximum = rangeValueFromOptions(jsonTag.Options)
				example = exampleValueFromOptions(jsonTag.Options, member.Type)
				defaultValue = defValueFromOptions(jsonTag.Options, member.Type)
				enum = enumsValueFromOptions(jsonTag.Options)
			}

			if required {
				requiredFields = append(requiredFields, jsonTagString)
			}
			var schema = spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description:          formatComment(member.Comment),
					Type:                 typeFromGoType(member.Type),
					Default:              defaultValue,
					Maximum:              maximum,
					ExclusiveMaximum:     exclusiveMaximum,
					Minimum:              minimum,
					ExclusiveMinimum:     exclusiveMinimum,
					Enum:                 enum,
					AdditionalProperties: mapFromGoType(member.Type),
				},
			}
			switch sampleTypeFromGoType(member.Type) {
			case swaggerTypeArray:
				schema.Items = itemsFromGoType(member.Type)
			case swaggerTypeObject:
				p, r := propertiesFromType(member.Type)
				schema.Properties = p
				schema.Required = r
			}

			properties[jsonTagString] = schema
		})
	}

	return properties, requiredFields
}
