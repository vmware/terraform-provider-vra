package tango

func expandCustomProperties(configCustomProperties map[string]interface{}) map[string]string {
	customProperties := make(map[string]string)

	for key, value := range configCustomProperties {
		customProperties[key] = value.(string)
	}

	return customProperties
}
