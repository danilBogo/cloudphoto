package services

func (iniConfig *IniConfig) ToAwsConfig() AwsConfig {
	return AwsConfig{
		AccessKey:   iniConfig.AccessKey,
		SecretKey:   iniConfig.SecretKey,
		Region:      iniConfig.Region,
		EndpointURL: iniConfig.EndpointURL,
	}
}
