# GSON simple golang json encoder/decoder

Simple light-weigth json encoder/decoder, that operates on interface{} types

# Methods

## gson.Decode(runes []rune) (interface{}, error)

Accepts slice of runes and returns []interface{} or map[string]interface{} in case of success or nil, error in case of failure

## gson.Encode(source interface{}) string

Encodes provided interface{} to json string, interface{} should be either []interface{} or map[string]interface{}


